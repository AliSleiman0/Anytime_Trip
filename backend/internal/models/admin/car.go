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
	ID        string    `json:"id" bson:"_id"`
	CarID     string    `json:"car_id" bson:"car_id"`     // Car identifier for API lookup
	CarName   string    `json:"car_name" bson:"car_name"` // Display name like "Toyota RAV4"
	Status    CarStatus `json:"status" bson:"status"`
	Cost      float64   `json:"cost" bson:"cost"`         // Rental cost
	Currency  string    `json:"currency" bson:"currency"` // e.g., "USD", "$"
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
