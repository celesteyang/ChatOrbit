package main

// Business logic for auth service
import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(email, username string, password string, ip string) error {
	// Check if email already exists
	exists, err := IsEmailExists(email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Email already registered.")
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
		CreateTime:   time.Now(),
		IPAddress:    ip,
	}

	return InsertUser(user)
}

func LoginUser(email, password string) (string, error) {
	user, err := FindUserByEmail(email)
	if err != nil {
		return "", errors.New("Email is incorrect.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("Password is incorrect.")
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

func FindUserByEmailByID(userID primitive.ObjectID) (*User, error) {
	filter := bson.M{"_id": userID}
	var user User
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func ChangePassword(userID string, oldPassword, newPassword string) error {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("Invalid user account ID.")
	}

	user, err := FindUserByEmailByID(uid)
	if err != nil {
		return errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword))
	if err != nil {
		return errors.New("Old password incorrect.")
	}

	newHashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = UpdateUserPassword(uid, string(newHashed))
	return err
}

// Register Test
// curl -X POST http://localhost:8080/register   -H "Content-Type: application/json"   -d '{"email":"abc@example.com", "username":"test", "password":"12345678"}'
// Login Test
// curl -X POST http://localhost:8080/login   -H "Content-Type: application/json"   -d '{"email":"abc@example.com", "password":"12345678"}'
