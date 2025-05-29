package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// New creates and returns a MongoDB client with ensured indexes.
func New(host, port, name, username, password string) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", username, password, host, port, name)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(uri).
			SetAuth(options.Credential{
				Username: username,
				Password: password,
			}).
			SetMaxPoolSize(100).
			SetMaxConnIdleTime(10*time.Second),
	)
	if err != nil {
		return nil, err
	}

	if err := ensureIndexes(client.Database(name)); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

// ensureUserIndexes ensures unique indexes exist on "email" and "username"
func ensureUserIndexes(usersCollection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true).SetName("unique_email"),
		},
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true).SetName("unique_username"),
		},
	}

	_, err := usersCollection.Indexes().CreateMany(ctx, indexes)
	return err
}

func ensureUrlIndexes(urlsCollection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"short_code": 1},
			Options: options.Index().SetUnique(true).SetName("unique_short_code"),
		},
		{
			Keys:    bson.M{"user_id": 1},
			Options: options.Index().SetName("by_user_id"),
		},
		{
			Keys:    bson.M{"expires_at": 1},
			Options: options.Index().SetExpireAfterSeconds(0).SetName("expires_index"), // optional TTL if using expiry
		},
	}

	_, err := urlsCollection.Indexes().CreateMany(ctx, indexes)
	return err
}

func ensureIndexes(db *mongo.Database) error {

	err := ensureUserIndexes(db.Collection("users"))
	if err != nil {
		return err
	}

	err = ensureUrlIndexes(db.Collection("urls"))
	if err != nil {
		return err
	}

	return nil
}
