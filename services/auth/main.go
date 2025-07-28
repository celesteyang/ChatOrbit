// OAuth2 login and user session management
package main

import (
	"context"
	"os"
	"time"

	"github.com/celesteyang/ChatOrbit/shared/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	logConfig := logger.LogConfig{
		Level:       getEnvOrDefault("LOG_LEVEL", "info"),
		ServiceName: "auth",
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}

	if err := logger.InitLogger(logConfig); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	defer logger.Sync()

	logger.Info("Starting auth service")
	// 連接 MongoDB
	mongoURI := getEnvOrDefault("MONGO_URL", "mongodb://localhost:27017")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Fatal("Mongo client creation failed", zap.Error(err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.Fatal("Mongo connection failed", zap.Error(err))
	}
	db := client.Database("chatorbit")
	InitUserCollection(db)

	r := gin.Default()
	r.POST("/register", RegisterHandler)
	r.POST("/login", LoginHandler)
	logger.Info("Auth service is running on port 8080")
	if err := r.Run(":8080"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
