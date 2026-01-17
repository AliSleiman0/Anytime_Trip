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

	collection := dbConn.DB.Collection("support_tickets")

	now := time.Now()

	// Tickets: TicketID length <= 6 characters
	tickets := []interface{}{
		bson.M{
			"_id":            primitive.NewObjectID(),
			"ticket_id":      "T001",
			"user_id":        primitive.NewObjectID(),
			"customer_name":  "John Smith",
			"customer_email": "john.smith@example.com",
			"category":       "Flight",
			"subject":        "Flight cancellation",
			"description":    "I need to cancel my flight booked for tomorrow but cannot in the app.",
			"status":         "Open",
			"priority":       "High",
			"replies":        []interface{}{},
			"created_at":     now.Add(-24 * time.Hour),
			"updated_at":     now.Add(-24 * time.Hour),
		},
		bson.M{
			"_id":            primitive.NewObjectID(),
			"ticket_id":      "T002",
			"user_id":        primitive.NewObjectID(),
			"customer_name":  "Sarah Johnson",
			"customer_email": "sarah.j@example.com",
			"category":       "Hotel",
			"subject":        "No confirmation email",
			"description":    "I haven't received my hotel booking confirmation (ref HB123456).",
			"status":         "In Progress",
			"priority":       "Medium",
			"replies": []interface{}{
				bson.M{"_id": primitive.NewObjectID(), "message": "We're checking your reservation.", "is_admin": true, "admin_name": "Support Team", "created_at": now.Add(-2 * time.Hour)},
			},
			"created_at": now.Add(-48 * time.Hour),
			"updated_at": now.Add(-2 * time.Hour),
		},
		bson.M{
			"_id":            primitive.NewObjectID(),
			"ticket_id":      "T003",
			"user_id":        primitive.NewObjectID(),
			"customer_name":  "Michael Chen",
			"customer_email": "m.chen@example.com",
			"category":       "Payment",
			"subject":        "Double charge",
			"description":    "I was charged twice for a booking, please refund one charge.",
			"status":         "Open",
			"priority":       "Urgent",
			"replies":        []interface{}{},
			"created_at":     now.Add(-6 * time.Hour),
			"updated_at":     now.Add(-6 * time.Hour),
		},
		bson.M{
			"_id":            primitive.NewObjectID(),
			"ticket_id":      "T004",
			"user_id":        primitive.NewObjectID(),
			"customer_name":  "Emily Davis",
			"customer_email": "emily.davis@example.com",
			"category":       "Car",
			"subject":        "Pickup location",
			"description":    "Where is the car pickup counter at LAX?",
			"status":         "Resolved",
			"priority":       "Low",
			"replies": []interface{}{
				bson.M{"_id": primitive.NewObjectID(), "message": "The counter is on ground floor Terminal 1.", "is_admin": true, "admin_name": "Admin", "created_at": now.Add(-4 * time.Hour)},
				bson.M{"_id": primitive.NewObjectID(), "message": "Thanks for the info!", "is_admin": false, "admin_name": "", "created_at": now.Add(-3 * time.Hour)},
			},
			"created_at": now.Add(-12 * time.Hour),
			"updated_at": now.Add(-3 * time.Hour),
		},
		bson.M{
			"_id":            primitive.NewObjectID(),
			"ticket_id":      "T005",
			"user_id":        primitive.NewObjectID(),
			"customer_name":  "David Wilson",
			"customer_email": "d.wilson@example.com",
			"category":       "Flight",
			"subject":        "Seat map not loading",
			"description":    "Seat map won't load for my flight AA1234. Tried multiple browsers.",
			"status":         "In Progress",
			"priority":       "Medium",
			"replies": []interface{}{
				bson.M{"_id": primitive.NewObjectID(), "message": "Escalated to tech team.", "is_admin": true, "admin_name": "Tech Support", "created_at": now.Add(-45 * time.Minute)},
			},
			"created_at": now.Add(-3 * time.Hour),
			"updated_at": now.Add(-45 * time.Minute),
		},
	}

	res, err := collection.InsertMany(ctx, tickets)
	if err != nil {
		log.Fatalf("Failed to insert tickets: %v", err)
	}

	fmt.Printf("Inserted %d support tickets. IDs:\n", len(res.InsertedIDs))
	for i, id := range res.InsertedIDs {
		fmt.Printf(" %d) %v\n", i+1, id)
	}

	fmt.Println("\n✅ Support ticket seed completed successfully!")
}
