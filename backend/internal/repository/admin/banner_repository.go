package admin

import (
	"context"
	"time"

	adminmodel "Anytime_Travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BannerRepository provides DB access for banner items
type BannerRepository struct {
	collection *mongo.Collection
}

// NewBannerRepository creates a new BannerRepository
func NewBannerRepository(db *mongo.Database) *BannerRepository {
	return &BannerRepository{collection: db.Collection("banners")}
}

// UpsertByID inserts or updates a banner by ID
func (r *BannerRepository) UpsertByID(ctx context.Context, item *adminmodel.Banner) error {
	now := time.Now()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": item.ID}, bson.M{"$set": item}, options.Update().SetUpsert(true))
	return err
}

// FindByID returns a banner by ID
func (r *BannerRepository) FindByID(ctx context.Context, id string) (*adminmodel.Banner, error) {
	var b adminmodel.Banner
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ListAll returns all banners sorted by ID
func (r *BannerRepository) ListAll(ctx context.Context) ([]*adminmodel.Banner, error) {
	cur, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var items []*adminmodel.Banner
	for cur.Next(ctx) {
		var b adminmodel.Banner
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

// DeleteByID removes a banner by ID
func (r *BannerRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
