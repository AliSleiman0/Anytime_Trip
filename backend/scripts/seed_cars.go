package main

import (
	"context"
	"log"
	"time"

	"Anytime_Travel/backend/config"
	"Anytime_Travel/backend/internal/database"
	"Anytime_Travel/backend/seeds"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	dbConn, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := dbConn.Client.Disconnect(ctx); err != nil {
			log.Fatal("Failed to disconnect from database:", err)
		}
	}()

	log.Println("Starting car seeding...")
	if err := seeds.SeedCars(dbConn.DB); err != nil {
		log.Fatal("Failed to seed cars:", err)
	}
	log.Println("Car seeding completed successfully!")
}
