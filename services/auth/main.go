// OAuth2 login and user session management
// @title Auth Service API
// @version 1.0
// @description This is the authentication service for ChatOrbit.
// @host localhost:8083
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
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// import docs

func main() {
	logConfig := logger.LogConfig{
		Level:       getEnvOrDefault("LOG_LEVEL", "info"),
		ServiceName: "auth",
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}

	if err := logger.InitLogger(logConfig); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	r := gin.Default()

	// 初始化 Swagger
	swagger.InitSwagger(r, "Auth Service")

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

	r.POST("/register", RegisterHandler)
	r.POST("/login", LoginHandler)

	logger.Info("Auth service is running on port 8080")
	if err := r.Run(":8080"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
	logger.Debug("Debugging information for auth service")

	r.GET("/test", testHandler)
}

// @Summary      Test the auth service
// @Description  This endpoint is used to test if the service is up
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /test [get]
func testHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "test",
	})

}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
