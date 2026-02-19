//go:build ignore
// +build ignore

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
	checkIn := now.AddDate(0, 0, 7)   // Check-in 7 days from now
	checkOut := now.AddDate(0, 0, 10) // Check-out 10 days from now (3 nights)

	bookings := []app.HotelBooking{
		{
			ID:        "hotel_bk_001",
			BookingID: "HBK-3001",
			UserID:    "user_001",
			HotelID:   "hotel_001",
			Status:    app.HotelBookingStatusConfirmed,
			Customer:  app.HotelCustomer{Name: "John Doe", Email: "john.doe@example.com"},
			Guests: []app.HotelGuest{
				{Name: "John Doe", Type: "Adult"},
				{Name: "Jane Doe", Type: "Adult"},
			},
			HotelName:    "Beirut Hotel",
			HotelImage:   "https://example.com/hotels/beirut-hotel.jpg",
			CheckInDate:  checkIn,
			CheckOutDate: checkOut,
			Nights:       3,
			Rooms:        1,
			Pricing: app.HotelPricing{
				NightPrice:     180.00,
				Taxes:          54.00,
				DestinationFee: 15.00,
				ServiceFee:     25.00,
				Total:          620.00,
			},
			Details:       "Deluxe room, 3 nights",
			BookingDate:   now.AddDate(0, 0, -4),
			Amount:        620.00,
			Currency:      "USD",
			PaymentStatus: app.HotelPaymentStatusPaid,
			CreatedAt:     now.AddDate(0, 0, -4),
			UpdatedAt:     now,
		},
		{
			ID:        "hotel_bk_002",
			BookingID: "HBK-3002",
			UserID:    "user_002",
			HotelID:   "hotel_002",
			Status:    app.HotelBookingStatusPending,
			Customer:  app.HotelCustomer{Name: "Jane Smith", Email: "jane.smith@example.com"},
			Guests: []app.HotelGuest{
				{Name: "Jane Smith", Type: "Adult"},
			},
			HotelName:    "Seaside Resort",
			HotelImage:   "https://example.com/hotels/seaside-resort.jpg",
			CheckInDate:  now.AddDate(0, 0, 14),
			CheckOutDate: now.AddDate(0, 0, 16),
			Nights:       2,
			Rooms:        1,
			Pricing: app.HotelPricing{
				NightPrice:     120.00,
				Taxes:          36.00,
				DestinationFee: 10.00,
				ServiceFee:     14.00,
				Total:          280.00,
			},
			Details:       "Standard room, weekend",
			BookingDate:   now.AddDate(0, 0, -2),
			Amount:        280.00,
			Currency:      "USD",
			PaymentStatus: app.HotelPaymentStatusPending,
			CreatedAt:     now.AddDate(0, 0, -2),
			UpdatedAt:     now,
		},
		{
			ID:        "hotel_bk_003",
			BookingID: "HBK-3003",
			UserID:    "user_003",
			HotelID:   "hotel_003",
			Status:    app.HotelBookingStatusCancelled,
			Customer:  app.HotelCustomer{Name: "Alice Johnson", Email: "alice.johnson@example.com"},
			Guests: []app.HotelGuest{
				{Name: "Alice Johnson", Type: "Adult"},
				{Name: "Bob Johnson", Type: "Adult"},
				{Name: "Charlie Johnson", Type: "Child"},
			},
			HotelName:    "Mountain View Hotel",
			HotelImage:   "https://example.com/hotels/mountain-view.jpg",
			CheckInDate:  now.AddDate(0, 0, -5),
			CheckOutDate: now.AddDate(0, 0, -3),
			Nights:       2,
			Rooms:        1,
			Pricing: app.HotelPricing{
				NightPrice:     200.00,
				Taxes:          60.00,
				DestinationFee: 12.00,
				ServiceFee:     28.00,
				Total:          480.00,
			},
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
