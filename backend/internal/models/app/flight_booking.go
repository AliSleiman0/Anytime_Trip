package app

import "time"

// FlightPaymentStatus represents the payment status of a flight booking
type FlightPaymentStatus string

const (
	FlightPaymentStatusPending   FlightPaymentStatus = "pending"
	FlightPaymentStatusPaid      FlightPaymentStatus = "paid"
	FlightPaymentStatusByWhish   FlightPaymentStatus = "by_whish"
	FlightPaymentStatusRefunded  FlightPaymentStatus = "refunded"
	FlightPaymentStatusCancelled FlightPaymentStatus = "cancelled"
)

// FlightBookingStatus represents the status of a flight booking
type FlightBookingStatus string

const (
	FlightBookingStatusConfirmed FlightBookingStatus = "confirmed"
	FlightBookingStatusPending   FlightBookingStatus = "pending"
	FlightBookingStatusCancelled FlightBookingStatus = "cancelled"
	FlightBookingStatusCompleted FlightBookingStatus = "completed"
)

// FlightCustomer contains customer information for a flight booking
type FlightCustomer struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

// FlightTraveler represents a traveler in the booking
type FlightTraveler struct {
	Name string `json:"name" bson:"name"`
	Type string `json:"type" bson:"type"` // Adult, Child, Infant
}

// FlightSegment represents a flight segment in the journey
type FlightSegment struct {
	Airline          string `json:"airline" bson:"airline"`
	FlightNumber     string `json:"flight_number" bson:"flight_number"`
	Class            string `json:"class" bson:"class"`
	DepartureCode    string `json:"departure_code" bson:"departure_code"`
	DepartureCity    string `json:"departure_city" bson:"departure_city"`
	DepartureAirport string `json:"departure_airport" bson:"departure_airport"`
	DepartureDate    string `json:"departure_date" bson:"departure_date"`
	DepartureTime    string `json:"departure_time" bson:"departure_time"`
	ArrivalCode      string `json:"arrival_code" bson:"arrival_code"`
	ArrivalCity      string `json:"arrival_city" bson:"arrival_city"`
	ArrivalAirport   string `json:"arrival_airport" bson:"arrival_airport"`
	ArrivalDate      string `json:"arrival_date" bson:"arrival_date"`
	ArrivalTime      string `json:"arrival_time" bson:"arrival_time"`
	Duration         string `json:"duration" bson:"duration"`
}

// FlightPricing represents the pricing breakdown for the booking
type FlightPricing struct {
	BasePrice  float64 `json:"base_price" bson:"base_price"`
	Taxes      float64 `json:"taxes" bson:"taxes"`
	ServiceFee float64 `json:"service_fee" bson:"service_fee"`
	Total      float64 `json:"total" bson:"total"`
	Currency   string  `json:"currency" bson:"currency"`
}

// FlightTripType represents the type of flight trip
type FlightTripType string

const (
	FlightTripTypeOneWay    FlightTripType = "one-way"
	FlightTripTypeRoundTrip FlightTripType = "round-trip"
	FlightTripTypeMultiCity FlightTripType = "multi-city"
)

// FlightBooking represents a user's flight booking
type FlightBooking struct {
	ID              string              `json:"id" bson:"_id"`
	BookingID       string              `json:"booking_id" bson:"booking_id"` // e.g., "ANT-2847-MC"
	UserID          string              `json:"user_id" bson:"user_id"`
	FlightID        string              `json:"flight_id,omitempty" bson:"flight_id,omitempty"` // Reference to flight in admin collection (optional)
	TripType        FlightTripType      `json:"trip_type" bson:"trip_type"`                     // one-way, round-trip, multi-city
	Status          FlightBookingStatus `json:"status" bson:"status"`
	Customer        FlightCustomer      `json:"customer" bson:"customer"`
	Travelers       []FlightTraveler    `json:"travelers" bson:"travelers"`
	OutboundFlights []FlightSegment     `json:"outbound_flights" bson:"outbound_flights"`
	ReturnFlights   []FlightSegment     `json:"return_flights,omitempty" bson:"return_flights,omitempty"`
	Pricing         FlightPricing       `json:"pricing" bson:"pricing"`
	Details         string              `json:"details,omitempty" bson:"details,omitempty"` // Additional details/notes
	BookingDate     time.Time           `json:"booking_date" bson:"booking_date"`
	PaymentStatus   FlightPaymentStatus `json:"payment_status" bson:"payment_status"`
	PaymentMethod   string              `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	CreatedAt       time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at" bson:"updated_at"`
}
