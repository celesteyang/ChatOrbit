package main

// Business logic for chat service
import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/celesteyang/ChatOrbit/shared/logger"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Heartbeat and presence timing configuration.
const (
	writeWait   = 10 * time.Second
	pongWait    = 25 * time.Second
	pingPeriod  = 10 * time.Second
	presenceTTL = 30 * time.Second
)

// Represents a single WebSocket connection, including user info and message send channel.
type client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	user   *UserClaims
	roomID string
}

// Coordinates all client connections and handles message broadcasting.
type Hub struct {
	clients    map[*client]bool
	broadcast  chan BroadcastMessage
	register   chan *client
	unregister chan *client
	redis      *redis.Client
	rooms      map[string]bool

	// roomsMu protects concurrent access to the rooms map when clients join
	// new rooms outside of the Hub event loop (e.g., when switching rooms).
	roomsMu sync.Mutex
}

type BroadcastMessage struct {
	RoomID  string
	Payload []byte
}

// Stores user information extracted from JWT.
type UserClaims struct {
	UserID string
	Email  string
}

// Creates and returns a new Hub instance.
func NewHub(redisClient *redis.Client) *Hub {
	return &Hub{
		clients:    make(map[*client]bool),
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *client),
		unregister: make(chan *client),
		redis:      redisClient,
		rooms:      make(map[string]bool),
	}
}

// Starts the main event loop for the Hub, listens for register, unregister, and broadcast events and handles them accordingly.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if client.roomID == "" {
				client.roomID = "general"
			}
			h.ensureRoomSubscription(client.roomID)
			h.clients[client] = true
			if err := h.trackPresence(context.Background(), client.roomID, client.user.UserID); err != nil {
				logger.Error("Failed to track presence", zap.Error(err))
			}
			logger.Info("Client registered", zap.String("userID", client.user.UserID))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				if err := h.removePresence(context.Background(), client.roomID, client.user.UserID); err != nil {
					logger.Error("Failed to remove presence", zap.Error(err))
				}
				logger.Info("Client unregistered", zap.String("userID", client.user.UserID))
			}
		case message := <-h.broadcast:
			// Broadcast the message received from Redis to all local connections.
			// Now it's simplified to send to all connections, futurely can filter by room.
			for client := range h.clients {
				if client.roomID != message.RoomID {
					continue
				}
				select {
				case client.send <- message.Payload:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Subscribes to a Redis Pub/Sub channel for a specific chat room, when there are new messages, forward them to Hub's broadcast channel.
func (h *Hub) subscribeToRoom(roomID string) {
	pubsub := h.redis.Subscribe(context.Background(), "chat_room:"+roomID)
	go func() {
		defer pubsub.Close()
		for {
			msg, err := pubsub.ReceiveMessage(context.Background())
			if err != nil {
				logger.Error("Error receiving message from Redis PubSub", zap.Error(err))
				break
			}
			h.broadcast <- BroadcastMessage{RoomID: roomID, Payload: []byte(msg.Payload)}
		}
	}()
}

// ensureRoomSubscription subscribes the Hub to a room channel if it has not done so yet.
func (h *Hub) ensureRoomSubscription(roomID string) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	if _, ok := h.rooms[roomID]; ok {
		return
	}

	h.rooms[roomID] = true
	h.subscribeToRoom(roomID)
	logger.Info("Subscribed to room", zap.String("roomID", roomID))
}

// Handlers reading messages from the WebSocket connection, saving them to the database, and publishing them to Redis.
func HandleClientMessages(c *client) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Set initial read deadline so connections that stop responding to heartbeats are cleaned up.
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("WebSocket read error", zap.Error(err))
			}
			break
		}

		var incomingMessage Message
		if err := json.Unmarshal(message, &incomingMessage); err != nil {
			logger.Error("Failed to parse incoming message", zap.Error(err))
			continue
		}

		incomingMessage.UserID = c.user.UserID
		targetRoom := incomingMessage.RoomID
		if targetRoom == "" {
			targetRoom = c.roomID
		}
		if targetRoom != c.roomID {
			if err := c.hub.switchClientRoom(context.Background(), c, targetRoom); err != nil {
				logger.Error("Failed to switch client room", zap.Error(err))
				targetRoom = c.roomID
			}
		}
		incomingMessage.RoomID = targetRoom
		incomingMessage.Timestamp = time.Now()

		// Save the message to the database by calling the model layer function.
		if err := InsertMessage(context.Background(), &incomingMessage); err != nil {
			logger.Error("Failed to insert message into DB", zap.Error(err))
			continue
		}

		// Publish the message to Redis.
		msgJSON, _ := json.Marshal(incomingMessage)
		if err := c.hub.redis.Publish(context.Background(), "chat_room:"+incomingMessage.RoomID, msgJSON).Err(); err != nil {
			logger.Error("Failed to publish message to Redis", zap.Error(err))
		}
	}
}

