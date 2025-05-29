package store

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Urls  interface{
		Create(context.Context, *ShortURL) error
		GetByShortCode(context.Context, string) (string, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetById(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
	}
}

func NewStorage(db *mongo.Database) Storage {
	return Storage{
		Urls:  &ShortUrlsStore{db.Collection("urls")},
		Users: &UserStore{db.Collection("users")},
	}
}
