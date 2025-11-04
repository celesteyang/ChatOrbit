// OAuth2 login and user session management
// @title Auth Service API
// @version 1.0
// @description This is the authentication service for ChatOrbit.
// @host localhost:8089
// @BasePath /
// @schemes http
package main

import (
	_ "auth/docs"
	"context"
	"os"
	"time"

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
		ServiceName: "auth",
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}
	servicePort := getEnvOrDefault("PORT", "")
	if servicePort == "" {
		logger.Fatal("PORT environment variable is not set")
	}

	if err := logger.InitLogger(logConfig); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	r := gin.Default()

	// 初始化 Swagger
	swagger.InitSwagger(r, "Auth Service")

	defer logger.Sync()

	logger.Info("Starting auth service....")

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
	InitUserCollection(db)

	// CORS 設定，允許前端跨域並攜帶 Cookie
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 前端網址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, // 允許 Cookie
	}))

	r.POST("/register", RegisterHandler)
	r.POST("/login", LoginHandler)
	r.POST("/change-password", AuthMiddleware(), ChangePasswordHandler)
	r.POST("/logout", AuthMiddleware(), LogoutHandler)

	logger.Debug("Debugging information for auth service")
	r.Run()
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
