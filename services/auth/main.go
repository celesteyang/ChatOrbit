// OAuth2 login and user session management
package main

import (
	"os"

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

	// Your gateway logic here...
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
