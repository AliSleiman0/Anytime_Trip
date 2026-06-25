package admin

import (
	"context"
	"time"

	adminmodel "travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PopularRepository provides DB access for popular location items
type PopularRepository struct {
	collection *mongo.Collection
}

// NewPopularRepository creates a new PopularRepository
func NewPopularRepository(db *mongo.Database) *PopularRepository {
	return &PopularRepository{collection: db.Collection("popular_locations")}
}

// UpsertByID inserts or updates a popular location by ID
func (r *PopularRepository) UpsertByID(ctx context.Context, item *adminmodel.PopularLocation) error {
	now := time.Now()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": item.ID}, bson.M{"$set": item}, options.Update().SetUpsert(true))
	return err
}

// FindByID returns a popular location by ID
func (r *PopularRepository) FindByID(ctx context.Context, id string) (*adminmodel.PopularLocation, error) {
	var t adminmodel.PopularLocation
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// ListAll returns all popular locations sorted by ID
func (r *PopularRepository) ListAll(ctx context.Context) ([]*adminmodel.PopularLocation, error) {
	cur, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []*adminmodel.PopularLocation
	for cur.Next(ctx) {
		var t adminmodel.PopularLocation
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

// DeleteByID removes a popular location by ID
func (r *PopularRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
