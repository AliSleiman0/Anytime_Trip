package main

import (
	"context"
	"log"
	"time"

	"Anytime_Travel/backend/config"
	"Anytime_Travel/backend/internal/database"
	"Anytime_Travel/backend/internal/models/app"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	collection := db.DB.Collection("flight_bookings")

	now := time.Now()
	bookings := []app.FlightBooking{
		{
			ID:            "flight_bk_001",
			BookingID:     "FBK-2001",
			UserID:        "user_001",
			FlightID:      "flight_001",
			Status:        app.FlightBookingStatusConfirmed,
			Customer:      app.FlightCustomer{Name: "John Doe", Email: "john.doe@example.com"},
			Details:       "ATL-JFK DL190", // example PNR-like detail
			BookingDate:   now.AddDate(0, 0, -3),
			Amount:        480.00,
			Currency:      "USD",
			PaymentStatus: app.FlightPaymentStatusPaid,
			CreatedAt:     now.AddDate(0, 0, -3),
			UpdatedAt:     now,
		},
		{
			ID:            "flight_bk_002",
			BookingID:     "FBK-2002",
			UserID:        "user_002",
			FlightID:      "flight_002",
			Status:        app.FlightBookingStatusPending,
			Customer:      app.FlightCustomer{Name: "Jane Smith", Email: "jane.smith@example.com"},
			Details:       "SFO-LAX UA450",
			BookingDate:   now.AddDate(0, 0, -1),
			Amount:        220.50,
			Currency:      "USD",
			PaymentStatus: app.FlightPaymentStatusPending,
			CreatedAt:     now.AddDate(0, 0, -1),
			UpdatedAt:     now,
		},
		{
			ID:            "flight_bk_003",
			BookingID:     "FBK-2003",
			UserID:        "user_003",
			FlightID:      "flight_003",
			Status:        app.FlightBookingStatusCancelled,
			Customer:      app.FlightCustomer{Name: "Alice Johnson", Email: "alice.johnson@example.com"},
			Details:       "LHR-CDG AF120",
			BookingDate:   now.AddDate(0, 0, -7),
			Amount:        350.00,
			Currency:      "EUR",
			PaymentStatus: app.FlightPaymentStatusRefunded,
			CreatedAt:     now.AddDate(0, 0, -7),
			UpdatedAt:     now,
		},
	}

	for _, booking := range bookings {
		opts := options.Replace().SetUpsert(true)
		_, err := collection.ReplaceOne(context.Background(), map[string]string{"_id": booking.ID}, booking, opts)
		if err != nil {
			log.Printf("Error upserting flight booking %s: %v", booking.BookingID, err)
		} else {
			log.Printf("Successfully upserted flight booking: %s", booking.BookingID)
		}
	}

	log.Println("Flight bookings seeding completed.")
}
