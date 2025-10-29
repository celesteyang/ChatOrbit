package main

// Business logic for chat service
import (
	"context"
	"encoding/json"
	"time"

	"github.com/celesteyang/ChatOrbit/shared/logger"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Represents a single WebSocket connection, including user info and message send channel.
type client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	user *UserClaims
}

// Coordinates all client connections and handles message broadcasting.
type Hub struct {
	clients    map[*client]bool
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
	redis      *redis.Client
	rooms      map[string]bool
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
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		redis:      redisClient,
		rooms:      make(map[string]bool),
	}
}

// Starts the main event loop for the Hub, listens for register, unregister, and broadcast events and handles them accordingly.
func (h *Hub) Run() {
	roomID := "general"
	h.rooms[roomID] = true
	h.subscribeToRoom(roomID)

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			logger.Info("Client registered", zap.String("userID", client.user.UserID))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				logger.Info("Client unregistered", zap.String("userID", client.user.UserID))
			}
		case message := <-h.broadcast:
			// Broadcast the message received from Redis to all local connections.
			// Now it's simplified to send to all connections, futurely can filter by room.
			for client := range h.clients {
				select {
				case client.send <- message:
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
			h.broadcast <- []byte(msg.Payload)
		}
	}()
}

// Handlers reading messages from the WebSocket connection, saving them to the database, and publishing them to Redis.
func HandleClientMessages(c *client) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

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
	defer c.conn.Close()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
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
