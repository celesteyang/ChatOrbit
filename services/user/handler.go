package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUserHandler godoc
// @Summary      Get user info
// @Description  Returns user profile by user ID
// @Tags         User
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  User
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /user/{id} [get]
func GetUserHandler(c *gin.Context) {
	id := c.Param("id")
	user, err := GetUserByID(context.Background(), id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}
