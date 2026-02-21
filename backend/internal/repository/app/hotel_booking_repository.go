package app

import (
	"context"
	"time"

	"Anytime_Travel/backend/internal/models/app"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// HotelBookingRepository handles database operations for hotel bookings
type HotelBookingRepository struct {
	collection *mongo.Collection
}

// NewHotelBookingRepository creates a new hotel booking repository
func NewHotelBookingRepository(db *mongo.Database) *HotelBookingRepository {
	return &HotelBookingRepository{
		collection: db.Collection("hotel_bookings"),
	}
}

// Create creates a new hotel booking
func (r *HotelBookingRepository) Create(ctx context.Context, booking *app.HotelBooking) error {
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, booking)
	return err
}

// FindByID finds a hotel booking by its internal ID
func (r *HotelBookingRepository) FindByID(ctx context.Context, id string) (*app.HotelBooking, error) {
	var booking app.HotelBooking
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByBookingID finds a hotel booking by its user-facing booking ID
func (r *HotelBookingRepository) FindByBookingID(ctx context.Context, bookingID string) (*app.HotelBooking, error) {
	var booking app.HotelBooking
	err := r.collection.FindOne(ctx, bson.M{"booking_id": bookingID}).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByUserID finds all hotel bookings for a specific user
func (r *HotelBookingRepository) FindByUserID(ctx context.Context, userID string) ([]app.HotelBooking, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.HotelBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByHotelID finds all bookings for a specific hotel
func (r *HotelBookingRepository) FindByHotelID(ctx context.Context, hotelID string) ([]app.HotelBooking, error) {
	filter := bson.M{"hotel_id": hotelID}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.HotelBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByStatus finds all hotel bookings with a specific status
func (r *HotelBookingRepository) FindByStatus(ctx context.Context, status app.HotelBookingStatus) ([]app.HotelBooking, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.HotelBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByPaymentStatus finds all hotel bookings with a specific payment status
func (r *HotelBookingRepository) FindByPaymentStatus(ctx context.Context, paymentStatus app.HotelPaymentStatus) ([]app.HotelBooking, error) {
	filter := bson.M{"payment_status": paymentStatus}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.HotelBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindAll finds all hotel bookings with optional pagination
func (r *HotelBookingRepository) FindAll(ctx context.Context, limit, skip int64) ([]app.HotelBooking, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.HotelBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// Update updates a hotel booking
func (r *HotelBookingRepository) Update(ctx context.Context, id string, booking *app.HotelBooking) error {
	booking.UpdatedAt = time.Now()
	update := bson.M{"$set": booking}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateStatus updates the status of a hotel booking
func (r *HotelBookingRepository) UpdateStatus(ctx context.Context, id string, status app.HotelBookingStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdatePaymentStatus updates the payment status of a hotel booking
func (r *HotelBookingRepository) UpdatePaymentStatus(ctx context.Context, id string, paymentStatus app.HotelPaymentStatus) error {
	update := bson.M{
		"$set": bson.M{
			"payment_status": paymentStatus,
			"updated_at":     time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a hotel booking
func (r *HotelBookingRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of hotel bookings
func (r *HotelBookingRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByUser returns the number of hotel bookings for a specific user
func (r *HotelBookingRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
}

// CountByStatus returns the number of hotel bookings by status
func (r *HotelBookingRepository) CountByStatus(ctx context.Context, status app.HotelBookingStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// CountBetween returns the number of hotel bookings created between start (inclusive) and end (exclusive).
func (r *HotelBookingRepository) CountBetween(ctx context.Context, start, end time.Time) (int64, error) {
	filter := bson.M{"created_at": bson.M{"$gte": start, "$lt": end}}
	return r.collection.CountDocuments(ctx, filter)
}

// UpdateCustomerEmail updates the customer email for a hotel booking
func (r *HotelBookingRepository) UpdateCustomerEmail(ctx context.Context, id string, email string) error {
	update := bson.M{
		"$set": bson.M{
			"customer.email": email,
			"updated_at":     time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateRefundAmount updates the refund amount for a hotel booking
func (r *HotelBookingRepository) UpdateRefundAmount(ctx context.Context, id string, amount float64) error {
	update := bson.M{
		"$set": bson.M{
			"refund_amount": amount,
			"updated_at":    time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateRefundRequested updates the refund requested status for a hotel booking
func (r *HotelBookingRepository) UpdateRefundRequested(ctx context.Context, id string, requested bool) error {
	update := bson.M{
		"$set": bson.M{
			"refund_requested": requested,
			"updated_at":       time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
