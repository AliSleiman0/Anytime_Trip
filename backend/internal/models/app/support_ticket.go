package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TicketReply represents a reply on a support ticket
type TicketReply struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Message     string             `bson:"message" json:"message"`
	IsAdmin     bool               `bson:"is_admin" json:"is_admin"`
	ReadByAdmin bool               `bson:"read_by_admin,omitempty" json:"read_by_admin,omitempty"`
	AdminName   string             `bson:"admin_name,omitempty" json:"admin_name,omitempty"`
	UserName    string             `bson:"user_name,omitempty" json:"user_name,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

// SupportTicket represents a support ticket with embedded replies
// Field names are exported to match template usage (e.g. CustomerName, TicketID, Replies)
type SupportTicket struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TicketID      string             `bson:"ticket_id" json:"ticket_id"` // short id (<=6 chars) for display
	CustomerName  string             `bson:"customer_name" json:"customer_name"`
	CustomerEmail string             `bson:"customer_email" json:"customer_email"`
	Category      string             `bson:"category" json:"category"`
	Subject       string             `bson:"subject" json:"subject"`
	Description   string             `bson:"description" json:"description"`
	Status        string             `bson:"status" json:"status"`
	Priority      string             `bson:"priority,omitempty" json:"priority,omitempty"`
	Replies       []TicketReply      `bson:"replies,omitempty" json:"replies,omitempty"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	// UnreadCount is computed at runtime (not stored in DB)
	UnreadCount int `bson:"-" json:"unread_count,omitempty"`
	// Last message preview and timestamp (computed at runtime)
	LastMessage   string    `bson:"-" json:"last_message,omitempty"`
	LastMessageAt time.Time `bson:"-" json:"last_message_at,omitempty"`
}
