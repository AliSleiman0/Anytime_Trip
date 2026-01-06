package admin

import (
	"context"
	"time"

	"Anytime_Travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FlightRepository handles database operations for flights
type FlightRepository struct {
	collection *mongo.Collection
}

// NewFlightRepository creates a new flight repository
func NewFlightRepository(db *mongo.Database) *FlightRepository {
	return &FlightRepository{
		collection: db.Collection("flights"),
	}
}

// Create creates a new flight
func (r *FlightRepository) Create(ctx context.Context, flight *admin.Flight) error {
	flight.CreatedAt = time.Now()
	flight.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, flight)
	return err
}

// FindByID finds a flight by its internal ID
func (r *FlightRepository) FindByID(ctx context.Context, id string) (*admin.Flight, error) {
	var flight admin.Flight
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&flight)
	if err != nil {
		return nil, err
	}
	return &flight, nil
}

// FindByFlightID finds a flight by its flight ID (for API lookup)
func (r *FlightRepository) FindByFlightID(ctx context.Context, flightID string) (*admin.Flight, error) {
	var flight admin.Flight
	err := r.collection.FindOne(ctx, bson.M{"flight_id": flightID}).Decode(&flight)
	if err != nil {
		return nil, err
	}
	return &flight, nil
}

// FindByFlightNumber finds a flight by its flight number
func (r *FlightRepository) FindByFlightNumber(ctx context.Context, flightNumber string) (*admin.Flight, error) {
	var flight admin.Flight
	err := r.collection.FindOne(ctx, bson.M{"flight_number": flightNumber}).Decode(&flight)
	if err != nil {
		return nil, err
	}
	return &flight, nil
}

// FindByStatus finds all flights with a specific status
func (r *FlightRepository) FindByStatus(ctx context.Context, status admin.FlightStatus) ([]admin.Flight, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var flights []admin.Flight
	if err := cursor.All(ctx, &flights); err != nil {
		return nil, err
	}
	return flights, nil
}

// FindAll finds all flights with optional pagination
func (r *FlightRepository) FindAll(ctx context.Context, limit, skip int64) ([]admin.Flight, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var flights []admin.Flight
	if err := cursor.All(ctx, &flights); err != nil {
		return nil, err
	}
	return flights, nil
}

// Update updates a flight
func (r *FlightRepository) Update(ctx context.Context, id string, flight *admin.Flight) error {
	flight.UpdatedAt = time.Now()
	update := bson.M{"$set": flight}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateStatus updates the status of a flight
func (r *FlightRepository) UpdateStatus(ctx context.Context, id string, status admin.FlightStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateCost updates the cost of a flight
func (r *FlightRepository) UpdateCost(ctx context.Context, id string, cost float64) error {
	update := bson.M{
		"$set": bson.M{
			"cost":       cost,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a flight
func (r *FlightRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of flights
func (r *FlightRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByStatus returns the number of flights by status
func (r *FlightRepository) CountByStatus(ctx context.Context, status admin.FlightStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}
