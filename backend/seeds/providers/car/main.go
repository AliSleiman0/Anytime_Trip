package main

import (
	"context"
	"log"
	"time"

	"travel/backend/config"
	"travel/backend/internal/database"
	adminModel "travel/backend/internal/models/admin"

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

	collection := db.DB.Collection("cars")

	now := time.Now()
	cars := []adminModel.Car{
		{
			ID:            "car_001",
			CarID:         "HERTZ-RAV4",
			CarName:       "Toyota RAV4",
			Status:        adminModel.CarStatusActive,
			Cost:          45.00,
			Currency:      "USD",
			ProviderName:  "Hertz",
			ProviderType:  "Car Rental",
			Location:      "Atlanta, GA",
			ContactEmail:  "support@hertz.com",
			Rating:        4.2,
			TotalBookings: 8200,
			Revenue:       740000,
			Profit:        118000,
			ProfitPercent: 16.0,
			CreatedAt:     now.AddDate(0, 0, -25),
			UpdatedAt:     now,
		},
		{
			ID:            "car_002",
			CarID:         "AVIS-CAMRY",
			CarName:       "Toyota Camry",
			Status:        adminModel.CarStatusActive,
			Cost:          40.00,
			Currency:      "USD",
			ProviderName:  "Avis",
			ProviderType:  "Car Rental",
			Location:      "Chicago, IL",
			ContactEmail:  "support@avis.com",
			Rating:        4.1,
			TotalBookings: 6950,
			Revenue:       610000,
			Profit:        91500,
			ProfitPercent: 15.0,
			CreatedAt:     now.AddDate(0, 0, -18),
			UpdatedAt:     now,
		},
		{
			ID:            "car_003",
			CarID:         "ENTERPRISE-FOCUS",
			CarName:       "Ford Focus",
			Status:        adminModel.CarStatusActive,
			Cost:          38.00,
			Currency:      "USD",
			ProviderName:  "Enterprise",
			ProviderType:  "Car Rental",
			Location:      "Dallas, TX",
			ContactEmail:  "support@enterprise.com",
			Rating:        4.0,
			TotalBookings: 7300,
			Revenue:       620000,
			Profit:        93000,
			ProfitPercent: 15.0,
			CreatedAt:     now.AddDate(0, 0, -12),
			UpdatedAt:     now,
		},
	}

	for _, car := range cars {
		opts := options.Replace().SetUpsert(true)
		_, err := collection.ReplaceOne(context.Background(), bson.M{"_id": car.ID}, car, opts)
		if err != nil {
			log.Printf("Error upserting car provider %s: %v", car.ProviderName, err)
		} else {
			log.Printf("Upserted car provider: %s", car.ProviderName)
		}
	}

	log.Println("Car providers seeding completed.")
}
