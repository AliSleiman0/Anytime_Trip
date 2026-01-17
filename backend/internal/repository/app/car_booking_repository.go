package app

import (
	"context"
	"time"

	"Anytime_Travel/backend/internal/models/app"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CarBookingRepository handles database operations for car bookings
type CarBookingRepository struct {
	collection *mongo.Collection
}

// NewCarBookingRepository creates a new car booking repository
func NewCarBookingRepository(db *mongo.Database) *CarBookingRepository {
	return &CarBookingRepository{
		collection: db.Collection("car_bookings"),
	}
}

// Create creates a new car booking
func (r *CarBookingRepository) Create(ctx context.Context, booking *app.CarBooking) error {
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, booking)
	return err
}

// FindByID finds a car booking by its internal ID
func (r *CarBookingRepository) FindByID(ctx context.Context, id string) (*app.CarBooking, error) {
	var booking app.CarBooking
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByBookingID finds a car booking by its user-facing booking ID
func (r *CarBookingRepository) FindByBookingID(ctx context.Context, bookingID string) (*app.CarBooking, error) {
	var booking app.CarBooking
	err := r.collection.FindOne(ctx, bson.M{"booking_id": bookingID}).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByUserID finds all car bookings for a specific user
func (r *CarBookingRepository) FindByUserID(ctx context.Context, userID string) ([]app.CarBooking, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.CarBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByCarID finds all bookings for a specific car
func (r *CarBookingRepository) FindByCarID(ctx context.Context, carID string) ([]app.CarBooking, error) {
	filter := bson.M{"car_id": carID}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.CarBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByStatus finds all car bookings with a specific status
func (r *CarBookingRepository) FindByStatus(ctx context.Context, status app.CarBookingStatus) ([]app.CarBooking, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.CarBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByPaymentStatus finds all car bookings with a specific payment status
func (r *CarBookingRepository) FindByPaymentStatus(ctx context.Context, paymentStatus app.CarPaymentStatus) ([]app.CarBooking, error) {
	filter := bson.M{"payment_status": paymentStatus}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.CarBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindAll finds all car bookings with optional pagination
func (r *CarBookingRepository) FindAll(ctx context.Context, limit, skip int64) ([]app.CarBooking, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.CarBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// Update updates a car booking
func (r *CarBookingRepository) Update(ctx context.Context, id string, booking *app.CarBooking) error {
	booking.UpdatedAt = time.Now()
	update := bson.M{"$set": booking}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateStatus updates the status of a car booking
func (r *CarBookingRepository) UpdateStatus(ctx context.Context, id string, status app.CarBookingStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdatePaymentStatus updates the payment status of a car booking
func (r *CarBookingRepository) UpdatePaymentStatus(ctx context.Context, id string, paymentStatus app.CarPaymentStatus) error {
	update := bson.M{
		"$set": bson.M{
			"payment_status": paymentStatus,
			"updated_at":     time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a car booking
func (r *CarBookingRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of car bookings
func (r *CarBookingRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByUser returns the number of car bookings for a specific user
func (r *CarBookingRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
}

// CountByStatus returns the number of car bookings by status
func (r *CarBookingRepository) CountByStatus(ctx context.Context, status app.CarBookingStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// UpdateCustomerEmail updates the customer email for a car booking
func (r *CarBookingRepository) UpdateCustomerEmail(ctx context.Context, id string, email string) error {
	update := bson.M{
		"$set": bson.M{
			"customer.email": email,
			"updated_at":     time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
