package main

// Service-specific data types for auth service
import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username,omitempty"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash,omitempty" json:"-"`
	CreateTime   time.Time          `bson:"created_at" json:"created_at"`
	LoginTime    time.Time          `bson:"login_time,omitempty" json:"login_time,omitempty"`
	IPAddress    string             `bson:"login_ip,omitempty" json:"login_ip,omitempty"`
	UpdateTime   time.Time          `bson:"update_time,omitempty" json:"update_time,omitempty"`
}

var userCollection *mongo.Collection

// Set the MongoDB collection for users.
func InitUserCollection(db *mongo.Database) {
	userCollection = db.Collection("users")
}

// Check if the given email is already registered.
func IsEmailExists(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	count, err := userCollection.CountDocuments(ctx, filter)
	return count > 0, err
}

// Insert a new user into the database.
func InsertUser(ctx context.Context, user User) error {
	_, err := userCollection.InsertOne(ctx, user)
	return err
}

// Find user by email.
func FindUserByEmail(ctx context.Context, email string) (*User, error) {
	filter := bson.M{"email": email}
	var user User
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update the password.
func UpdateUserPassword(ctx context.Context, userID primitive.ObjectID, newHashedPwd string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"password_hash": newHashedPwd,
			"update_time":   time.Now(),
		},
	}
	_, err := userCollection.UpdateOne(ctx, filter, update)
	return err
}

// Update the login info.
func UpdateLoginInfo(ctx context.Context, userID primitive.ObjectID, ip string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"login_time": time.Now(),
			"login_ip":   ip,
		},
	}
	_, err := userCollection.UpdateOne(ctx, filter, update)
	return err
}

// Find user by ObjectID.
func FindUserByID(ctx context.Context, userID primitive.ObjectID) (*User, error) {
	filter := bson.M{"_id": userID}
	var user User
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
