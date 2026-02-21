package app

import (
	"context"
	"time"

	"Anytime_Travel/backend/internal/models/app"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransferBookingRepository handles database operations for transfer bookings
type TransferBookingRepository struct {
	collection *mongo.Collection
}

// NewTransferBookingRepository creates a new transfer booking repository
func NewTransferBookingRepository(db *mongo.Database) *TransferBookingRepository {
	return &TransferBookingRepository{
		collection: db.Collection("transfer_bookings"),
	}
}

// Create creates a new transfer booking
func (r *TransferBookingRepository) Create(ctx context.Context, booking *app.TransferBooking) error {
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, booking)
	return err
}

// FindByID finds a transfer booking by its internal ID
func (r *TransferBookingRepository) FindByID(ctx context.Context, id string) (*app.TransferBooking, error) {
	var booking app.TransferBooking
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByBookingID finds a transfer booking by its user-facing booking ID
func (r *TransferBookingRepository) FindByBookingID(ctx context.Context, bookingID string) (*app.TransferBooking, error) {
	var booking app.TransferBooking
	err := r.collection.FindOne(ctx, bson.M{"booking_id": bookingID}).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByUserID finds all transfer bookings for a specific user
func (r *TransferBookingRepository) FindByUserID(ctx context.Context, userID string) ([]app.TransferBooking, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.TransferBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByTransferID finds all bookings for a specific transfer
func (r *TransferBookingRepository) FindByTransferID(ctx context.Context, transferID string) ([]app.TransferBooking, error) {
	filter := bson.M{"transfer_id": transferID}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.TransferBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByStatus finds all transfer bookings with a specific status
func (r *TransferBookingRepository) FindByStatus(ctx context.Context, status app.TransferBookingStatus) ([]app.TransferBooking, error) {
	filter := bson.M{"status": status}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.TransferBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindByPaymentStatus finds all transfer bookings with a specific payment status
func (r *TransferBookingRepository) FindByPaymentStatus(ctx context.Context, paymentStatus app.TransferPaymentStatus) ([]app.TransferBooking, error) {
	filter := bson.M{"payment_status": paymentStatus}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.TransferBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// FindAll finds all transfer bookings with optional pagination
func (r *TransferBookingRepository) FindAll(ctx context.Context, limit, skip int64) ([]app.TransferBooking, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit).
		SetSkip(skip)

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []app.TransferBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

// Update updates a transfer booking
func (r *TransferBookingRepository) Update(ctx context.Context, id string, booking *app.TransferBooking) error {
	booking.UpdatedAt = time.Now()
	update := bson.M{"$set": booking}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateStatus updates the status of a transfer booking
func (r *TransferBookingRepository) UpdateStatus(ctx context.Context, id string, status app.TransferBookingStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdatePaymentStatus updates the payment status of a transfer booking
func (r *TransferBookingRepository) UpdatePaymentStatus(ctx context.Context, id string, paymentStatus app.TransferPaymentStatus) error {
	update := bson.M{
		"$set": bson.M{
			"payment_status": paymentStatus,
			"updated_at":     time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a transfer booking
func (r *TransferBookingRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of transfer bookings
func (r *TransferBookingRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByUser returns the number of transfer bookings for a specific user
func (r *TransferBookingRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
}

// CountByStatus returns the number of transfer bookings by status
func (r *TransferBookingRepository) CountByStatus(ctx context.Context, status app.TransferBookingStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// CountBetween returns the number of transfer bookings created between start (inclusive) and end (exclusive).
func (r *TransferBookingRepository) CountBetween(ctx context.Context, start, end time.Time) (int64, error) {
	filter := bson.M{"created_at": bson.M{"$gte": start, "$lt": end}}
	return r.collection.CountDocuments(ctx, filter)
}

// UpdateCustomerEmail updates the customer email for a transfer booking
func (r *TransferBookingRepository) UpdateCustomerEmail(ctx context.Context, id string, email string) error {
	update := bson.M{
		"$set": bson.M{
			"customer.email": email,
			"updated_at":     time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateRefundRequested updates the refund requested status for a transfer booking
func (r *TransferBookingRepository) UpdateRefundRequested(ctx context.Context, id string, requested bool) error {
	update := bson.M{
		"$set": bson.M{
			"refund_requested": requested,
			"updated_at":       time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
