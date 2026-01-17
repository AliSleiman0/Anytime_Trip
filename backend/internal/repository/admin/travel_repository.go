package admin

import (
	"context"
	"time"

	adminmodel "Anytime_Travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TravelRepository provides DB access for travel experience items
type TravelRepository struct {
	collection *mongo.Collection
}

// NewTravelRepository creates a new TravelRepository
func NewTravelRepository(db *mongo.Database) *TravelRepository {
	return &TravelRepository{collection: db.Collection("travel_experiences")}
}

// UpsertByID inserts or updates a travel experience by ID
func (r *TravelRepository) UpsertByID(ctx context.Context, item *adminmodel.TravelExperience) error {
	now := time.Now()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": item.ID}, bson.M{"$set": item}, options.Update().SetUpsert(true))
	return err
}

// FindByID returns a travel experience by ID
func (r *TravelRepository) FindByID(ctx context.Context, id string) (*adminmodel.TravelExperience, error) {
	var t adminmodel.TravelExperience
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ListAll returns all travel experiences sorted by ID
func (r *TravelRepository) ListAll(ctx context.Context) ([]*adminmodel.TravelExperience, error) {
	cur, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []*adminmodel.TravelExperience
	for cur.Next(ctx) {
		var t adminmodel.TravelExperience
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		items = append(items, &t)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// DeleteByID removes a travel experience by ID
func (r *TravelRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
