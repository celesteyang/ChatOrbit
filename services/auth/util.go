// JWT utility functions for auth service
package main

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret stores the secret key for signing JWT tokens.
var jwtSecret []byte

// init loads the JWT secret from environment variable at startup.
func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable not set")
	}
	jwtSecret = []byte(secret)
}

// GenerateJWT creates a JWT token for the given user ID and email.
// The token expires in 24 hours.
func GenerateJWT(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
