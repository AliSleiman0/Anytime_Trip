package admin

import (
	"context"
	"html/template"
	"path/filepath"
	"time"

	appmodels "Anytime_Travel/backend/internal/models/app"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetSupport renders the support page
func (h *AdminHandler) GetSupport(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "support.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Customer Support",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetSupportFragment serves the support fragment with ticket data
func (h *AdminHandler) GetSupportFragment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tickets, err := h.supportTicketRepo.GetAllTickets(ctx)
	if err != nil {
		return c.Status(500).SendString("Error fetching tickets: " + err.Error())
	}

	// compute unread counts for admin (user replies not yet read)
	for i := range tickets {
		count := 0
		// determine last message and timestamp
		if len(tickets[i].Replies) > 0 {
			last := tickets[i].Replies[len(tickets[i].Replies)-1]
			tickets[i].LastMessage = last.Message
			tickets[i].LastMessageAt = last.CreatedAt
		} else {
			tickets[i].LastMessage = tickets[i].Subject
			tickets[i].LastMessageAt = tickets[i].CreatedAt
		}

		for _, r := range tickets[i].Replies {
			if !r.IsAdmin && !r.ReadByAdmin {
				count++
			}
		}
		tickets[i].UnreadCount = count
	}

	// Create template with custom functions
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	}

	tmpl, err := template.New("support-frag.html").Funcs(funcMap).ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "support-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template: " + err.Error())
	}

	data := fiber.Map{"Tickets": tickets}
	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetAllTickets returns all support tickets as JSON
func (h *AdminHandler) GetAllTickets(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tickets, err := h.supportTicketRepo.GetAllTickets(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tickets", "details": err.Error()})
	}
	return c.JSON(tickets)
}

// GetTicketDetail returns a single ticket by ObjectID or ticket_id
func (h *AdminHandler) GetTicketDetail(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Params("id")
	var ticket *appmodels.SupportTicket
	// try hex ObjectID first
	if objID, err := primitive.ObjectIDFromHex(idParam); err == nil {
		t, err := h.supportTicketRepo.GetTicketByID(ctx, objID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
		}
		ticket = t
	} else {
		t, err := h.supportTicketRepo.GetTicketByTicketID(ctx, idParam)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
		}
		ticket = t
	}

	return c.JSON(ticket)
}

// ReplyToTicket adds a reply to a ticket (Admin reply expected)
func (h *AdminHandler) ReplyToTicket(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Params("id")
	// find ticket to get ObjectID
	var ticket *appmodels.SupportTicket
	if objID, err := primitive.ObjectIDFromHex(idParam); err == nil {
		t, err := h.supportTicketRepo.GetTicketByID(ctx, objID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
		}
		ticket = t
	} else {
		t, err := h.supportTicketRepo.GetTicketByTicketID(ctx, idParam)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
		}
		ticket = t
	}

	var req struct {
		Message string `json:"message" form:"message"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.Message == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Message is required"})
	}

	adminName := "Admin"
	if v := c.Locals("admin_name"); v != nil {
		if s, ok := v.(string); ok && s != "" {
			adminName = s
		}
	}

	reply := appmodels.TicketReply{Message: req.Message, IsAdmin: true, ReadByAdmin: true, AdminName: adminName, CreatedAt: time.Now()}

	if err := h.supportTicketRepo.AddReply(ctx, ticket.ID, reply); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add reply"})
	}
	return c.JSON(fiber.Map{"success": true})
}

// UpdateTicketStatus updates status for a ticket
func (h *AdminHandler) UpdateTicketStatus(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Params("id")
	var ticket *appmodels.SupportTicket
	if objID, err := primitive.ObjectIDFromHex(idParam); err == nil {
		t, err := h.supportTicketRepo.GetTicketByID(ctx, objID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
		}
		ticket = t
	} else {
		t, err := h.supportTicketRepo.GetTicketByTicketID(ctx, idParam)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
		}
		ticket = t
	}

	var req struct {
		Status string `json:"status" form:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.Status == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Status is required"})
	}

	if err := h.supportTicketRepo.UpdateTicketStatus(ctx, ticket.ID, req.Status); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update status"})
	}
	return c.JSON(fiber.Map{"success": true})
}

// UpdateTicketPriority updates priority for a ticket
func (h *AdminHandler) UpdateTicketPriority(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Params("id")
	var ticket *appmodels.SupportTicket
	if objID, err := primitive.ObjectIDFromHex(idParam); err == nil {
		t, err := h.supportTicketRepo.GetTicketByID(ctx, objID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
		}
		ticket = t
	} else {
		t, err := h.supportTicketRepo.GetTicketByTicketID(ctx, idParam)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Ticket not found"})
		}
		ticket = t
	}

	var req struct {
		Priority string `json:"priority" form:"priority"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.Priority == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Priority is required"})
	}

	if err := h.supportTicketRepo.UpdateTicketPriority(ctx, ticket.ID, req.Priority); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update priority"})
	}
	return c.JSON(fiber.Map{"success": true})
}
