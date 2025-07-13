package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongo_url := "mongodb://localhost:27019"
	redis_url := "localhost:6381"

	// Test MongoDB connection
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongo_url))
	if err != nil {
		fmt.Println("MongoDB connection error:", err)
	} else {
		err = mongoClient.Ping(ctx, nil)
		if err != nil {
			fmt.Println("MongoDB ping error:", err)
		} else {
			fmt.Println("MongoDB connection: OK")
		}
		_ = mongoClient.Disconnect(ctx)
	}

	// Test Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: redis_url,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis connection error:", err)
	} else {
		fmt.Println("Redis connection: OK")
	}
}
