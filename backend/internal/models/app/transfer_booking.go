package app

import "time"

// TransferPaymentStatus represents the payment status of a transfer booking
type TransferPaymentStatus string

const (
	TransferPaymentStatusPending   TransferPaymentStatus = "pending"
	TransferPaymentStatusPaid      TransferPaymentStatus = "paid"
	TransferPaymentStatusByWhish   TransferPaymentStatus = "by_whish"
	TransferPaymentStatusRefunded  TransferPaymentStatus = "refunded"
	TransferPaymentStatusCancelled TransferPaymentStatus = "cancelled"
)

// TransferBookingStatus represents the status of a transfer booking
type TransferBookingStatus string

const (
	TransferBookingStatusConfirmed TransferBookingStatus = "confirmed"
	TransferBookingStatusPending   TransferBookingStatus = "pending"
	TransferBookingStatusCancelled TransferBookingStatus = "cancelled"
	TransferBookingStatusCompleted TransferBookingStatus = "completed"
)

// TransferCustomer contains customer information for a transfer booking
type TransferCustomer struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

// TransferBooking represents a user's transfer booking
type TransferBooking struct {
	ID            string                `json:"id" bson:"_id"`
	BookingID     string                `json:"booking_id" bson:"booking_id"` // e.g., "BK001"
	UserID        string                `json:"user_id" bson:"user_id"`
	TransferID    string                `json:"transfer_id" bson:"transfer_id"` // Reference to transfer in admin collection
	Status        TransferBookingStatus `json:"status" bson:"status"`
	Customer      TransferCustomer      `json:"customer" bson:"customer"`
	Details       string                `json:"details" bson:"details"` // e.g., "Airport pickup - Terminal 3"
	BookingDate   time.Time             `json:"booking_date" bson:"booking_date"`
	Amount        float64               `json:"amount" bson:"amount"`
	Currency      string                `json:"currency" bson:"currency"`
	PaymentStatus TransferPaymentStatus `json:"payment_status" bson:"payment_status"`
	CreatedAt     time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at" bson:"updated_at"`
}
