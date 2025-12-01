package main

// HTTP/WebSocket handlers for chat service
import (
	"net/http"
	"strings"

	"github.com/celesteyang/ChatOrbit/shared/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type createRoomRequest struct {
	RoomID string `json:"room_id" binding:"required"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Reads messages from the WebSocket connection and broadcasts them to the Hub.
// Verify JWT token from query parameter or Authorization header.
// On success, register the client and start read/write goroutines.
// Pass the Hub instance to manage the client connection.
// Usage: r.GET("/ws/chat", ChatWebSocketHandler(hub))
func ChatWebSocketHandler(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		roomID := c.DefaultQuery("room_id", "general")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
			return
		}

		if err := EnsureRoomExists(c.Request.Context(), roomID); err != nil {
			logger.Error("Failed to ensure room exists", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
			return
		}

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token claims"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error("Failed to upgrade WebSocket connection", zap.Error(err))
			return
		}

		client := &client{
			hub:  hub,
			conn: conn,
			send: make(chan []byte, 256),
			user: &UserClaims{
				UserID: userID,
				Email:  claims["email"].(string),
			},
			roomID: roomID,
		}

		// Register the client
		client.hub.register <- client

		// Start goroutines to handle read and write
		go HandleClientWrites(client)
		go HandleClientMessages(client)
	}
}

// @Summary Get room presence
// @Description Returns the number of users currently connected to a room.
// @Tags Chat
// @Produce json
// @Param roomID path string true "Chat Room ID"
// @Success 200 {object} map[string]interface{}
// @Router /chat/rooms/{roomID}/presence [get]
func GetRoomPresenceHandler(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("roomID")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
			return
		}

		count, err := hub.GetRoomPresenceCount(c.Request.Context(), roomID)
		if err != nil {
			logger.Error("Failed to get room presence", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get presence"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"room_id": roomID, "online": count})
	}
}

// @Summary Get chat history
// @Description Retrieves chat messages from a specific room.
// @Tags Chat
// @Accept json
// @Produce json
// @Param roomID path string true "Chat Room ID"
// @Success 200 {array} Message
// @Router /chat/history/{roomID} [get]
func GetChatHistoryHandler(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	if err := EnsureRoomExists(c.Request.Context(), roomID); err != nil {
		logger.Error("Failed to ensure room exists", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
		return
	}

	// Call the model function to get messages.
	// You might want to add pagination (e.g., limit, offset) here.
	messages, err := GetMessagesByRoom(c.Request.Context(), roomID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// @Summary Create chat room
// @Description Creates a chat room by ID. Idempotent: returns success if the room already exists.
// @Tags Chat
// @Accept json
// @Produce json
// @Param room body createRoomRequest true "Room info"
// @Success 200 {object} map[string]interface{}
// @Router /chat/rooms [post]
func CreateRoomHandler(c *gin.Context) {
	var req createRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id is required"})
		return
	}

	if err := EnsureRoomExists(c.Request.Context(), req.RoomID); err != nil {
		logger.Error("Failed to ensure room exists", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room_id": req.RoomID})
}
