package admin

import (
	"context"
	"time"

	"travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// HotelRepository handles database operations for hotels
type HotelRepository struct {
	collection *mongo.Collection
}

// NewHotelRepository creates a new hotel repository
func NewHotelRepository(db *mongo.Database) *HotelRepository {
	return &HotelRepository{
		collection: db.Collection("hotels"),
	}
}

// Create creates a new hotel
func (r *HotelRepository) Create(ctx context.Context, hotel *admin.Hotel) error {
	hotel.CreatedAt = time.Now()
	hotel.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, hotel)
	return err
}

// FindByID finds a hotel by its internal ID
func (r *HotelRepository) FindByID(ctx context.Context, id string) (*admin.Hotel, error) {
	var hotel admin.Hotel
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&hotel)
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

// FindByHotelID finds a hotel by its hotel ID (for API lookup)
func (r *HotelRepository) FindByHotelID(ctx context.Context, hotelID string) (*admin.Hotel, error) {
	var hotel admin.Hotel
	err := r.collection.FindOne(ctx, bson.M{"hotel_id": hotelID}).Decode(&hotel)
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

// FindByHotelName finds hotels by name
func (r *HotelRepository) FindByHotelName(ctx context.Context, hotelName string) ([]admin.Hotel, error) {
	filter := bson.M{"hotel_name": hotelName}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var hotels []admin.Hotel
	if err := cursor.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

// FindByStatus finds all hotels with a specific status
func (r *HotelRepository) FindByStatus(ctx context.Context, status admin.HotelStatus) ([]admin.Hotel, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var hotels []admin.Hotel
	if err := cursor.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

// FindAll finds all hotels with optional pagination
func (r *HotelRepository) FindAll(ctx context.Context, limit, skip int64) ([]admin.Hotel, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var hotels []admin.Hotel
	if err := cursor.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

// Update updates a hotel
func (r *HotelRepository) Update(ctx context.Context, id string, hotel *admin.Hotel) error {
	hotel.UpdatedAt = time.Now()
	update := bson.M{"$set": hotel}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateStatus updates the status of a hotel
func (r *HotelRepository) UpdateStatus(ctx context.Context, id string, status admin.HotelStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateCost updates the cost of a hotel room
func (r *HotelRepository) UpdateCost(ctx context.Context, id string, cost float64) error {
	update := bson.M{
		"$set": bson.M{
			"cost":       cost,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateProfitPercent updates the profit percent for a hotel provider
func (r *HotelRepository) UpdateProfitPercent(ctx context.Context, id string, profitPercent float64) error {
	update := bson.M{
		"$set": bson.M{
			"profit_percent": profitPercent,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a hotel
func (r *HotelRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of hotels
func (r *HotelRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByStatus returns the number of hotels by status
func (r *HotelRepository) CountByStatus(ctx context.Context, status admin.HotelStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// UpdateFreeze updates the is_freezed flag for a hotel provider
func (r *HotelRepository) UpdateFreeze(ctx context.Context, id string, freeze bool) error {
	update := bson.M{
		"$set": bson.M{
			"is_freezed": freeze,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
