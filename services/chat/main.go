// Chat logic: rooms, messaging, typing indicators
// @title Chat Service API
// @version 1.0
// @description  A scalable real-time chat backend for live streams.
// @host localhost:8088
// @BasePath /
// @schemes http
package main

import (
	_ "chat/docs"
	"context"
	"os"
	"time"

	"github.com/celesteyang/ChatOrbit/shared/logger"
	"github.com/celesteyang/ChatOrbit/shared/swagger"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	logConfig := logger.LogConfig{
		Level:       getEnvOrDefault("LOG_LEVEL", "info"),
		ServiceName: "chat",
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}
	if err := logger.InitLogger(logConfig); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	} else {
		logger.Info("Logger initialized", zap.String("level", logConfig.Level))
	}
	defer logger.Sync()

	servicePort := getEnvOrDefault("PORT", "")
	if servicePort == "" {
		logger.Fatal("PORT environment variable is not set")
	} else {
		logger.Info("Service port", zap.String("port", servicePort))
	}

	// connect MongoDB
	mongoURI := getEnvOrDefault("MONGO_URL", "")
	if mongoURI == "" {
		logger.Fatal("MONGO_URL environment variable is not set")
	}
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Fatal("Mongo client creation failed", zap.Error(err))
	} else {
		logger.Info("Mongo client created")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = mongoClient.Connect(ctx); err != nil {
		logger.Fatal("Mongo connection failed", zap.Error(err))
	} else {
		logger.Info("Mongo connected")
	}
	mongoDB := mongoClient.Database("chatorbit")
	InitCollections(mongoDB)

	// connect Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: getEnvOrDefault("REDIS_ADDR", ""),
		DB:   0,
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		logger.Fatal("Redis connection failed", zap.Error(err))
	}

	logger.Info("Starting chat service")

	r := gin.Default()
	swagger.InitSwagger(r, "Chat Service")
	hub := NewHub(redisClient)
	// hub instance run in a separate goroutine
	go hub.Run()

	// Define routes and pass Hub instance to handlers
	r.GET("/ws/chat", ChatWebSocketHandler(hub))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World!"})
	})
	// RESTful API for chat history
	r.GET("/chat/history/:roomID", GetChatHistoryHandler)
	r.POST("/chat/rooms", CreateRoomHandler)
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
