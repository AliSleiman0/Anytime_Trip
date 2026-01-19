package superadmin

import "time"

// SystemConfig represents system-wide configuration
type SystemConfig struct {
	ID          int       `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SuperAdminUser represents a super admin user
type SuperAdminUser struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	FullAccess bool      `json:"full_access"`
	CreatedAt  time.Time `json:"created_at"`
}

// PredefinedAnswer represents a predefined answer for customer support
type PredefinedAnswer struct {
	ID        int       `json:"id" bson:"id"`
	Shortcut  string    `json:"shortcut" bson:"shortcut"`
	Answer    string    `json:"answer" bson:"answer"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
