package main

// HTTP/WebSocket handlers for chat service
import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ChatWebSocketHandler(c *gin.Context) {
	// Accept JWT via query param or header
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

	// Validate JWT
	_, err := ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// TODO: Add connection to room, handle messages, heartbeat, etc.
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// Echo for now
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}
