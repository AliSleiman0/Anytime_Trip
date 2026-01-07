package main

import (
	"context"
	"log"
	"time"

	"Anytime_Travel/backend/config"
	"Anytime_Travel/backend/internal/database"
	adminModel "Anytime_Travel/backend/internal/models/admin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	collection := db.DB.Collection("hotels")

	now := time.Now()
	hotels := []adminModel.Hotel{
		{
			ID:            "hotel_001",
			HotelID:       "MARRIOTT-TIMES-SQ",
			HotelName:     "Marriott Times Square",
			Status:        adminModel.HotelStatusActive,
			Cost:          210.00,
			Currency:      "USD",
			ProviderName:  "Marriott",
			ProviderType:  "Hotel",
			Location:      "New York, NY",
			ContactEmail:  "contact@marriott.com",
			Rating:        4.6,
			TotalBookings: 10230,
			Revenue:       3200000,
			CreatedAt:     now.AddDate(0, 0, -40),
			UpdatedAt:     now,
		},
		{
			ID:            "hotel_002",
			HotelID:       "HILTON-DOWNTOWN-CHI",
			HotelName:     "Hilton Downtown",
			Status:        adminModel.HotelStatusActive,
			Cost:          180.00,
			Currency:      "USD",
			ProviderName:  "Hilton",
			ProviderType:  "Hotel",
			Location:      "Chicago, IL",
			ContactEmail:  "contact@hilton.com",
			Rating:        4.4,
			TotalBookings: 8950,
			Revenue:       2750000,
			CreatedAt:     now.AddDate(0, 0, -28),
			UpdatedAt:     now,
		},
		{
			ID:            "hotel_003",
			HotelID:       "IHG-MIDTOWN-ATL",
			HotelName:     "IHG Midtown",
			Status:        adminModel.HotelStatusActive,
			Cost:          160.00,
			Currency:      "USD",
			ProviderName:  "IHG",
			ProviderType:  "Hotel",
			Location:      "Atlanta, GA",
			ContactEmail:  "contact@ihg.com",
			Rating:        4.3,
			TotalBookings: 7600,
			Revenue:       2100000,
			CreatedAt:     now.AddDate(0, 0, -21),
			UpdatedAt:     now,
		},
	}

	for _, hotel := range hotels {
		opts := options.Replace().SetUpsert(true)
		_, err := collection.ReplaceOne(context.Background(), bson.M{"_id": hotel.ID}, hotel, opts)
		if err != nil {
			log.Printf("Error upserting hotel provider %s: %v", hotel.ProviderName, err)
		} else {
			log.Printf("Upserted hotel provider: %s", hotel.ProviderName)
		}
	}

	log.Println("Hotel providers seeding completed.")
}
