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
		}

		// Register the client
		client.hub.register <- client

		// Start goroutines to handle read and write
		go HandleClientWrites(client)
		go HandleClientMessages(client)
	}
}
