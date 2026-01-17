package admin

import "time"

// AdminUser represents an admin user with extended permissions
type AdminUser struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	Currency    string    `json:"currency"`
	Language    string    `json:"language"`
	CreatedAt   time.Time `json:"created_at"`
}
