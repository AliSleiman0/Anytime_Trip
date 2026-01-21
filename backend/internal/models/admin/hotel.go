package admin

import "time"

// HotelStatus represents the status of a hotel
type HotelStatus string

const (
	HotelStatusActive      HotelStatus = "active"
	HotelStatusInactive    HotelStatus = "inactive"
	HotelStatusMaintenance HotelStatus = "maintenance"
	HotelStatusFullyBooked HotelStatus = "fully_booked"
)

// Hotel represents a hotel with minimal stored data
// Full hotel details are fetched from external API using HotelID
type Hotel struct {
	ID        string      `json:"id" bson:"_id"`
	HotelID   string      `json:"hotel_id" bson:"hotel_id"`     // Hotel identifier for API lookup
	HotelName string      `json:"hotel_name" bson:"hotel_name"` // Display name
	Status    HotelStatus `json:"status" bson:"status"`
	Cost      float64     `json:"cost" bson:"cost"`         // Room cost per night
	Currency  string      `json:"currency" bson:"currency"` // e.g., "USD", "$"
	// Provider metadata for admin service-provider table
	ProviderName    string    `json:"provider_name" bson:"provider_name"`
	ProviderType    string    `json:"type" bson:"type"` // e.g., "Hotel"
	Location        string    `json:"location" bson:"location"`
	ContactEmail    string    `json:"contact_email" bson:"contact_email"`
	Rating          float64   `json:"rating" bson:"rating"`
	TotalBookings   int64     `json:"total_bookings" bson:"total_bookings"`
	Revenue         float64   `json:"revenue" bson:"revenue"`
	Profit          float64   `json:"profit" bson:"profit"`
	ProfitPercent   float64   `json:"profit_percent" bson:"profit_percent"`
	PropertyDetails string    `json:"property_details" bson:"property_details"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
	IsFreezed       bool      `json:"is_freezed" bson:"is_freezed"`
}
