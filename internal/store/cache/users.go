package cache

import (
	"Url-Shortener/internal/store"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (s *UserStore) Get(ctx context.Context, userID primitive.ObjectID) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%s", userID.Hex())

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%s", user.ID.Hex())

	jsonn, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.SetEx(ctx, cacheKey, jsonn, UserExpTime).Err()
}

func (s *UserStore) Delete(ctx context.Context, userID primitive.ObjectID) {
	cacheKey := fmt.Sprintf("user-%s", userID.Hex())
	s.rdb.Del(ctx, cacheKey)
}
