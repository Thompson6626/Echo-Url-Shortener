package store

import (
	"Url-Shortener/internal/base62"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

// 15 days
var defaultExpiration = 24 * time.Hour * 15

var mu sync.Mutex
var currentSequence int64 = 0
var machineID int64 = 164

const maxSequence int64 = 9999

func nextID() int64 {
	mu.Lock()
	defer mu.Unlock()

	currentSequence = (currentSequence + 1) % maxSequence
	id := machineID*10000 + currentSequence
	return id
}

type ShortURL struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	ShortCode   string             `bson:"short_code" json:"short_code"`                     // Unique short identifier
	OriginalURL string             `bson:"original_url" json:"original_url"`                 // The original long URL
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`                           // Reference to User
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`                     // Creation timestamp
	ExpiresAt   *time.Time         `bson:"expires_at,omitempty" json:"expires_at,omitempty"` // Optional expiration timestamp
	VisitCount  uint64             `bson:"visit_count" json:"visit_count"`                   // Total visit count
}

type ShortUrlsStore struct {
	collection *mongo.Collection
}

func (s *ShortUrlsStore) Create(ctx context.Context, shortURL *ShortURL) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	now := time.Now()

	id := nextID()
	shortCode := base62.Encode(id)

	shortURL.ShortCode = shortCode

	if shortURL.CreatedAt.IsZero() {
		shortURL.CreatedAt = now
	}
	if shortURL.ExpiresAt == nil {
		exp := now.Add(defaultExpiration)
		shortURL.ExpiresAt = &exp
	}

	shortURL.VisitCount = 0

	_, err := s.collection.InsertOne(ctx, shortURL)
	return err
}

func (s *ShortUrlsStore) GetByShortCode(ctx context.Context, shortCode string) (*ShortURL, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// Find and increment VisitCount atomically
	filter := bson.M{"short_code": shortCode}
	update := bson.M{"$inc": bson.M{"visit_count": 1}}

	var updated ShortURL
	err := s.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().
			SetReturnDocument(options.After), // Return updated doc
	).Decode(&updated)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &updated, nil
}

func (s *ShortUrlsStore) GetAllUrlsByUser(ctx context.Context, userID primitive.ObjectID) ([]ShortURL, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	filter := bson.M{"user_id": userID}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var urls []ShortURL
	if err = cursor.All(ctx, &urls); err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *ShortUrlsStore) Delete(ctx context.Context, shortCode string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	filter := bson.M{"short_code": shortCode}
	_, err := s.collection.DeleteOne(ctx, filter)
	return err
}
