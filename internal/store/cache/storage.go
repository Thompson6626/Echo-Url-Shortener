package cache

import (
	"Url-Shortener/internal/store"
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Storage struct {
	Users interface {
		Get(context.Context, primitive.ObjectID) (*store.User, error)
		Set(context.Context, *store.User) error
		Delete(context.Context, primitive.ObjectID)
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rbd},
	}
}
