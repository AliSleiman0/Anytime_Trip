package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentMethod represents a saved payment method for a user
type PaymentMethod struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"user_id" json:"user_id"`
	CardHolderName  string             `bson:"card_holder_name" json:"card_holder_name"`   // Can be stored encrypted in production
	CardNumberLast4 string             `bson:"card_number_last4" json:"card_number_last4"` // Only last 4 digits
	CardBrand       string             `bson:"card_brand" json:"card_brand"`               // Visa, Mastercard, Amex, etc
	ExpiryMonth     int                `bson:"expiry_month" json:"expiry_month"`
	ExpiryYear      int                `bson:"expiry_year" json:"expiry_year"`
	BillingAddress  string             `bson:"billing_address" json:"billing_address"`
	City            string             `bson:"city" json:"city"`
	State           string             `bson:"state" json:"state"`
	PostalCode      string             `bson:"postal_code" json:"postal_code"`
	Country         string             `bson:"country" json:"country"`
	PhoneNumber     string             `bson:"phone_number" json:"phone_number"`
	IsDefault       bool               `bson:"is_default" json:"is_default"` // Set as default payment method
	IsActive        bool               `bson:"is_active" json:"is_active"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// AddPaymentMethodRequest represents the add payment method request payload
// NOTE: Full card number and CVV are only received for this single transaction
// and should be immediately validated/tokenized by payment processor.
// Only last 4 digits are stored in database.
type AddPaymentMethodRequest struct {
	CardHolderName string `json:"card_holder_name" validate:"required"`
	CardNumber     string `json:"card_number" validate:"required,len=16"` // Full number only in memory, not stored
	CardBrand      string `json:"card_brand" validate:"required,oneof=Visa Mastercard Amex Discover"`
	ExpiryMonth    int    `json:"expiry_month" validate:"required,min=1,max=12"`
	ExpiryYear     int    `json:"expiry_year" validate:"required,min=2024"`
	CVV            string `json:"cvv" validate:"required"` // Only for initial validation, NEVER stored
	BillingAddress string `json:"billing_address" validate:"required"`
	City           string `json:"city" validate:"required"`
	State          string `json:"state" validate:"required"`
	PostalCode     string `json:"postal_code" validate:"required"`
	Country        string `json:"country" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required"`
	IsDefault      bool   `json:"is_default"`
}

// UpdatePaymentMethodRequest represents the update payment method request payload
type UpdatePaymentMethodRequest struct {
	CardHolderName string `json:"card_holder_name"`
	BillingAddress string `json:"billing_address"`
	City           string `json:"city"`
	State          string `json:"state"`
	PostalCode     string `json:"postal_code"`
	Country        string `json:"country"`
	PhoneNumber    string `json:"phone_number"`
	IsDefault      bool   `json:"is_default"`
}
