package main

// Service-specific data types for chat service
import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	messageCollection *mongo.Collection
)

// Message represents a chat message stored in MongoDB.
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID    string             `bson:"room_id" json:"room_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Content   string             `bson:"content" json:"content"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

// InitCollections sets up the MongoDB collections and creates necessary indexes.
func InitCollections(db *mongo.Database) {
	messageCollection = db.Collection("messages")

	// Create room_id index to optimize queries.
	_, err := messageCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "room_id", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	)
	if err != nil {
		panic("Failed to create index on messages collection: " + err.Error())
	}
}

// Insert the message to the database.
func InsertMessage(ctx context.Context, msg *Message) error {
	_, err := messageCollection.InsertOne(ctx, msg)
	return err
}

// func GetMessages(ctx context.Context, roomID string) ([]Message, error) {
// 	cursor, err := messageCollection.Find(ctx, bson.M{"room_id": roomID})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var messages []Message
// 	for cursor.Next(ctx) {
// 		var msg Message
// 		if err := cursor.Decode(&msg); err != nil {
// 			return nil, err
// 		}
// 		messages = append(messages, msg)
// 	}
// 	return messages, nil
// }

// GetMessagesByRoom retrieves chat messages for a specific room with pagination.
// It sorts messages by timestamp in descending order (newest first).
func GetMessagesByRoom(ctx context.Context, roomID string, limit int64) ([]Message, error) {
	// Find options: sort by timestamp descending, and limit the results.
	findOptions := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(limit)

	cursor, err := messageCollection.Find(ctx, bson.M{"room_id": roomID}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
