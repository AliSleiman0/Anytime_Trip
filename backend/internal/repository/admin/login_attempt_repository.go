package admin

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LoginAttemptRepository handles login attempts and temporary blocks.
type LoginAttemptRepository struct {
	attemptsColl *mongo.Collection
	blocksColl   *mongo.Collection
}

// NewLoginAttemptRepository creates a new repository for login attempts/blocks.
func NewLoginAttemptRepository(db *mongo.Database) *LoginAttemptRepository {
	return &LoginAttemptRepository{
		attemptsColl: db.Collection("login_attempts"),
		blocksColl:   db.Collection("login_blocks"),
	}
}

// RecordAttempt inserts a timestamped login attempt for the given key (e.g. "ip:1.2.3.4" or "email:foo@x").
func (r *LoginAttemptRepository) RecordAttempt(ctx context.Context, key string) error {
	_, err := r.attemptsColl.InsertOne(ctx, bson.M{"key": key, "created_at": time.Now()})
	return err
}

// CountAttemptsSince counts attempts for the key since the provided time.
func (r *LoginAttemptRepository) CountAttemptsSince(ctx context.Context, key string, since time.Time) (int64, error) {
	return r.attemptsColl.CountDocuments(ctx, bson.M{"key": key, "created_at": bson.M{"$gte": since}})
}

// CreateBlock creates or updates a block record for the given key until the provided time.
func (r *LoginAttemptRepository) CreateBlock(ctx context.Context, key string, until time.Time) error {
	_, err := r.blocksColl.UpdateOne(ctx, bson.M{"key": key}, bson.M{"$set": bson.M{"key": key, "blocked_until": until}}, options.Update().SetUpsert(true))
	return err
}

// IsBlocked returns whether the key is currently blocked and the block expiry time.
func (r *LoginAttemptRepository) IsBlocked(ctx context.Context, key string) (bool, time.Time, error) {
	var res struct {
		Key          string    `bson:"key"`
		BlockedUntil time.Time `bson:"blocked_until"`
	}
	err := r.blocksColl.FindOne(ctx, bson.M{"key": key}).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, time.Time{}, nil
		}
		return false, time.Time{}, err
	}
	return res.BlockedUntil.After(time.Now()), res.BlockedUntil, nil
}
