package admin

import "time"

// TransferStatus represents the status of a transfer
type TransferStatus string

const (
	TransferStatusActive    TransferStatus = "active"
	TransferStatusInactive  TransferStatus = "inactive"
	TransferStatusCancelled TransferStatus = "cancelled"
	TransferStatusDelayed   TransferStatus = "delayed"
)

// Transfer represents a transfer with minimal stored data
// Full transfer details are fetched from external API using TransferID
type Transfer struct {
	ID              string         `json:"id" bson:"_id"`
	TransferID      string         `json:"transfer_id" bson:"transfer_id"`   // Transfer identifier for API lookup
	ServiceName     string         `json:"service_name" bson:"service_name"` // Display name like "Airport Express"
	Status          TransferStatus `json:"status" bson:"status"`
	Cost            float64        `json:"cost" bson:"cost"`         // Top-up or base cost
	Currency        string         `json:"currency" bson:"currency"` // e.g., "USD", "$"
	PickupTime      time.Time      `json:"pickup_time" bson:"pickup_time"`
	DropoffTime     time.Time      `json:"dropoff_time" bson:"dropoff_time"`
	PickupLocation  string         `json:"pickup_location" bson:"pickup_location"`
	DropoffLocation string         `json:"dropoff_location" bson:"dropoff_location"`
	// Provider metadata for admin service-provider table
	ProviderName    string    `json:"provider_name" bson:"provider_name"`
	ProviderType    string    `json:"type" bson:"type"` // e.g., "Transfer Service"
	Location        string    `json:"location" bson:"location"`
	ContactEmail    string    `json:"contact_email" bson:"contact_email"`
	Rating          float64   `json:"rating" bson:"rating"`
	TotalBookings   int64     `json:"total_bookings" bson:"total_bookings"`
	Revenue         float64   `json:"revenue" bson:"revenue"`
	Profit          float64   `json:"profit" bson:"profit"`
	ProfitPercent   float64   `json:"profit_percent" bson:"profit_percent"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
	IsFreezed       bool      `json:"is_freezed" bson:"is_freezed"`
	Refundable      bool      `json:"refundable" bson:"refundable"`
	AllowChanges    bool      `json:"allow_changes" bson:"allow_changes"`
	MeetAndGreet    bool      `json:"meet_and_greet" bson:"meet_and_greet"`
	LuggageIncluded bool      `json:"luggage_included" bson:"luggage_included"`
}
