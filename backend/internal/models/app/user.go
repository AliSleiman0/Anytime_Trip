package app

import "time"

// User represents an app user
type User struct {
	ID          string    `json:"id" bson:"_id"`
	Name        string    `json:"name" bson:"name"`
	Email       string    `json:"email" bson:"email"`
	PhoneNumber string    `json:"phone_number" bson:"phone_number"`
	IsActive    bool      `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
