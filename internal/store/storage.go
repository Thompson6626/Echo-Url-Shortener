package store

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Urls interface{}
	Users interface{}
}

func NewStorage(db *mongo.Database) Storage {
	return Storage {
		Urls: &ShortUrlsStore{db.Collection("")},
		Users: &UsersStore{db.Collection("")},
	}
}