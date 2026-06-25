package admin

import (
	"context"
	"time"

	"travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationRepository handles database operations for admin notification preferences
type NotificationRepository struct {
	collection *mongo.Collection
}

// NewNotificationRepository creates a new notification preferences repository
func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	return &NotificationRepository{
		collection: db.Collection("admin_notification_preferences"),
	}
}

// Create creates notification preferences for an admin
func (r *NotificationRepository) Create(ctx context.Context, preferences *admin.NotificationPreferences) error {
	_, err := r.collection.InsertOne(ctx, preferences)
	return err
}

// FindByAdminID retrieves notification preferences for a specific admin
func (r *NotificationRepository) FindByAdminID(ctx context.Context, adminID string) (*admin.NotificationPreferences, error) {
	var preferences admin.NotificationPreferences
	err := r.collection.FindOne(ctx, bson.M{"admin_id": adminID}).Decode(&preferences)
	if err != nil {
		return nil, err
	}
	return &preferences, nil
}

// Update updates all notification preferences for an admin
func (r *NotificationRepository) Update(ctx context.Context, adminID string, preferences admin.NotificationPreferences) error {
	preferences.UpdatedAt = time.Now()
	update := bson.M{
		"$set": preferences,
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, bson.M{"admin_id": adminID}, update, opts)
	return err
}

// UpdateBookingNotifications updates only booking notification settings
func (r *NotificationRepository) UpdateBookingNotifications(ctx context.Context, adminID string, settings admin.BookingNotificationSettings) error {
	update := bson.M{
		"$set": bson.M{
			"bookings":   settings,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"admin_id": adminID}, update)
	return err
}

// UpdatePaymentNotifications updates only payment notification settings
func (r *NotificationRepository) UpdatePaymentNotifications(ctx context.Context, adminID string, settings admin.PaymentNotificationSettings) error {
	update := bson.M{
		"$set": bson.M{
			"payments":   settings,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"admin_id": adminID}, update)
	return err
}

// ToggleNotificationSetting toggles a specific notification setting
// settingPath should be in dot notation, e.g., "bookings.new_booking" or "payments.payment_pending"
func (r *NotificationRepository) ToggleNotificationSetting(ctx context.Context, adminID string, settingPath string, value bool) error {
	update := bson.M{
		"$set": bson.M{
			settingPath:  value,
			"updated_at": time.Now(),
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"admin_id": adminID}, update)
	return err
}

// Delete removes notification preferences for an admin
func (r *NotificationRepository) Delete(ctx context.Context, adminID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"admin_id": adminID})
	return err
}

// CreateDefaultForAdmin creates default notification preferences for a new admin
func (r *NotificationRepository) CreateDefaultForAdmin(ctx context.Context, adminID string) error {
	defaults := admin.GetDefaultNotificationPreferences(adminID)
	return r.Create(ctx, &defaults)
}
