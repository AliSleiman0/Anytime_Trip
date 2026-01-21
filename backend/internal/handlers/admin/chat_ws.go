package admin

import (
	"Anytime_Travel/backend/internal/models/app"
	"Anytime_Travel/backend/internal/ws"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// HandleChatWebSocket upgrades HTTP to WebSocket for admin chat
func (h *AdminHandler) HandleChatWebSocket(c *websocket.Conn) {
	// Get query parameters from URL
	ticketID := c.Query("ticket_id")
	adminID := c.Query("admin_id")
	adminName := c.Query("admin_name")
	userID := c.Query("user_id")
	userName := c.Query("user_name")

	// Require ticket and at least one identity (admin or user)
	if ticketID == "" || (adminID == "" && userID == "") {
		log.Printf("[WebSocket] Missing ticket_id or admin_id/user_id")
		return
	}

	// Determine role: if user_id present, treat as user even if admin_id also supplied
	isAdmin := false
	clientID := ""
	if userID != "" {
		isAdmin = false
		clientID = userID
	} else {
		isAdmin = true
		clientID = adminID
	}

	// Create client
	client := &ws.Client{
		ID:       clientID,
		TicketID: ticketID,
		IsAdmin:  isAdmin,
		Hub:      h.chatHub,
		Conn:     c,
		Send:     make(chan []byte, 256),
	}

	// Register client
	h.chatHub.Register <- client

	// Start write pump in goroutine
	go client.WritePump()

	// Read pump in current goroutine (blocking)
	defer func() {
		h.chatHub.Unregister <- client
		c.Close()
	}()

	for {
		// Read message as map for flexibility
		var msg map[string]interface{}
		if err := c.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket] Read error: %v", err)
			}
			break
		}

		// Extract message text
		messageText, ok := msg["message"].(string)
		if !ok || messageText == "" {
			continue
		}

		// Create reply
		reply := app.TicketReply{
			Message:   messageText,
			IsAdmin:   isAdmin,
			CreatedAt: time.Now(),
		}
		if isAdmin {
			reply.ReadByAdmin = true
			reply.AdminName = adminName
		} else {
			// user message; ensure ReadByAdmin is false
			reply.ReadByAdmin = false
			// Record the user's display name separately so history can show it
			reply.UserName = userName
		}

		// Save to database asynchronously
		go h.saveMessageToDB(ticketID, reply)

		// Prepare broadcast message - preserve whether sender is admin or user
		broadcastMsg := map[string]interface{}{
			"ticket_id":  ticketID,
			"message":    messageText,
			"is_admin":   isAdmin,
			"created_at": reply.CreatedAt.Format(time.RFC3339),
		}
		if isAdmin {
			broadcastMsg["admin_id"] = adminID
			broadcastMsg["admin_name"] = adminName
		} else {
			broadcastMsg["user_id"] = userID
			broadcastMsg["user_name"] = userName
		}

		// Broadcast to room
		msgBytes, _ := json.Marshal(broadcastMsg)
		h.chatHub.Broadcast <- &ws.BroadcastMessage{
			TicketID: ticketID,
			Message:  msgBytes,
		}
	}
}

// saveMessageToDB saves message to database
func (h *AdminHandler) saveMessageToDB(ticketID string, reply app.TicketReply) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticket, err := h.supportTicketRepo.GetTicketByTicketID(ctx, ticketID)
	if err != nil {
		log.Printf("[SaveMessage] Ticket not found: %v", err)
		return
	}

	if err := h.supportTicketRepo.AddReply(ctx, ticket.ID, reply); err != nil {
		log.Printf("[SaveMessage] Save error: %v", err)
	} else {
		log.Printf("[SaveMessage] Saved message for ticket %s", ticketID)
		// If this was a user message, notify all connected admins about the unlocked ticket
		if !reply.IsAdmin && h.notifyHub != nil {
			// Minimal notification payload
			notif := map[string]interface{}{
				"type":      "ticket_unlocked",
				"ticket_id": ticketID,
				"user_name": reply.UserName,
				"message":   reply.Message,
				"timestamp": reply.CreatedAt.Format(time.RFC3339),
			}
			if data, err := json.Marshal(notif); err == nil {
				h.notifyHub.Broadcast <- data
			}
		}
	}
}

// GetChatHistory returns chat history for a ticket
func (h *AdminHandler) GetChatHistory(c *fiber.Ctx) error {
	ticketID := c.Query("ticket_id")
	if ticketID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ticket_id required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ticket, err := h.supportTicketRepo.GetTicketByTicketID(ctx, ticketID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
	}

	// Mark user replies as read by admin (so unread count becomes zero)
	if err := h.supportTicketRepo.MarkRepliesReadByAdmin(ctx, ticket.ID); err != nil {
		log.Printf("[GetChatHistory] Failed to mark replies read: %v", err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"replies": ticket.Replies,
		"ticket":  ticket,
	})
}

// HandleNotifyWebSocket upgrades an admin client to a WebSocket for global notifications
func (h *AdminHandler) HandleNotifyWebSocket(c *websocket.Conn) {
	adminID := c.Query("admin_id")
	if adminID == "" || h.notifyHub == nil {
		log.Printf("[NotifyWS] Missing admin_id or notify hub not configured")
		return
	}

	client := &ws.AdminClient{
		ID:   adminID,
		Conn: c,
		Send: make(chan []byte, 256),
	}

	// Register client
	h.notifyHub.Register <- client

	// Start write pump
	go client.WritePump()

	// Read loop to detect disconnects
	defer func() {
		h.notifyHub.Unregister <- client
		c.Close()
	}()

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			// client disconnected or error
			break
		}
	}
}
