package admin

import (
	"context"
	"errors"
	"strings"
	"time"

	"travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// PasswordResetRepository handles database operations for password reset tokens
type PasswordResetRepository struct {
	collection *mongo.Collection
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *mongo.Database) *PasswordResetRepository {
	return &PasswordResetRepository{
		collection: db.Collection("password_resets"),
	}
}

// Create inserts a new password reset token
func (r *PasswordResetRepository) Create(ctx context.Context, reset *admin.PasswordReset) error {
	_, err := r.collection.InsertOne(ctx, reset)
	return err
}

// FindByToken finds a password reset token by token string
func (r *PasswordResetRepository) FindByToken(ctx context.Context, token string) (*admin.PasswordReset, error) {
	var reset admin.PasswordReset
	err := r.collection.FindOne(ctx, bson.M{"token": token, "used": false}).Decode(&reset)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid or expired reset token")
		}
		return nil, err
	}

	// Check if token has expired
	if time.Now().After(reset.ExpiresAt) {
		return nil, errors.New("reset token has expired")
	}

	return &reset, nil
}

// FindByCode finds a password reset record by email and code
func (r *PasswordResetRepository) FindByCode(ctx context.Context, email, code string) (*admin.PasswordReset, error) {
	var reset admin.PasswordReset
	err := r.collection.FindOne(ctx, bson.M{"email": strings.ToLower(email), "code": code, "used": false}).Decode(&reset)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid or expired verification code")
		}
		return nil, err
	}

	// Check if code has expired
	if time.Now().After(reset.ExpiresAt) {
		return nil, errors.New("verification code has expired")
	}

	return &reset, nil
}

// UpdateTokenByID sets the token and updates expiry for a reset record
func (r *PasswordResetRepository) UpdateTokenByID(ctx context.Context, id, token string, expiresAt time.Time) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"token": token, "expires_at": expiresAt}},
	)
	return err
}

// MarkAsUsed marks a reset token as used
func (r *PasswordResetRepository) MarkAsUsed(ctx context.Context, token string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"token": token},
		bson.M{"$set": bson.M{"used": true}},
	)
	return err
}

// DeleteExpired deletes all expired reset tokens
func (r *PasswordResetRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.collection.DeleteMany(
		ctx,
		bson.M{"expires_at": bson.M{"$lt": time.Now()}},
	)
	return err
}

// FindByEmail finds a non-expired, unused reset token by email
func (r *PasswordResetRepository) FindByEmail(ctx context.Context, email string) (*admin.PasswordReset, error) {
	var reset admin.PasswordReset
	err := r.collection.FindOne(ctx, bson.M{
		"email":      strings.ToLower(email),
		"used":       false,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&reset)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &reset, nil
}
