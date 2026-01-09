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

	collection := db.DB.Collection("flights")

	now := time.Now()
	flights := []adminModel.Flight{
		{
			ID:            "flight_001",
			FlightID:      "DL-190",
			FlightNumber:  "DL190",
			Status:        adminModel.FlightStatusActive,
			Cost:          0,
			Currency:      "USD",
			ProviderName:  "Delta Air Lines",
			ProviderType:  "Airline",
			Location:      "Atlanta, GA",
			ContactEmail:  "contact@delta.com",
			Rating:        4.5,
			TotalBookings: 15420,
			Revenue:       2345000,
			Profit:        352000,
			ProfitPercent: 15.0,
			CreatedAt:     now.AddDate(0, 0, -30),
			UpdatedAt:     now,
		},
		{
			ID:            "flight_002",
			FlightID:      "UA-450",
			FlightNumber:  "UA450",
			Status:        adminModel.FlightStatusActive,
			Cost:          0,
			Currency:      "USD",
			ProviderName:  "United Airlines",
			ProviderType:  "Airline",
			Location:      "Chicago, IL",
			ContactEmail:  "contact@united.com",
			Rating:        4.3,
			TotalBookings: 12800,
			Revenue:       1980000,
			Profit:        297000,
			ProfitPercent: 15.0,
			CreatedAt:     now.AddDate(0, 0, -20),
			UpdatedAt:     now,
		},
		{
			ID:            "flight_003",
			FlightID:      "AA-320",
			FlightNumber:  "AA320",
			Status:        adminModel.FlightStatusActive,
			Cost:          0,
			Currency:      "USD",
			ProviderName:  "American Airlines",
			ProviderType:  "Airline",
			Location:      "Dallas, TX",
			ContactEmail:  "contact@aa.com",
			Rating:        4.2,
			TotalBookings: 11750,
			Revenue:       1765000,
			Profit:        265000,
			ProfitPercent: 15.0,
			CreatedAt:     now.AddDate(0, 0, -15),
			UpdatedAt:     now,
		},
	}

	for _, flight := range flights {
		opts := options.Replace().SetUpsert(true)
		_, err := collection.ReplaceOne(context.Background(), bson.M{"_id": flight.ID}, flight, opts)
		if err != nil {
			log.Printf("Error upserting flight provider %s: %v", flight.ProviderName, err)
		} else {
			log.Printf("Upserted flight provider: %s", flight.ProviderName)
		}
	}

	log.Println("Flight providers seeding completed.")
}
