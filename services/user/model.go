package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

type User struct {
	ID    string `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"username"`
	Email string `json:"email" bson:"email"`
}

// InitCollections sets up the MongoDB collections and creates necessary indexes.
func InitCollections(db *mongo.Database) {
	userCollection = db.Collection("users")

	// Create email index to optimize queries.
	_, err := userCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		panic("Failed to create index on users collection: " + err.Error())
	}
}

// GetUserByID fetches a user from the database by ID.
func GetUserByID(ctx context.Context, id string) (*User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user User
	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
