//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"Anytime_Travel/backend/config"
	"Anytime_Travel/backend/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	cfg := config.LoadConfig()

	dbConn, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := dbConn.DB.Collection("payments")

	now := time.Now()

	docs := []interface{}{
		bson.M{
			"_id":            primitive.NewObjectID(),
			"booking_id":     "BK001",
			"type":           "Flight",
			"customer_name":  "John Smith",
			"customer_email": "john.smith@gmail.com",
			"booking_date":   now.AddDate(0, -1, 0),
			"amount":         850.00,
			"payment_method": "Card",
			"payment_status": "Paid",
		},
		bson.M{
			"_id":            primitive.NewObjectID(),
			"booking_id":     "BK002",
			"type":           "Hotel",
			"customer_name":  "Sarah Johnson",
			"customer_email": "sarah.j@example.com",
			"booking_date":   now.AddDate(0, 0, -10),
			"amount":         420.50,
			"payment_method": "Bank Transfer",
			"payment_status": "Pending",
		},
		bson.M{
			"_id":            primitive.NewObjectID(),
			"booking_id":     "BK003",
			"type":           "Car",
			"customer_name":  "Michael Chen",
			"customer_email": "m.chen@example.com",
			"booking_date":   now.AddDate(0, 0, -3),
			"amount":         120.00,
			"payment_method": "Wallet",
			"payment_status": "Refunded",
		},
		bson.M{
			"_id":            primitive.NewObjectID(),
			"booking_id":     "BK004",
			"type":           "Flight",
			"customer_name":  "Emily Davis",
			"customer_email": "emily.davis@example.com",
			"booking_date":   now.AddDate(0, 0, -30),
			"amount":         610.00,
			"payment_method": "Card",
			"payment_status": "Paid",
		},
	}

	res, err := coll.InsertMany(ctx, docs)
	if err != nil {
		log.Fatalf("Failed to insert payments: %v", err)
	}

	fmt.Printf("Inserted %d payment records. IDs:\n", len(res.InsertedIDs))
	for i, id := range res.InsertedIDs {
		fmt.Printf(" %d) %v\n", i+1, id)
	}

	fmt.Println("\n✅ Payments seed completed successfully!")
}
