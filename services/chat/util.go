package main

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// Note: jwtSecret is used in ValidateJWT function below.
// jwtSecret is declared in handler.go and used here.
// init loads the JWT secret from environment variable at startup.
func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable not set")
	}
	jwtSecret = []byte(secret)
}

// ValidateJWT parses and validates a JWT token string.
// Returns claims (user info) if valid, or error if invalid.
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}
