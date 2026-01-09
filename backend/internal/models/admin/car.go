package admin

import "time"

// CarStatus represents the status of a car rental
type CarStatus string

const (
	CarStatusActive      CarStatus = "active"
	CarStatusInactive    CarStatus = "inactive"
	CarStatusMaintenance CarStatus = "maintenance"
	CarStatusRented      CarStatus = "rented"
)

// Car represents a car rental with minimal stored data
// Full car details are fetched from external API using CarID
type Car struct {
	ID       string    `json:"id" bson:"_id"`
	CarID    string    `json:"car_id" bson:"car_id"`     // Car identifier for API lookup
	CarName  string    `json:"car_name" bson:"car_name"` // Display name like "Toyota RAV4"
	Status   CarStatus `json:"status" bson:"status"`
	Cost     float64   `json:"cost" bson:"cost"`         // Rental cost
	Currency string    `json:"currency" bson:"currency"` // e.g., "USD", "$"
	// Provider metadata for admin service-provider table
	ProviderName  string    `json:"provider_name" bson:"provider_name"`
	ProviderType  string    `json:"type" bson:"type"` // e.g., "Car Rental"
	Location      string    `json:"location" bson:"location"`
	ContactEmail  string    `json:"contact_email" bson:"contact_email"`
	Rating        float64   `json:"rating" bson:"rating"`
	TotalBookings int64     `json:"total_bookings" bson:"total_bookings"`
	Revenue       float64   `json:"revenue" bson:"revenue"`
	Profit        float64   `json:"profit" bson:"profit"`
	ProfitPercent float64   `json:"profit_percent" bson:"profit_percent"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}
