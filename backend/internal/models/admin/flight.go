package admin

import "time"

// FlightStatus represents the status of a flight
type FlightStatus string

const (
	FlightStatusActive    FlightStatus = "active"
	FlightStatusInactive  FlightStatus = "inactive"
	FlightStatusCancelled FlightStatus = "cancelled"
	FlightStatusDelayed   FlightStatus = "delayed"
)

// Flight represents a flight with minimal stored data
// Full flight details are fetched from external API using FlightID
type Flight struct {
	ID                string       `json:"id" bson:"_id"`
	FlightID          string       `json:"flight_id" bson:"flight_id"`         // Flight identifier for API lookup
	FlightNumber      string       `json:"flight_number" bson:"flight_number"` // Display number like "ME 225"
	Status            FlightStatus `json:"status" bson:"status"`
	Cost              float64      `json:"cost" bson:"cost"`         // Top-up or base cost
	Currency          string       `json:"currency" bson:"currency"` // e.g., "USD", "$"
	DepartureTime     time.Time    `json:"departure_time" bson:"departure_time"`
	ArrivalTime       time.Time    `json:"arrival_time" bson:"arrival_time"`
	DepartureLocation string       `json:"departure_location" bson:"departure_location"`
	ArrivalLocation   string       `json:"arrival_location" bson:"arrival_location"`
	// Provider metadata for admin service-provider table
	ProviderName  string    `json:"provider_name" bson:"provider_name"`
	ProviderType  string    `json:"type" bson:"type"` // e.g., "Airline"
	Location      string    `json:"location" bson:"location"`
	ContactEmail  string    `json:"contact_email" bson:"contact_email"`
	Rating        float64   `json:"rating" bson:"rating"`
	TotalBookings int64     `json:"total_bookings" bson:"total_bookings"`
	Revenue       float64   `json:"revenue" bson:"revenue"`
	Profit        float64   `json:"profit" bson:"profit"`
	ProfitPercent float64   `json:"profit_percent" bson:"profit_percent"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
	IsFreezed     bool      `json:"is_freezed" bson:"is_freezed"`
}
