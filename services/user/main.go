// @title User Service API
// @version 1.0
// @description API for user profile and presence.
// @host localhost:8087
// @BasePath /
// @schemes http
package main

import (
	"context"
	"os"
	"time"
	_ "user/docs"

	"github.com/celesteyang/ChatOrbit/shared/logger"
	"github.com/celesteyang/ChatOrbit/shared/swagger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	logConfig := logger.LogConfig{
		Level:       getEnvOrDefault("LOG_LEVEL", "info"),
		ServiceName: "user",
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}

	if err := logger.InitLogger(logConfig); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	defer logger.Sync()

	servicePort := getEnvOrDefault("PORT", "")
	if servicePort == "" {
		logger.Fatal("PORT environment variable is not set")
	} else {
		logger.Info("Service port", zap.String("port", servicePort))
	}

	logger.Info("Starting user service")
	logger.Debug("Debugging information for user service")
	println("User service is running...")

	// 連接 MongoDB
	mongoURI := getEnvOrDefault("MONGO_URL", "")
	if mongoURI == "" {
		logger.Fatal("MONGO_URL environment variable is not set")
	}

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
	InitCollections(db)

	r := gin.Default()
	r.Use(cors.Default())
	swagger.InitSwagger(r, "User Service")
	r.GET("/user/:id", GetUserHandler)
	// Run the server
	if err := r.Run(":" + servicePort); err != nil {
		logger.Fatal("Failed to run server", zap.Error(err))
	}

}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
