package app

import (
	"context"
	"time"

	"Anytime_Travel/backend/internal/models/app"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// OTPRepository handles database operations for OTPs
type OTPRepository struct {
	collection *mongo.Collection
}

// NewOTPRepository creates a new OTP repository
func NewOTPRepository(db *mongo.Database) *OTPRepository {
	return &OTPRepository{
		collection: db.Collection("otps"),
	}
}

// Create creates a new OTP
func (r *OTPRepository) Create(ctx context.Context, otp *app.OTP) error {
	_, err := r.collection.InsertOne(ctx, otp)
	return err
}

// FindByEmailAndCode finds an OTP by email and code
func (r *OTPRepository) FindByEmailAndCode(ctx context.Context, email, code string) (*app.OTP, error) {
	var otp app.OTP
	err := r.collection.FindOne(ctx, bson.M{
		"email": email,
		"code":  code,
		"type":  app.OTPTypeEmail,
	}).Decode(&otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// FindByPhoneAndCode finds an OTP by phone and code
func (r *OTPRepository) FindByPhoneAndCode(ctx context.Context, phone, code string) (*app.OTP, error) {
	var otp app.OTP
	err := r.collection.FindOne(ctx, bson.M{
		"phone": phone,
		"code":  code,
		"type":  app.OTPTypePhone,
	}).Decode(&otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// MarkAsVerified marks an OTP as verified
func (r *OTPRepository) MarkAsVerified(ctx context.Context, id string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{"verified": true},
	})
	return err
}

// DeleteExpired deletes expired OTPs
func (r *OTPRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
	})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// DeleteByUserAndType deletes all OTPs for a user and type (to prevent multiple active OTPs)
func (r *OTPRepository) DeleteByUserAndType(ctx context.Context, userID string, otpType app.OTPType) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{
		"user_id": userID,
		"type":    otpType,
	})
	return err
}
