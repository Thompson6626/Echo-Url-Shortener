package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
func New(host string, port string) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri).SetMaxPoolSize(100).SetMaxConnIdleTime(10*time.Second),
	)

	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
