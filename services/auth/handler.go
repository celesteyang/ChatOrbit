package main

// HTTP/WebSocket handlers for auth service
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

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
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Router       /register [post]
// @Param        request body RegisterRequest true "Register request body"
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input"})
		return
	}
	// Get client IP address
	ip := c.ClientIP()

	if err := RegisterUser(req.Email, req.Username, req.Password, ip); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Registration successful"})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

// @Summary      Login
// @Description  Login a user with email and password
// @Tags         Login
// @Accept       json
// @Produce      json
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /login [post]
// @Param        request body LoginRequest true "Login request body"
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input"})
		return
	}

	token, err := LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}
	// Set JWT in HttpOnly Cookie
	// c.SetCookie("token", token, 3600*24, "/", "localhost", true, true)
	// without frontend cookie
	c.SetCookie("token", token, 3600*24, "/", "", false, true)
	c.JSON(200, gin.H{"token": token})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// @Summary      Change Password
// @Description  Change the password of the logged-in user
// @Tags         ChangePassword
// @Accept       json
// @Produce      json
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /change-password [post]
// @Param        request body ChangePasswordRequest true "Change Password request body"
func ChangePasswordHandler(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input"})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	err := ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Password changed successfully"})
}

// @Summary      Logout
// @Description  Logout the user by clearing the JWT cookie
// @Tags         Logout
// @Accept       json
// @Produce      json
// @Success      200  {object}  MessageResponse
// @Router       /logout [post]
func LogoutHandler(c *gin.Context) {
	// c.SetCookie("token", "", -1, "/", "localhost", true, true)
	// without frontend cookie
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, MessageResponse{Message: "Logged out successfully"})
}
