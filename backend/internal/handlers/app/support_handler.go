package app

import (
	"context"
	"log"
	"math/rand"
	"time"

	"Anytime_Travel/backend/internal/models/app"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateSupportTicket creates a new support ticket for the user
func (h *AppHandler) CreateSupportTicket(c *fiber.Ctx) error {
	var req struct {
		CustomerName  string `json:"customer_name" validate:"required"`
		CustomerEmail string `json:"customer_email" validate:"required,email"`
		Category      string `json:"category" validate:"required"`
		Subject       string `json:"subject" validate:"required"`
		Description   string `json:"description" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.CustomerName == "" || req.CustomerEmail == "" || req.Category == "" ||
		req.Subject == "" || req.Description == "" {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required"})
	}

	// Get user ID from context
	userID := c.Locals("user_id")
	var userIDStr string
	if userID != nil {
		userIDStr = userID.(string)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate ticket ID (6 characters)
	ticketID := generateTicketID()

	ticket := &app.SupportTicket{
		TicketID:      ticketID,
		UserID:        userIDStr, // Store user_id for faster lookups
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		Category:      req.Category,
		Subject:       req.Subject,
		Description:   req.Description,
		Status:        "Open",
		Priority:      "Medium",
		Replies:       []app.TicketReply{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Insert ticket
	result, err := h.supportTicketRepo.Create(ctx, ticket)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create ticket", "details": err.Error()})
	}

	ticket.ID = result

	// Prepare response with string ID
	response := fiber.Map{
		"success": true,
		"message": "Support ticket created successfully",
		"ticket": fiber.Map{
			"id":             ticket.ID.Hex(),
			"ticket_id":      ticket.TicketID,
			"customer_name":  ticket.CustomerName,
			"customer_email": ticket.CustomerEmail,
			"category":       ticket.Category,
			"subject":        ticket.Subject,
			"description":    ticket.Description,
			"status":         ticket.Status,
			"priority":       ticket.Priority,
			"replies":        ticket.Replies,
			"created_at":     ticket.CreatedAt.Format(time.RFC3339),
			"updated_at":     ticket.UpdatedAt.Format(time.RFC3339),
		},
	}

	return c.Status(201).JSON(response)
}

// GetUserSupportTickets retrieves all support tickets for the logged-in user
func (h *AppHandler) GetUserSupportTickets(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID := c.Locals("user_id")
	if userID == nil {
		log.Printf("[GetUserTickets] No user_id in context")
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	log.Printf("[GetUserTickets] Fetching tickets for user_id: %s", userID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tickets, err := h.supportTicketRepo.GetTicketsByUserID(ctx, userID.(string))
	if err != nil {
		log.Printf("[GetUserTickets] Error fetching tickets: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tickets", "details": err.Error()})
	}

	log.Printf("[GetUserTickets] Found %d tickets", len(tickets))

	return c.JSON(fiber.Map{
		"success": true,
		"tickets": tickets,
	})
}

// GetSupportTicketDetails retrieves details of a specific ticket
func (h *AppHandler) GetSupportTicketDetails(c *fiber.Ctx) error {
	ticketID := c.Params("ticketId")
	if ticketID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Ticket ID is required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to get ticket by ticket_id (short ID)
	ticket, err := h.supportTicketRepo.GetTicketByTicketID(ctx, ticketID)
	if err != nil {
		// If not found, try by ObjectID
		objID, err2 := primitive.ObjectIDFromHex(ticketID)
		if err2 != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
		}
		ticket, err = h.supportTicketRepo.GetTicketByID(ctx, objID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
		}
	}

	return c.JSON(ticket)
}

// ReplyToSupportTicket adds a user reply to a ticket
func (h *AppHandler) ReplyToSupportTicket(c *fiber.Ctx) error {
	ticketID := c.Params("ticketId")
	if ticketID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Ticket ID is required"})
	}

	var req struct {
		Message string `json:"message" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Message == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Message is required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get ticket
	ticket, err := h.supportTicketRepo.GetTicketByTicketID(ctx, ticketID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
	}

	// Create reply
	reply := app.TicketReply{
		Message:     req.Message,
		IsAdmin:     false,
		ReadByAdmin: false,
		CreatedAt:   time.Now(),
	}

	// Add reply to ticket
	if err := h.supportTicketRepo.AddReply(ctx, ticket.ID, reply); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add reply", "details": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Reply added successfully",
	})
}

// GetChatHistory retrieves the chat history for a ticket
func (h *AppHandler) GetChatHistory(c *fiber.Ctx) error {
	ticketID := c.Query("ticket_id")
	if ticketID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ticket_id parameter is required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get ticket by ticket_id
	ticket, err := h.supportTicketRepo.GetTicketByTicketID(ctx, ticketID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
	}

	// Return ticket with all replies
	return c.JSON(fiber.Map{
		"success": true,
		"ticket":  ticket,
	})
}

// generateTicketID generates a random 6-character ticket ID
func generateTicketID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
