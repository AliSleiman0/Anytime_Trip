package admin

import (
	"context"
	"time"

	"travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationPreferencesRepository handles database operations for notification preferences
type NotificationPreferencesRepository struct {
	collection *mongo.Collection
}

// NewNotificationPreferencesRepository creates a new notification preferences repository
func NewNotificationPreferencesRepository(db *mongo.Database) *NotificationPreferencesRepository {
	return &NotificationPreferencesRepository{
		collection: db.Collection("notification_preferences"),
	}
}

// FindByAdminID finds notification preferences for an admin
func (r *NotificationPreferencesRepository) FindByAdminID(ctx context.Context, adminID string) (*admin.NotificationPreferences, error) {
	var prefs admin.NotificationPreferences
	err := r.collection.FindOne(ctx, bson.M{"admin_id": adminID}).Decode(&prefs)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return default preferences if not found
			defaultPrefs := admin.GetDefaultNotificationPreferences(adminID)
			return &defaultPrefs, nil
		}
		return nil, err
	}
	return &prefs, nil
}

// Upsert creates or updates notification preferences for an admin
func (r *NotificationPreferencesRepository) Upsert(ctx context.Context, prefs *admin.NotificationPreferences) error {
	prefs.UpdatedAt = time.Now()

	filter := bson.M{"admin_id": prefs.AdminID}
	update := bson.M{
		"$set": bson.M{
			"admin_id":   prefs.AdminID,
			"bookings":   prefs.Bookings,
			"payments":   prefs.Payments,
			"updated_at": prefs.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// UpdateBookingPreferences updates only booking notification preferences
func (r *NotificationPreferencesRepository) UpdateBookingPreferences(ctx context.Context, adminID string, bookings admin.BookingNotificationSettings) error {
	filter := bson.M{"admin_id": adminID}
	update := bson.M{
		"$set": bson.M{
			"bookings":   bookings,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdatePaymentPreferences updates only payment notification preferences
func (r *NotificationPreferencesRepository) UpdatePaymentPreferences(ctx context.Context, adminID string, payments admin.PaymentNotificationSettings) error {
	filter := bson.M{"admin_id": adminID}
	update := bson.M{
		"$set": bson.M{
			"payments":   payments,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes notification preferences for an admin
func (r *NotificationPreferencesRepository) Delete(ctx context.Context, adminID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"admin_id": adminID})
	return err
}
