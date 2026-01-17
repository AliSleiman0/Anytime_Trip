package app

import (
	"context"

	appmodels "Anytime_Travel/backend/internal/models/app"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PaymentRepository provides DB access for payments
type PaymentRepository struct {
	collection *mongo.Collection
}

// NewPaymentRepository creates a new repository
func NewPaymentRepository(db *mongo.Database) *PaymentRepository {
	return &PaymentRepository{collection: db.Collection("payments")}
}

// GetAllPayments returns all payments sorted by booking_date desc
func (r *PaymentRepository) GetAllPayments(ctx context.Context) ([]appmodels.Payment, error) {
	opts := options.Find().SetSort(bson.D{{Key: "booking_date", Value: -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var out []appmodels.Payment
	if err := cursor.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// InsertMany inserts multiple payment documents (used by seed)
func (r *PaymentRepository) InsertMany(ctx context.Context, docs []interface{}) (*mongo.InsertManyResult, error) {
	return r.collection.InsertMany(ctx, docs)
}
