package app

import "time"

// CarPaymentStatus represents the payment status of a car booking
type CarPaymentStatus string

const (
	CarPaymentStatusPending   CarPaymentStatus = "pending"
	CarPaymentStatusPaid      CarPaymentStatus = "paid"
	CarPaymentStatusByWhish   CarPaymentStatus = "by_whish"
	CarPaymentStatusRefunded  CarPaymentStatus = "refunded"
	CarPaymentStatusCancelled CarPaymentStatus = "cancelled"
)

// CarBookingStatus represents the status of a car booking
type CarBookingStatus string

const (
	CarBookingStatusConfirmed CarBookingStatus = "confirmed"
	CarBookingStatusPending   CarBookingStatus = "pending"
	CarBookingStatusCancelled CarBookingStatus = "cancelled"
	CarBookingStatusCompleted CarBookingStatus = "completed"
)

// CarCustomer contains customer information for a car booking
type CarCustomer struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

// CarBooking represents a user's car rental booking
type CarBooking struct {
	ID            string           `json:"id" bson:"_id"`
	BookingID     string           `json:"booking_id" bson:"booking_id"` // e.g., "BK001"
	UserID        string           `json:"user_id" bson:"user_id"`
	CarID         string           `json:"car_id" bson:"car_id"` // Reference to car in admin collection
	Status        CarBookingStatus `json:"status" bson:"status"`
	Customer      CarCustomer      `json:"customer" bson:"customer"`
	Details       string           `json:"details,omitempty" bson:"details,omitempty"`
	BookingDate   time.Time        `json:"booking_date" bson:"booking_date"`
	Amount        float64          `json:"amount" bson:"amount"`
	Currency      string           `json:"currency" bson:"currency"`
	PaymentStatus CarPaymentStatus `json:"payment_status" bson:"payment_status"`
	CreatedAt     time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" bson:"updated_at"`
}
