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
	"os"

	"github.com/celesteyang/ChatOrbit/shared/logger"
	"github.com/celesteyang/ChatOrbit/shared/swagger"
	"github.com/gin-gonic/gin"
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

	r.Run(":8083")

	defer logger.Sync()

	logger.Info("Starting auth service")
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
