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

// CarDriver contains driver information
type CarDriver struct {
	Name          string `json:"name" bson:"name"`
	PhoneNumber   string `json:"phone_number" bson:"phone_number"`
	LicenseNumber string `json:"license_number" bson:"license_number"`
}

// CarPickupDropoff contains location and time details
type CarPickupDropoff struct {
	Location string    `json:"location" bson:"location"` // e.g., "BEY Airport"
	Address  string    `json:"address" bson:"address"`   // Full address
	Date     time.Time `json:"date" bson:"date"`         // Pickup/Dropoff date
	Time     string    `json:"time" bson:"time"`         // e.g., "1:30 AM"
}

// CarPricing contains pricing details for car booking
type CarPricing struct {
	RentalPrice float64 `json:"rental_price" bson:"rental_price"` // Base rental price
	Taxes       float64 `json:"taxes" bson:"taxes"`               // Taxes
	Total       float64 `json:"total" bson:"total"`               // Total price
}

// CarBooking represents a user's car rental booking
type CarBooking struct {
	ID            string           `json:"id" bson:"_id"`
	BookingID     string           `json:"booking_id" bson:"booking_id"` // e.g., "CBK-12345678"
	UserID        string           `json:"user_id" bson:"user_id"`
	CarID         string           `json:"car_id" bson:"car_id"` // Reference to car in admin collection
	Status        CarBookingStatus `json:"status" bson:"status"`
	Customer      CarCustomer      `json:"customer" bson:"customer"`
	CarType       string           `json:"car_type" bson:"car_type"` // e.g., "SUV", "Sedan"
	Passengers    int              `json:"passengers" bson:"passengers"`
	Pickup        CarPickupDropoff `json:"pickup" bson:"pickup"`
	Dropoff       CarPickupDropoff `json:"dropoff" bson:"dropoff"`
	Driver        CarDriver        `json:"driver" bson:"driver"`
	Pricing       CarPricing       `json:"pricing" bson:"pricing"`
	Details       string           `json:"details,omitempty" bson:"details,omitempty"`
	BookingDate   time.Time        `json:"booking_date" bson:"booking_date"`
	Amount        float64          `json:"amount" bson:"amount"`
	Currency      string           `json:"currency" bson:"currency"`
	PaymentStatus CarPaymentStatus `json:"payment_status" bson:"payment_status"`
	CreatedAt     time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" bson:"updated_at"`
}
