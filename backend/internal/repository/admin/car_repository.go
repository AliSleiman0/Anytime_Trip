package admin

import (
	"context"
	"time"

	"Anytime_Travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CarRepository handles database operations for car rentals
type CarRepository struct {
	collection *mongo.Collection
}

// NewCarRepository creates a new car repository
func NewCarRepository(db *mongo.Database) *CarRepository {
	return &CarRepository{
		collection: db.Collection("cars"),
	}
}

// Create creates a new car rental
func (r *CarRepository) Create(ctx context.Context, car *admin.Car) error {
	car.CreatedAt = time.Now()
	car.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, car)
	return err
}

// FindByID finds a car by its internal ID
func (r *CarRepository) FindByID(ctx context.Context, id string) (*admin.Car, error) {
	var car admin.Car
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&car)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

// FindByCarID finds a car by its car ID (for API lookup)
func (r *CarRepository) FindByCarID(ctx context.Context, carID string) (*admin.Car, error) {
	var car admin.Car
	err := r.collection.FindOne(ctx, bson.M{"car_id": carID}).Decode(&car)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

// FindByCarName finds cars by name
func (r *CarRepository) FindByCarName(ctx context.Context, carName string) ([]admin.Car, error) {
	filter := bson.M{"car_name": carName}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cars []admin.Car
	if err := cursor.All(ctx, &cars); err != nil {
		return nil, err
	}
	return cars, nil
}

// FindByStatus finds all cars with a specific status
func (r *CarRepository) FindByStatus(ctx context.Context, status admin.CarStatus) ([]admin.Car, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cars []admin.Car
	if err := cursor.All(ctx, &cars); err != nil {
		return nil, err
	}
	return cars, nil
}

// FindAll finds all cars with optional pagination
func (r *CarRepository) FindAll(ctx context.Context, limit, skip int64) ([]admin.Car, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cars []admin.Car
	if err := cursor.All(ctx, &cars); err != nil {
		return nil, err
	}
	return cars, nil
}

// Update updates a car
func (r *CarRepository) Update(ctx context.Context, id string, car *admin.Car) error {
	car.UpdatedAt = time.Now()
	update := bson.M{"$set": car}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateStatus updates the status of a car
func (r *CarRepository) UpdateStatus(ctx context.Context, id string, status admin.CarStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateCost updates the cost of a car rental
func (r *CarRepository) UpdateCost(ctx context.Context, id string, cost float64) error {
	update := bson.M{
		"$set": bson.M{
			"cost":       cost,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateProfitPercent updates the profit percent for a car provider
func (r *CarRepository) UpdateProfitPercent(ctx context.Context, id string, profitPercent float64) error {
	update := bson.M{
		"$set": bson.M{
			"profit_percent": profitPercent,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateFreeze updates the is_freezed flag for a car provider
func (r *CarRepository) UpdateFreeze(ctx context.Context, id string, freeze bool) error {
	update := bson.M{
		"$set": bson.M{
			"is_freezed": freeze,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a car
func (r *CarRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of cars
func (r *CarRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByStatus returns the number of cars by status
func (r *CarRepository) CountByStatus(ctx context.Context, status admin.CarStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// SearchCars searches for available cars based on filters
func (r *CarRepository) SearchCars(ctx context.Context, pickupLocation, dropoffLocation string, pickupTime, dropoffTime time.Time, carType string, passengers int) ([]admin.Car, error) {
	filter := bson.M{
		"status": bson.M{"$in": []admin.CarStatus{admin.CarStatusActive}},
	}

	// Filter by pickup location if provided
	if pickupLocation != "" {
		filter["pickup_location"] = bson.M{"$regex": pickupLocation, "$options": "i"}
	}

	// Filter by car type if provided
	if carType != "" {
		filter["car_type"] = bson.M{"$regex": carType, "$options": "i"}
	}

	// Filter by passengers if provided
	if passengers > 0 {
		filter["passengers"] = bson.M{"$gte": passengers}
	}

	// Find cars matching the criteria
	opts := options.Find().SetSort(bson.D{{Key: "cost", Value: 1}}) // Sort by cost ascending
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cars []admin.Car
	if err := cursor.All(ctx, &cars); err != nil {
		return nil, err
	}

	return cars, nil
}
