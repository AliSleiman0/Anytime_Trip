package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatbotRating represents a user's rating and feedback for the chatbot
type ChatbotRating struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string             `bson:"user_id" json:"user_id"`
	UserEmail string             `bson:"user_email" json:"user_email"`
	UserName  string             `bson:"user_name,omitempty" json:"user_name,omitempty"`
	Rating    int                `bson:"rating" json:"rating"` // 1-5 stars
	Feedback  string             `bson:"feedback,omitempty" json:"feedback,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
