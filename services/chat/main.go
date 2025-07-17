// Chat logic: rooms, messaging, typing indicators
package main

import (
	"os"
	"time"

	"github.com/celesteyang/ChatOrbit/shared/logger"
)

func main() {
	// This is the main entry point for the chat service.
	// The actual implementation would go here, such as setting up routes,
	// initializing the database connections, and starting the server.

	logConfig := logger.LogConfig{
		Level:       getEnvOrDefault("LOG_LEVEL", "info"),
		ServiceName: "chat",
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}

	if err := logger.InitLogger(logConfig); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	defer logger.Sync()

	logger.Info("Starting chat service")
	logger.Debug("Debugging information for chat service")
	println("Chat service is running...")

	for {
		time.Sleep(1 * time.Second)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
