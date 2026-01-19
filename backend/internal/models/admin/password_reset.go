package admin

import "time"

// PasswordReset represents a password reset token
type PasswordReset struct {
	ID        string    `bson:"_id" json:"id"`
	Email     string    `bson:"email" json:"email"`
	Token     string    `bson:"token" json:"token"`
	Code      string    `bson:"code" json:"code"`
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	Used      bool      `bson:"used" json:"used"`
}
