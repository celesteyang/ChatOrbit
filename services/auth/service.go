package main

// Business logic for auth service
import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(email, username string, password string) error {
	// Check if email already exists
	exists, err := IsEmailExists(email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already registered")
	}

	// Password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create User object
	user := User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	return InsertUser(user)
}

func LoginUser(email, password string) (string, error) {
	user, err := FindUserByEmail(email)
	if err != nil {
		return "", errors.New("email or password is incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("email or password is incorrect")
	}

	token, err := GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

var jwtSecret = []byte("jwtTestY&771765454330an")

func GenerateJWT(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Register Test
// curl -X POST http://localhost:8080/register   -H "Content-Type: application/json"   -d '{"email":"abc@example.com", "username":"test", "password":"12345678"}'
// Login Test
// curl -X POST http://localhost:8080/login   -H "Content-Type: application/json"   -d '{"email":"abc@example.com", "password":"12345678"}'
