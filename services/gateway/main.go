// HTTP/WebSocket entry point for ChatOrbit
package main

import (
	"os"
	"time"

	"github.com/celesteyang/ChatOrbit/shared/logger"
)

func main() {
	logConfig := logger.LogConfig{
		Level:       getEnvOrDefault("LOG_LEVEL", "info"),
		ServiceName: "gateway",
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}

	if err := logger.InitLogger(logConfig); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	defer logger.Sync()

	logger.Info("Starting gateway service")
	logger.Debug("Debugging information for gateway service")
	println("Gateway service is running...")

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

// func main() {
// 	// This is the main entry point for the gateway service.
// 	// The actual implementation would go here, such as setting up routes,
// 	// initializing the database connections, and starting the server.
// 	// Create a logger with configurable log level

// 	cfg := zap.NewProductionConfig()
// 	loggingLevel := os.Getenv("LOGGING_LEVEL")
// 	var level zap.AtomicLevel
// 	switch loggingLevel {
// 	case "debug":
// 		level = zap.NewAtomicLevelAt(zap.DebugLevel)
// 	case "warn":
// 		level = zap.NewAtomicLevelAt(zap.WarnLevel)
// 	case "error":
// 		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
// 	default:
// 		level = zap.NewAtomicLevelAt(zap.InfoLevel)
// 	}

// 	cfg.Level = level
// 	logger, _ := cfg.Build()
// 	defer logger.Sync()

// 	logger.Debug("This is a debug message")
// 	logger.Info("This is an info message")
// 	logger.Warn("This is a warning")
// 	logger.Error("This is an error")

// 	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
// 	chatServiceURL := os.Getenv("CHAT_SERVICE_URL")
// 	if authServiceURL == "" {
// 		authServiceURL = "http://auth-service:8080"
// 	}
// 	if chatServiceURL == "" {
// 		chatServiceURL = "http://chat-service:8080"
// 	}
// 	logger.Info("Auth Service URL:", zap.String("url", authServiceURL))
// 	logger.Info("Chat Service URL:", zap.String("url", chatServiceURL))

// 	for {
// 		time.Sleep(1 * time.Second)
// 	}
// }