// Write messages from the Hub to the WebSocket connection.
func HandleClientWrites(c *client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Reverse the order to get oldest first (for UI display)
func ReverseMessages(messages []Message) {
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
}

// switchClientRoom moves a connection into a different room, updating presence and subscriptions.
func (h *Hub) switchClientRoom(ctx context.Context, c *client, newRoomID string) error {
	if newRoomID == "" || newRoomID == c.roomID {
		return nil
	}

	if err := EnsureRoomExists(ctx, newRoomID); err != nil {
		return err
	}

	if err := h.trackPresence(ctx, newRoomID, c.user.UserID); err != nil {
		return err
	}

	if err := h.removePresence(ctx, c.roomID, c.user.UserID); err != nil {
		return err
	}

	c.roomID = newRoomID
	h.ensureRoomSubscription(newRoomID)

	return nil
}

func presenceKey(roomID string) string {
	return "presence:room:" + roomID
}

func presenceMemberKey(roomID, userID string) string {
	return fmt.Sprintf("presence:room:%s:user:%s", roomID, userID)
}

func (h *Hub) trackPresence(ctx context.Context, roomID, userID string) error {
	pipe := h.redis.TxPipeline()
	pipe.SAdd(ctx, presenceKey(roomID), userID)
	pipe.Set(ctx, presenceMemberKey(roomID, userID), "1", presenceTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (h *Hub) refreshPresence(ctx context.Context, roomID, userID string) error {
	pipe := h.redis.TxPipeline()
	pipe.SAdd(ctx, presenceKey(roomID), userID)
	pipe.Set(ctx, presenceMemberKey(roomID, userID), "1", presenceTTL)
	_, err := pipe.Exec(ctx)
	return err
}

func (h *Hub) removePresence(ctx context.Context, roomID, userID string) error {
	pipe := h.redis.TxPipeline()
	pipe.SRem(ctx, presenceKey(roomID), userID)
	pipe.Del(ctx, presenceMemberKey(roomID, userID))
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

func (h *Hub) GetRoomPresenceCount(ctx context.Context, roomID string) (int64, error) {
	userIDs, err := h.redis.SMembers(ctx, presenceKey(roomID)).Result()
	if err != nil {
		return 0, err
	}

	var (
		activeCount int64
		staleUsers  []interface{}
	)

	for _, userID := range userIDs {
		ttl, err := h.redis.TTL(ctx, presenceMemberKey(roomID, userID)).Result()
		if err != nil {
			return 0, err
		}
		if ttl > 0 {
			activeCount++
			continue
		}
		staleUsers = append(staleUsers, userID)
	}

	if len(staleUsers) > 0 {
		if err := h.redis.SRem(ctx, presenceKey(roomID), staleUsers...).Err(); err != nil {
			logger.Error("Failed to clean stale presence entries", zap.Error(err))
		}
	}

	return activeCount, nil
}
