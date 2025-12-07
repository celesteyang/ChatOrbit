package main

// Business logic for auth service
import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles the user registration process.
//
// It first checks if the email already exists in the database. If not,
// it hashes the user's password using bcrypt, creates a new User object,
// and inserts it into the database.
//
// Parameters:
//
//	ctx: The context for the request.
//	email: The email address for the new user.
//	username: The username for the new user.
//	password: The plaintext password.
//	ip: The user's IP address.
//
// Returns:
//
//	An error if the email is already registered, password hashing fails,
//	or the database insertion fails.
func RegisterUser(ctx context.Context, email, username string, password string, ip string) error {
	// Check if email already exists
	exists, err := IsEmailExists(ctx, email)
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

	return InsertUser(ctx, user)
}

// LoginUser authenticates a user by their email and password.
//
// It first finds a user by their email, then compares the provided password
// with the stored hashed password. If authentication is successful,
// it generates and returns a new JWT.
//
// Parameters:
//
//	ctx: The context for the request.
//	email: The user's email address.
//	password: The plaintext password provided by the user.
//
// Returns:
//
//	A JWT string if the login is successful.
//	An error if the user is not found, the password is incorrect, or token generation fails.
func LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := FindUserByEmail(ctx, email)
	if err != nil {
		log.Println("Login failed (email)", email, err)
		return "", errors.New("Email is incorrect. Error: " + err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Println("Login failed (password)", email, err)
		return "", errors.New("Password is incorrect. Error: " + err.Error())
	}

	token, err := GenerateJWT(user.ID.Hex(), user.Email)
	if err != nil {
		log.Println("JWT error", email, err)
		return "", errors.New("Failed to generate JWT: " + err.Error())
	}

	log.Println("Login success:", user.Email)

	return token, nil
}

// ChangePassword handles updating a user's password.
//
// It first validates the userID, then retrieves the user from the database.
// It verifies the old password against the stored hash. If it matches,
// it hashes the new password and updates the user's record in the database.
//
// Parameters:
//
//	ctx: The context for the request.
//	userID: The unique ID of the user to change the password for.
//	oldPassword: The user's current plaintext password.
//	newPassword: The new plaintext password.
//
// Returns:
//
//	An error if the userID is invalid, the user is not found, the old password
//	is incorrect, or the database update fails.
func ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("Invalid user account ID.")
	}

	user, err := FindUserByID(ctx, uid)
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

	err = UpdateUserPassword(ctx, uid, string(newHashed))
	return err
}
