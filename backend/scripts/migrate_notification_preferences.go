//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"
	"time"

	"Anytime_Travel/backend/config"
	"Anytime_Travel/backend/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// This script migrates old notification preferences to the new simplified structure
// Old structure had 10 options, new structure has only 3: Email, SMS, Chatbot

func main() {
	cfg := config.LoadConfig()
	db, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	collection := db.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Update all users with the new notification preferences structure
	// Map old preferences to new ones:
	// - EmailAlerts -> Email
	// - SMSUpdates -> SMS
	// - ChatbotMessages -> Chatbot

	pipeline := mongo.Pipeline{
		{
			{Key: "$set", Value: bson.D{
				{Key: "notification_preferences", Value: bson.D{
					{Key: "Email", Value: bson.D{
						{Key: "$ifNull", Value: bson.A{
							"$notification_preferences.email_alerts",
							"$notification_preferences.EmailAlerts",
							true,
						}},
					}},
					{Key: "SMS", Value: bson.D{
						{Key: "$ifNull", Value: bson.A{
							"$notification_preferences.sms_updates",
							"$notification_preferences.SMSUpdates",
							true,
						}},
					}},
					{Key: "Chatbot", Value: bson.D{
						{Key: "$ifNull", Value: bson.A{
							"$notification_preferences.chatbot_messages",
							"$notification_preferences.ChatbotMessages",
							true,
						}},
					}},
				}},
			}},
		},
	}

	result, err := collection.UpdateMany(
		ctx,
		bson.M{}, // Match all documents
		pipeline,
	)

	if err != nil {
		log.Fatal("Failed to update notification preferences:", err)
	}

	log.Printf("Successfully migrated notification preferences for %d users", result.ModifiedCount)
	log.Println("Migration completed.")
}
