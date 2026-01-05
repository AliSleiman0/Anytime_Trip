package app

import "time"

// HotelPaymentStatus represents the payment status of a hotel booking
type HotelPaymentStatus string

const (
	HotelPaymentStatusPending   HotelPaymentStatus = "pending"
	HotelPaymentStatusPaid      HotelPaymentStatus = "paid"
	HotelPaymentStatusByWhish   HotelPaymentStatus = "by_whish"
	HotelPaymentStatusRefunded  HotelPaymentStatus = "refunded"
	HotelPaymentStatusCancelled HotelPaymentStatus = "cancelled"
)

// HotelBookingStatus represents the status of a hotel booking
type HotelBookingStatus string

const (
	HotelBookingStatusConfirmed HotelBookingStatus = "confirmed"
	HotelBookingStatusPending   HotelBookingStatus = "pending"
	HotelBookingStatusCancelled HotelBookingStatus = "cancelled"
	HotelBookingStatusCompleted HotelBookingStatus = "completed"
)

// HotelCustomer contains customer information for a hotel booking
type HotelCustomer struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

// HotelBooking represents a user's hotel booking
type HotelBooking struct {
	ID            string             `json:"id" bson:"_id"`
	BookingID     string             `json:"booking_id" bson:"booking_id"` // e.g., "BK001"
	UserID        string             `json:"user_id" bson:"user_id"`
	HotelID       string             `json:"hotel_id" bson:"hotel_id"` // Reference to hotel in admin collection
	Status        HotelBookingStatus `json:"status" bson:"status"`
	Customer      HotelCustomer      `json:"customer" bson:"customer"`
	Details       string             `json:"details,omitempty" bson:"details,omitempty"`
	BookingDate   time.Time          `json:"booking_date" bson:"booking_date"`
	Amount        float64            `json:"amount" bson:"amount"`
	Currency      string             `json:"currency" bson:"currency"`
	PaymentStatus HotelPaymentStatus `json:"payment_status" bson:"payment_status"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}
