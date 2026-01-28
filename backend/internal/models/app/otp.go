package app

import "time"

// OTPType represents the type of OTP
type OTPType string

const (
	OTPTypeEmail OTPType = "email"
	OTPTypePhone OTPType = "phone"
)

// OTP represents an OTP verification code
type OTP struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	Phone     string    `json:"phone,omitempty" bson:"phone,omitempty"`
	Code      string    `json:"-" bson:"code"` // Hidden from JSON
	Type      OTPType   `json:"type" bson:"type"`
	Verified  bool      `json:"verified" bson:"verified"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// SendOTPRequest represents the request to send OTP
type SendOTPRequest struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Type  string `json:"type" validate:"required,oneof=email phone"`
}

// VerifyOTPRequest represents the request to verify OTP
type VerifyOTPRequest struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Code  string `json:"code" validate:"required"`
	Type  string `json:"type" validate:"required,oneof=email phone"`
}
