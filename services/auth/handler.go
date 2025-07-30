package main

// HTTP/WebSocket handlers for auth service
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required,min=6"`
}

// @Summary      Register a new user
// @Description  Register a new user with email, username, and password
// @Tags         Register
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /register [post]
// @Param        request body RegisterRequest true "Register request body"
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := RegisterUser(req.Email, req.Username, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// @Summary      Login
// @Description  Login a user with email and password
// @Tags         Login
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /login [post]
// @Param        request body LoginRequest true "Login request body"
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	token, err := LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
