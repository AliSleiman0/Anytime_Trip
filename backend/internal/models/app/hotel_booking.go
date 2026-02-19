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

// HotelGuest represents a guest staying at the hotel
type HotelGuest struct {
	Name string `json:"name" bson:"name"`
	Type string `json:"type" bson:"type"` // e.g., "Adult", "Child"
}

// HotelPricing contains pricing details for hotel booking
type HotelPricing struct {
	NightPrice     float64 `json:"night_price" bson:"night_price"`         // Price per night
	Taxes          float64 `json:"taxes" bson:"taxes"`                     // Taxes
	DestinationFee float64 `json:"destination_fee" bson:"destination_fee"` // Destination fee
	ServiceFee     float64 `json:"service_fee" bson:"service_fee"`         // Service fee
	Total          float64 `json:"total" bson:"total"`                     // Total price
}

// HotelBooking represents a user's hotel booking
type HotelBooking struct {
	ID            string             `json:"id" bson:"_id"`
	BookingID     string             `json:"booking_id" bson:"booking_id"` // e.g., "HBK-12345678"
	UserID        string             `json:"user_id" bson:"user_id"`
	HotelID       string             `json:"hotel_id" bson:"hotel_id"` // Reference to hotel in admin collection
	Status        HotelBookingStatus `json:"status" bson:"status"`
	Customer      HotelCustomer      `json:"customer" bson:"customer"`
	Guests        []HotelGuest       `json:"guests" bson:"guests"`           // List of guests
	HotelName     string             `json:"hotel_name" bson:"hotel_name"`   // Hotel name
	HotelImage    string             `json:"hotel_image" bson:"hotel_image"` // Hotel image URL
	CheckInDate   time.Time          `json:"check_in_date" bson:"check_in_date"`
	CheckOutDate  time.Time          `json:"check_out_date" bson:"check_out_date"`
	Nights        int                `json:"nights" bson:"nights"`   // Number of nights
	Rooms         int                `json:"rooms" bson:"rooms"`     // Number of rooms
	Pricing       HotelPricing       `json:"pricing" bson:"pricing"` // Pricing details
	Details       string             `json:"details,omitempty" bson:"details,omitempty"`
	BookingDate   time.Time          `json:"booking_date" bson:"booking_date"`
	Amount        float64            `json:"amount" bson:"amount"`
	Currency      string             `json:"currency" bson:"currency"`
	PaymentStatus HotelPaymentStatus `json:"payment_status" bson:"payment_status"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}
