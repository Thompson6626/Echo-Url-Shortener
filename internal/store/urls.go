package store

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ShortURL struct {
    ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`              // Unique identifier
    ShortCode   string     `bson:"short_code" json:"short_code"`         // The short code part of the URL
    OriginalURL string     `bson:"original_url" json:"original_url"`     // The full original URL
    UserID      string     `bson:"user_id" json:"user_id"`                // Owner of the short URL
    CreatedAt   time.Time  `bson:"created_at" json:"created_at"`          // Creation timestamp
    ExpiresAt   *time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"` // Optional expiry time
    VisitCount  uint32     `bson:"visit_count" json:"visit_count"`        // Number of visits
}

type ShortUrlsStore struct {
	db *mongo.Collection
}





