package app

import "time"

// NotificationPreferences represents user notification settings
type NotificationPreferences struct {
	AllNotifications     bool `json:"all_notifications" bson:"all_notifications"`
	PushNotifications    bool `json:"push_notifications" bson:"push_notifications"`
	EmailAlerts          bool `json:"email_alerts" bson:"email_alerts"`
	SMSUpdates           bool `json:"sms_updates" bson:"sms_updates"`
	InAppMessages        bool `json:"in_app_messages" bson:"in_app_messages"`
	SocialMediaAlerts    bool `json:"social_media_alerts" bson:"social_media_alerts"`
	WebhookNotifications bool `json:"webhook_notifications" bson:"webhook_notifications"`
	BrowserNotifications bool `json:"browser_notifications" bson:"browser_notifications"`
	RSSFeedUpdates       bool `json:"rss_feed_updates" bson:"rss_feed_updates"`
	ChatbotMessages      bool `json:"chatbot_messages" bson:"chatbot_messages"`
}

// SecurityPreferences represents user security settings
type SecurityPreferences struct {
	TwoFactorAuth bool `json:"two_factor_auth" bson:"two_factor_auth"`
}

// User represents an app user
type User struct {
	ID                      string                  `json:"id" bson:"_id"`
	Name                    string                  `json:"name" bson:"name"`
	Email                   string                  `json:"email" bson:"email"`
	PhoneNumber             string                  `json:"phone_number" bson:"phone_number"`
	PasswordHash            string                  `json:"-" bson:"password_hash"` // Hidden from JSON responses
	ProfileImage            string                  `json:"profile_image" bson:"profile_image,omitempty"`
	Sex                     string                  `json:"sex" bson:"sex"`
	Country                 string                  `json:"country" bson:"country"`
	IsActive                bool                    `json:"is_active" bson:"is_active"`
	IsFreezed               bool                    `json:"is_freezed" bson:"is_freezed"`
	TotalBookings           int                     `json:"total_bookings" bson:"total_bookings"`
	LastBooking             time.Time               `json:"last_booking" bson:"last_booking,omitempty"`
	NotificationPreferences NotificationPreferences `json:"notification_preferences" bson:"notification_preferences,omitempty"`
	SecurityPreferences     SecurityPreferences     `json:"security_preferences" bson:"security_preferences,omitempty"`
	CreatedAt               time.Time               `json:"created_at" bson:"created_at"`
	LastLogin               time.Time               `json:"last_login" bson:"last_login,omitempty"` // New field to track last login
}

// SignupRequest represents the signup request payload
type SignupRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Sex             string `json:"sex" validate:"required,oneof=Male Female Other"`
	Country         string `json:"country" validate:"required"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
