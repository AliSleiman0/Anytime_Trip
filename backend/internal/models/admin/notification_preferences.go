package admin

import "time"

// BookingNotificationSettings represents booking-related notification preferences
type BookingNotificationSettings struct {
	NewBooking          bool `json:"new_booking" bson:"new_booking"`
	BookingCancelled    bool `json:"booking_cancelled" bson:"booking_cancelled"`
	BookingConfirmation bool `json:"booking_confirmation" bson:"booking_confirmation"`
	BookingRescheduled  bool `json:"booking_rescheduled" bson:"booking_rescheduled"`
	BookingReminder     bool `json:"booking_reminder" bson:"booking_reminder"`
	BookingUpdated      bool `json:"booking_updated" bson:"booking_updated"`
}

// PaymentNotificationSettings represents payment-related notification preferences
type PaymentNotificationSettings struct {
	PaymentsReceived   bool `json:"payments_received" bson:"payments_received"`
	PaymentsConfirmed  bool `json:"payments_confirmed" bson:"payments_confirmed"`
	PaymentReceived    bool `json:"payment_received" bson:"payment_received"`
	BookingRescheduled bool `json:"booking_rescheduled" bson:"booking_rescheduled"`
	PaymentPending     bool `json:"payment_pending" bson:"payment_pending"`
}

// NotificationPreferences contains all notification settings for an admin
// Stored in a separate collection and linked to admin by AdminID
type NotificationPreferences struct {
	AdminID   string                      `json:"admin_id" bson:"admin_id"`
	Bookings  BookingNotificationSettings `json:"bookings" bson:"bookings"`
	Payments  PaymentNotificationSettings `json:"payments" bson:"payments"`
	CreatedAt time.Time                   `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time                   `json:"updated_at" bson:"updated_at"`
}

// GetDefaultNotificationPreferences returns default notification settings for a new admin
func GetDefaultNotificationPreferences(adminID string) NotificationPreferences {
	now := time.Now()
	return NotificationPreferences{
		AdminID: adminID,
		Bookings: BookingNotificationSettings{
			NewBooking:          false,
			BookingCancelled:    true,
			BookingConfirmation: false,
			BookingRescheduled:  false,
			BookingReminder:     true,
			BookingUpdated:      false,
		},
		Payments: PaymentNotificationSettings{
			PaymentsReceived:   false,
			PaymentsConfirmed:  true,
			PaymentReceived:    false,
			BookingRescheduled: false,
			PaymentPending:     true,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
