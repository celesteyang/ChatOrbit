// OAuth2 login and user session management
package main

import (
	"os"

	"github.com/celesteyang/ChatOrbit/tree/commonUtilities/shared/logger"
)

func main() {
	// This is the main entry point for the auth service.
	// The actual implementation would go here, such as setting up routes,
	// initializing the database connections, and starting the server.
	println("Auth service is running...")

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

	// Your gateway logic here...
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
