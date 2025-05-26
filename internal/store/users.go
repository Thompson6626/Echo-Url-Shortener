package store

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Password Password `bson:"password,omitempty" json:"-"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	IsActive  bool  `bson:"is_active" json:"is_active"`
}

type Password struct {
	Text *string `bson:"text,omitempty" json:"-"`
	Hash []byte  `bson:"hash" json:"-"`
}

type UsersStore struct {
	db *mongo.Collection
}

