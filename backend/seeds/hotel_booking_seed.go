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

	collection := db.DB.Collection("hotel_bookings")

	now := time.Now()
	bookings := []app.HotelBooking{
		{
			ID:            "hotel_bk_001",
			BookingID:     "HBK-3001",
			UserID:        "user_001",
			HotelID:       "hotel_001",
			Status:        app.HotelBookingStatusConfirmed,
			Customer:      app.HotelCustomer{Name: "John Doe", Email: "john.doe@example.com"},
			Details:       "Deluxe room, 3 nights",
			BookingDate:   now.AddDate(0, 0, -4),
			Amount:        620.00,
			Currency:      "USD",
			PaymentStatus: app.HotelPaymentStatusPaid,
			CreatedAt:     now.AddDate(0, 0, -4),
			UpdatedAt:     now,
		},
		{
			ID:            "hotel_bk_002",
			BookingID:     "HBK-3002",
			UserID:        "user_002",
			HotelID:       "hotel_002",
			Status:        app.HotelBookingStatusPending,
			Customer:      app.HotelCustomer{Name: "Jane Smith", Email: "jane.smith@example.com"},
			Details:       "Standard room, weekend",
			BookingDate:   now.AddDate(0, 0, -2),
			Amount:        280.00,
			Currency:      "USD",
			PaymentStatus: app.HotelPaymentStatusPending,
			CreatedAt:     now.AddDate(0, 0, -2),
			UpdatedAt:     now,
		},
		{
			ID:            "hotel_bk_003",
			BookingID:     "HBK-3003",
			UserID:        "user_003",
			HotelID:       "hotel_003",
			Status:        app.HotelBookingStatusCancelled,
			Customer:      app.HotelCustomer{Name: "Alice Johnson", Email: "alice.johnson@example.com"},
			Details:       "Suite, 2 nights",
			BookingDate:   now.AddDate(0, 0, -10),
			Amount:        480.00,
			Currency:      "USD",
			PaymentStatus: app.HotelPaymentStatusRefunded,
			CreatedAt:     now.AddDate(0, 0, -10),
			UpdatedAt:     now,
		},
	}

	for _, booking := range bookings {
		opts := options.Replace().SetUpsert(true)
		_, err := collection.ReplaceOne(context.Background(), map[string]string{"_id": booking.ID}, booking, opts)
		if err != nil {
			log.Printf("Error upserting hotel booking %s: %v", booking.BookingID, err)
		} else {
			log.Printf("Successfully upserted hotel booking: %s", booking.BookingID)
		}
	}

	log.Println("Hotel bookings seeding completed.")
}
