package app

import "time"

// User represents an app user
type User struct {
	ID            string    `json:"id" bson:"_id"`
	Name          string    `json:"name" bson:"name"`
	Email         string    `json:"email" bson:"email"`
	PhoneNumber   string    `json:"phone_number" bson:"phone_number"`
	IsActive      bool      `json:"is_active" bson:"is_active"`
	TotalBookings int       `json:"total_bookings" bson:"total_bookings"`
	LastBooking   time.Time `json:"last_booking" bson:"last_booking,omitempty"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	LastLogin     time.Time `json:"last_login" bson:"last_login,omitempty"` // New field to track last login
}
