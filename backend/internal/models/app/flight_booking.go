package app

import "time"

// FlightPaymentStatus represents the payment status of a flight booking
type FlightPaymentStatus string

const (
	FlightPaymentStatusPending   FlightPaymentStatus = "pending"
	FlightPaymentStatusPaid      FlightPaymentStatus = "paid"
	FlightPaymentStatusByWhish   FlightPaymentStatus = "by_whish"
	FlightPaymentStatusRefunded  FlightPaymentStatus = "refunded"
	FlightPaymentStatusCancelled FlightPaymentStatus = "cancelled"
)

// FlightBookingStatus represents the status of a flight booking
type FlightBookingStatus string

const (
	FlightBookingStatusConfirmed FlightBookingStatus = "confirmed"
	FlightBookingStatusPending   FlightBookingStatus = "pending"
	FlightBookingStatusCancelled FlightBookingStatus = "cancelled"
	FlightBookingStatusCompleted FlightBookingStatus = "completed"
)

// FlightCustomer contains customer information for a flight booking
type FlightCustomer struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

// FlightBooking represents a user's flight booking
type FlightBooking struct {
	ID            string              `json:"id" bson:"_id"`
	BookingID     string              `json:"booking_id" bson:"booking_id"` // e.g., "BK001"
	UserID        string              `json:"user_id" bson:"user_id"`
	FlightID      string              `json:"flight_id" bson:"flight_id"` // Reference to flight in admin collection
	Status        FlightBookingStatus `json:"status" bson:"status"`
	Customer      FlightCustomer      `json:"customer" bson:"customer"`
	Details       string              `json:"details" bson:"details"` // e.g., "1DLRT190/6MARSAN-SMITH"
	BookingDate   time.Time           `json:"booking_date" bson:"booking_date"`
	Amount        float64             `json:"amount" bson:"amount"`
	Currency      string              `json:"currency" bson:"currency"`
	PaymentStatus FlightPaymentStatus `json:"payment_status" bson:"payment_status"`
	CreatedAt     time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at" bson:"updated_at"`
}
