//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"
	"time"

	"travel/backend/config"
	"travel/backend/internal/database"
	"travel/backend/internal/models/app"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	collection := db.DB.Collection("car_bookings")

	now := time.Now()
	bookings := []app.CarBooking{
		{
			ID:            "car_bk_001",
			BookingID:     "CBK-1001",
			UserID:        "user_001",
			CarID:         "car_001",
			Status:        app.CarBookingStatusConfirmed,
			Customer:      app.CarCustomer{Name: "John Doe", Email: "john.doe@example.com"},
			Details:       "SUV rental - 3 days",
			BookingDate:   now.AddDate(0, 0, -2),
			Amount:        220.50,
			Currency:      "USD",
			PaymentStatus: app.CarPaymentStatusPaid,
			CreatedAt:     now.AddDate(0, 0, -2),
			UpdatedAt:     now,
		},
		{
			ID:            "car_bk_002",
			BookingID:     "CBK-1002",
			UserID:        "user_002",
			CarID:         "car_002",
			Status:        app.CarBookingStatusPending,
			Customer:      app.CarCustomer{Name: "Jane Smith", Email: "jane.smith@example.com"},
			Details:       "Compact rental - weekend",
			BookingDate:   now.AddDate(0, 0, -1),
			Amount:        140.00,
			Currency:      "USD",
			PaymentStatus: app.CarPaymentStatusPending,
			CreatedAt:     now.AddDate(0, 0, -1),
			UpdatedAt:     now,
		},
		{
			ID:            "car_bk_003",
			BookingID:     "CBK-1003",
			UserID:        "user_003",
			CarID:         "car_003",
			Status:        app.CarBookingStatusCancelled,
			Customer:      app.CarCustomer{Name: "Alice Johnson", Email: "alice.johnson@example.com"},
			Details:       "Sedan rental - 5 days",
			BookingDate:   now.AddDate(0, 0, -5),
			Amount:        310.75,
			Currency:      "USD",
			PaymentStatus: app.CarPaymentStatusRefunded,
			CreatedAt:     now.AddDate(0, 0, -5),
			UpdatedAt:     now,
		},
	}

	for _, booking := range bookings {
		opts := options.Replace().SetUpsert(true)
		_, err := collection.ReplaceOne(context.Background(), map[string]string{"_id": booking.ID}, booking, opts)
		if err != nil {
			log.Printf("Error upserting car booking %s: %v", booking.BookingID, err)
		} else {
			log.Printf("Successfully upserted car booking: %s", booking.BookingID)
		}
	}

	log.Println("Car bookings seeding completed.")
}
