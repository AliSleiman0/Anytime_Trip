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
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phone_number" bson:"phone_number"`
}

// TransferVehicle contains vehicle information
type TransferVehicle struct {
	Type       string  `json:"type" bson:"type"`             // e.g., "Midsize SUV"
	Model      string  `json:"model" bson:"model"`           // e.g., "Toyota RAV 4 or similar"
	Passengers int     `json:"passengers" bson:"passengers"` // Number of passengers
	Duration   string  `json:"duration" bson:"duration"`     // e.g., "40 min"
	MeetGreet  bool    `json:"meet_greet" bson:"meet_greet"`
	Image      string  `json:"image" bson:"image"` // Vehicle image URL
	Price      float64 `json:"price" bson:"price"` // Base price for display
}

// TransferLocation contains location and time details
type TransferLocation struct {
	Name    string    `json:"name" bson:"name"`       // e.g., "BEY Airport", "Beirut Hotel"
	Address string    `json:"address" bson:"address"` // Full address
	Date    time.Time `json:"date" bson:"date"`       // Pickup/Dropoff date
	Time    string    `json:"time" bson:"time"`       // e.g., "1:30 PM"
}

// TransferPricing contains pricing details
type TransferPricing struct {
	BasePrice float64 `json:"base_price" bson:"base_price"`   // Base transfer price
	Taxes     float64 `json:"taxes" bson:"taxes"`             // Taxes
	TaxOnFees float64 `json:"tax_on_fees" bson:"tax_on_fees"` // Tax on fees
	Total     float64 `json:"total" bson:"total"`             // Total price
}

// TransferBooking represents a user's transfer booking
type TransferBooking struct {
	ID              string                `json:"id" bson:"_id"`
	BookingID       string                `json:"booking_id" bson:"booking_id"` // e.g., "TBK-12345678"
	UserID          string                `json:"user_id" bson:"user_id"`
	TransferID      string                `json:"transfer_id" bson:"transfer_id"` // Reference to transfer in admin collection
	Status          TransferBookingStatus `json:"status" bson:"status"`
	Customer        TransferCustomer      `json:"customer" bson:"customer"`
	Vehicle         TransferVehicle       `json:"vehicle" bson:"vehicle"`
	Pickup          TransferLocation      `json:"pickup" bson:"pickup"`
	Dropoff         TransferLocation      `json:"dropoff" bson:"dropoff"`
	Pricing         TransferPricing       `json:"pricing" bson:"pricing"`
	Details         string                `json:"details" bson:"details"` // e.g., "Airport pickup - Terminal 3"
	BookingDate     time.Time             `json:"booking_date" bson:"booking_date"`
	Amount          float64               `json:"amount" bson:"amount"`
	Currency        string                `json:"currency" bson:"currency"`
	PaymentStatus   TransferPaymentStatus `json:"payment_status" bson:"payment_status"`
	RefundRequested bool                  `json:"refund_requested" bson:"refund_requested"`
	RefundAmount    float64               `json:"refund_amount,omitempty" bson:"refund_amount,omitempty"`
	CreatedAt       time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at" bson:"updated_at"`
}
