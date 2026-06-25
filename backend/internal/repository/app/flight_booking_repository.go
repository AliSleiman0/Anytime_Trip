package app

import (
	"context"
	"time"

	"travel/backend/internal/models/app"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FlightBookingRepository handles database operations for flight bookings
type FlightBookingRepository struct {
	collection *mongo.Collection
}

// NewFlightBookingRepository creates a new flight booking repository
func NewFlightBookingRepository(db *mongo.Database) *FlightBookingRepository {
	return &FlightBookingRepository{
		collection: db.Collection("flight_bookings"),
	}
}

// Create creates a new flight booking
func (r *FlightBookingRepository) Create(ctx context.Context, booking *app.FlightBooking) error {
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, booking)
	return err
}

// FindByID finds a flight booking by its internal ID
func (r *FlightBookingRepository) FindByID(ctx context.Context, id string) (*app.FlightBooking, error) {
	var booking app.FlightBooking
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByBookingID finds a flight booking by its user-facing booking ID
func (r *FlightBookingRepository) FindByBookingID(ctx context.Context, bookingID string) (*app.FlightBooking, error) {
	var booking app.FlightBooking
	err := r.collection.FindOne(ctx, bson.M{"booking_id": bookingID}).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByUserID finds all flight bookings for a specific user
func (r *FlightBookingRepository) FindByUserID(ctx context.Context, userID string) ([]app.FlightBooking, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.FlightBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByFlightID finds all bookings for a specific flight
func (r *FlightBookingRepository) FindByFlightID(ctx context.Context, flightID string) ([]app.FlightBooking, error) {
	filter := bson.M{"flight_id": flightID}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.FlightBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByStatus finds all flight bookings with a specific status
func (r *FlightBookingRepository) FindByStatus(ctx context.Context, status app.FlightBookingStatus) ([]app.FlightBooking, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.FlightBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByPaymentStatus finds all flight bookings with a specific payment status
func (r *FlightBookingRepository) FindByPaymentStatus(ctx context.Context, paymentStatus app.FlightPaymentStatus) ([]app.FlightBooking, error) {
	filter := bson.M{"payment_status": paymentStatus}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.FlightBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindAll finds all flight bookings with optional pagination
func (r *FlightBookingRepository) FindAll(ctx context.Context, limit, skip int64) ([]app.FlightBooking, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.FlightBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// Update updates a flight booking
func (r *FlightBookingRepository) Update(ctx context.Context, id string, booking *app.FlightBooking) error {
	booking.UpdatedAt = time.Now()
	update := bson.M{"$set": booking}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateStatus updates the status of a flight booking
func (r *FlightBookingRepository) UpdateStatus(ctx context.Context, id string, status app.FlightBookingStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdatePaymentStatus updates the payment status of a flight booking
func (r *FlightBookingRepository) UpdatePaymentStatus(ctx context.Context, id string, paymentStatus app.FlightPaymentStatus) error {
	update := bson.M{
		"$set": bson.M{
			"payment_status": paymentStatus,
			"updated_at":     time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a flight booking
func (r *FlightBookingRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of flight bookings
func (r *FlightBookingRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByUser returns the number of flight bookings for a specific user
func (r *FlightBookingRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
}

// CountByStatus returns the number of flight bookings by status
func (r *FlightBookingRepository) CountByStatus(ctx context.Context, status app.FlightBookingStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// CountBetween returns the number of flight bookings created between start (inclusive) and end (exclusive).
func (r *FlightBookingRepository) CountBetween(ctx context.Context, start, end time.Time) (int64, error) {
	filter := bson.M{"created_at": bson.M{"$gte": start, "$lt": end}}
	return r.collection.CountDocuments(ctx, filter)
}

// UpdateCustomerEmail updates the customer email for a flight booking
func (r *FlightBookingRepository) UpdateCustomerEmail(ctx context.Context, id string, email string) error {
	update := bson.M{
		"$set": bson.M{
			"customer.email": email,
			"updated_at":     time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateRefundAmount updates the refund amount for a flight booking
func (r *FlightBookingRepository) UpdateRefundAmount(ctx context.Context, id string, amount float64) error {
	update := bson.M{
		"$set": bson.M{
			"refund_amount": amount,
			"updated_at":    time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
