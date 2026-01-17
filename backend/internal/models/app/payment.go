package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Payment represents a recorded payment or transaction
type Payment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookingID     string             `bson:"booking_id" json:"booking_id"`
	Type          string             `bson:"type" json:"type"`
	CustomerName  string             `bson:"customer_name" json:"customer_name"`
	CustomerEmail string             `bson:"customer_email" json:"customer_email"`
	BookingDate   time.Time          `bson:"booking_date" json:"booking_date"`
	Amount        float64            `bson:"amount" json:"amount"`
	PaymentMethod string             `bson:"payment_method" json:"payment_method"`
	PaymentStatus string             `bson:"payment_status" json:"payment_status"`
}
