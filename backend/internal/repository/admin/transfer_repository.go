package admin

import (
	"context"
	"time"

	"travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransferRepository handles database operations for transfers
type TransferRepository struct {
	collection *mongo.Collection
}

// NewTransferRepository creates a new transfer repository
func NewTransferRepository(db *mongo.Database) *TransferRepository {
	return &TransferRepository{
		collection: db.Collection("transfers"),
	}
}

// Create creates a new transfer
func (r *TransferRepository) Create(ctx context.Context, transfer *admin.Transfer) error {
	transfer.CreatedAt = time.Now()
	transfer.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, transfer)
	return err
}

// FindByID finds a transfer by its internal ID
func (r *TransferRepository) FindByID(ctx context.Context, id string) (*admin.Transfer, error) {
	var transfer admin.Transfer
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&transfer)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

// FindByTransferID finds a transfer by its transfer ID
func (r *TransferRepository) FindByTransferID(ctx context.Context, transferID string) (*admin.Transfer, error) {
	var transfer admin.Transfer
	err := r.collection.FindOne(ctx, bson.M{"transfer_id": transferID}).Decode(&transfer)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

// FindByStatus finds all transfers with a specific status
func (r *TransferRepository) FindByStatus(ctx context.Context, status admin.TransferStatus) ([]admin.Transfer, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transfers []admin.Transfer
	if err := cursor.All(ctx, &transfers); err != nil {
		return nil, err
	}
	return transfers, nil
}

// FindAll finds all transfers with optional pagination
func (r *TransferRepository) FindAll(ctx context.Context, limit, skip int64) ([]admin.Transfer, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transfers []admin.Transfer
	if err := cursor.All(ctx, &transfers); err != nil {
		return nil, err
	}
	return transfers, nil
}

// Update updates a transfer
func (r *TransferRepository) Update(ctx context.Context, id string, transfer *admin.Transfer) error {
	transfer.UpdatedAt = time.Now()
	update := bson.M{"$set": transfer}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateStatus updates the status of a transfer
func (r *TransferRepository) UpdateStatus(ctx context.Context, id string, status admin.TransferStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateProfitPercent updates the profit percent for a transfer provider
func (r *TransferRepository) UpdateProfitPercent(ctx context.Context, id string, profitPercent float64) error {
	update := bson.M{
		"$set": bson.M{
			"profit_percent": profitPercent,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a transfer
func (r *TransferRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of transfers
func (r *TransferRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByStatus returns the number of transfers by status
func (r *TransferRepository) CountByStatus(ctx context.Context, status admin.TransferStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// UpdateFreeze updates the is_freezed flag for a transfer provider
func (r *TransferRepository) UpdateFreeze(ctx context.Context, id string, freeze bool) error {
	update := bson.M{
		"$set": bson.M{
			"is_freezed": freeze,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
