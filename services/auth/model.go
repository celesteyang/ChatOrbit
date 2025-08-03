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

func InitUserCollection(db *mongo.Database) {
	userCollection = db.Collection("users")
}

func IsEmailExists(email string) (bool, error) {
	filter := bson.M{"email": email}
	count, err := userCollection.CountDocuments(context.TODO(), filter)
	return count > 0, err
}

func InsertUser(user User) error {
	_, err := userCollection.InsertOne(context.TODO(), user)
	return err
}

func FindUserByEmail(email string) (*User, error) {
	filter := bson.M{"email": email}
	var user User
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserPassword(userID primitive.ObjectID, newHashedPwd string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"password_hash": newHashedPwd,
			"update_time":   time.Now(),
		},
	}
	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func UpdateLoginInfo(userID primitive.ObjectID, ip string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"login_time": time.Now(),
			"login_ip":   ip,
		},
	}
	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	return err
}
