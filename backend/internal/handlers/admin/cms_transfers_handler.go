package admin

import (
	"context"
	"html/template"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CMSTransfers renders the transfers management page
func (h *AdminHandler) CMSTransfers(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "cms-transfers.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Transfers Management",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetCMSTransfersFragment serves the transfers management fragment for HTMX partial loads
func (h *AdminHandler) GetCMSTransfersFragment(c *fiber.Ctx) error {
	ctx := context.Background()

	// Fetch all transfers from database
	transfers, err := h.transferRepo.FindAll(ctx, 100, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching transfers")
	}

	// Create template with custom functions
	funcMap := template.FuncMap{
		"formatDateTime": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02T15:04")
		},
		"formatTime": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("Mon, Jan 2, 3:04 PM")
		},
	}

	tmpl, err := template.New("cms-transfers-frag.html").Funcs(funcMap).ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-transfers-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Transfers": transfers,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// UpdateTransfer handles transfer update requests from the CMS modal
func (h *AdminHandler) UpdateTransfer(c *fiber.Ctx) error {
	ctx := context.Background()

	type TransferUpdateRequest struct {
		ID              string  `json:"id"`
		ServiceName     string  `json:"service_name"`
		PickupLocation  string  `json:"pickup_location"`
		PickupTime      string  `json:"pickup_time"`
		DropoffLocation string  `json:"dropoff_location"`
		DropoffTime     string  `json:"dropoff_time"`
		ProfitPercent   float64 `json:"profit_percent"`
		Refundable      bool    `json:"refundable"`
		AllowChanges    bool    `json:"allow_changes"`
		MeetAndGreet    bool    `json:"meet_and_greet"`
		LuggageIncluded bool    `json:"luggage_included"`
	}

	var req TransferUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Fetch existing transfer
	transfer, err := h.transferRepo.FindByID(ctx, req.ID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Transfer not found",
		})
	}

	// Update transfer fields
	transfer.ServiceName = req.ServiceName
	transfer.PickupLocation = req.PickupLocation
	transfer.DropoffLocation = req.DropoffLocation
	transfer.ProfitPercent = req.ProfitPercent
	transfer.Refundable = req.Refundable
	transfer.AllowChanges = req.AllowChanges
	transfer.MeetAndGreet = req.MeetAndGreet
	transfer.LuggageIncluded = req.LuggageIncluded

	// Parse and update pickup/dropoff times if provided
	if req.PickupTime != "" {
		if pickupTime, err := time.Parse(time.RFC3339, req.PickupTime); err == nil {
			transfer.PickupTime = pickupTime
		}
	}
	if req.DropoffTime != "" {
		if dropoffTime, err := time.Parse(time.RFC3339, req.DropoffTime); err == nil {
			transfer.DropoffTime = dropoffTime
		}
	}

	// Update in database
	if err := h.transferRepo.Update(ctx, req.ID, transfer); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update transfer",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Transfer updated successfully",
	})
}

// DeleteTransfer handles transfer deletion from CMS
func (h *AdminHandler) DeleteTransfer(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Query("id")

	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Transfer ID is required",
		})
	}

	if err := h.transferRepo.Delete(ctx, id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete transfer",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Transfer deleted successfully",
	})
}
