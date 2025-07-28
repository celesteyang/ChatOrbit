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
