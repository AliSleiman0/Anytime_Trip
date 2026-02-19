package admin

import (
	"context"
	"time"

	adminmodel "Anytime_Travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SearchBannerRepository provides DB access for search banner items
type SearchBannerRepository struct {
	collection *mongo.Collection
}

// NewSearchBannerRepository creates a new SearchBannerRepository
func NewSearchBannerRepository(db *mongo.Database) *SearchBannerRepository {
	return &SearchBannerRepository{collection: db.Collection("search_banners")}
}

// UpsertByID inserts or updates a search banner by ID
func (r *SearchBannerRepository) UpsertByID(ctx context.Context, item *adminmodel.SearchBanner) error {
	now := time.Now()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": item.ID}, bson.M{"$set": item}, options.Update().SetUpsert(true))
	return err
}

// FindByID returns a search banner by ID
func (r *SearchBannerRepository) FindByID(ctx context.Context, id string) (*adminmodel.SearchBanner, error) {
	var b adminmodel.SearchBanner
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ListAll returns all search banners sorted by ID
func (r *SearchBannerRepository) ListAll(ctx context.Context) ([]*adminmodel.SearchBanner, error) {
	cur, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []*adminmodel.SearchBanner
	for cur.Next(ctx) {
		var b adminmodel.SearchBanner
		if err := cur.Decode(&b); err != nil {
			return nil, err
		}
		items = append(items, &b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// DeleteByID removes a search banner by ID
func (r *SearchBannerRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
