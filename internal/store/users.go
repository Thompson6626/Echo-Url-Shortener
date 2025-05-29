package store

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  password           `bson:"password,omitempty" json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
}

type password struct {
	text *string `bson:"text,omitempty" json:"-"`
	Hash []byte  `bson:"hash" json:"-"`
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.Hash = hash

	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.Hash, []byte(text))
}

type UserStore struct {
	collection *mongo.Collection
}

func (s *UserStore) GetById(c context.Context, id int64) (*User, error) {

	return nil, nil
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user.CreatedAt = time.Now()
	user.IsActive = true

	_, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			for _, we := range writeErr.WriteErrors {
				if we.Code == 11000 { // Duplicate key error
					msg := strings.ToLower(we.Message) // normalize case
					if strings.Contains(msg, "unique_email") {
						return ErrDuplicateEmail
					}
					if strings.Contains(msg, "unique_username") {
						return ErrDuplicateUsername
					}
					// Optionally return a generic duplicate key error here if needed
				}
			}
		}
		return err
	}

	return nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
