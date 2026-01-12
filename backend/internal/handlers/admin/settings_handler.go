package admin

import (
	"html/template"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// GetSettings renders the settings page
func (h *AdminHandler) GetSettings(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "settings.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Settings",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetSettingsFragment serves the settings fragment for HTMX partial loads
func (h *AdminHandler) GetSettingsFragment(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get admin ID from context
	adminID := c.Locals("user_id")
	if adminID == nil {
		return c.Status(401).SendString("Unauthorized")
	}

	// Convert adminID to string
	adminIDStr, ok := adminID.(string)
	if !ok {
		return c.Status(500).SendString("Invalid admin ID format")
	}

	// Fetch admin data
	admin, err := h.adminRepo.FindByID(ctx, adminIDStr)
	if err != nil {
		// If admin not found, show basic page without data
		tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "settings-frag.html")))
		if err != nil {
			return c.Status(500).SendString("Error loading template")
		}

		data := fiber.Map{
			"Name":     "Admin",
			"Email":    "",
			"Password": "••••••••••••••",
			"Photo":    "",
			"Currency": "USD",
			"Language": "English",
		}

		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c, data)
	}

	// Parse and execute template with admin data
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "settings-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Name":     admin.Username,
		"Email":    admin.Email,
		"Password": "••••••••••••••", // Never expose real password
		"Photo":    "",               // Add photo URL if available in your model
		"Currency": admin.Currency,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetNotificationPreferences gets notification preferences for the logged-in admin
func (h *AdminHandler) GetNotificationPreferences(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get admin ID from session/JWT (you'll need to implement this based on your auth)
	adminID := c.Locals("user_id")
	if adminID == nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	prefs, err := h.notificationPrefsRepo.FindByAdminID(ctx, adminID.(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve notification preferences",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    prefs,
	})
}

// ChangeLanguage handles admin language change
func (h *AdminHandler) ChangeLanguage(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get admin ID from JWT
	adminID := c.Locals("user_id")
	if adminID == nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Parse request body
	var req struct {
		Language string `json:"language"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Validate language
	if req.Language != "English" && req.Language != "Arabic" && req.Language != "French" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid language. Allowed values: English, Arabic, French",
		})
	}

	// Convert admin ID to string
	adminIDStr, ok := adminID.(string)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Invalid admin ID format",
		})
	}

	// Update language in repository
	err := h.adminRepo.UpdateLanguage(ctx, adminIDStr, req.Language)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Language updated successfully",
	})
}

// ChangeCurrency handles admin currency change
func (h *AdminHandler) ChangeCurrency(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get admin ID from JWT
	adminID := c.Locals("user_id")
	if adminID == nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Parse request body
	var req struct {
		Currency string `json:"currency"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Validate currency
	if req.Currency != "USD" && req.Currency != "LBP" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid currency. Allowed values: USD, LBP",
		})
	}

	// Convert admin ID to string
	adminIDStr, ok := adminID.(string)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Invalid admin ID format",
		})
	}

	// Update currency in repository
	err := h.adminRepo.UpdateCurrency(ctx, adminIDStr, req.Currency)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Currency updated successfully",
	})
}

// ChangePassword handles admin password change
func (h *AdminHandler) ChangePassword(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get admin ID from JWT
	adminID := c.Locals("user_id")
	if adminID == nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	// Parse request body
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Validate input
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Current password and new password are required",
		})
	}

	if len(req.NewPassword) < 8 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "New password must be at least 8 characters long",
		})
	}

	// Convert admin ID to string
	adminIDStr, ok := adminID.(string)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Invalid admin ID format",
		})
	}

	// Update password in repository
	err := h.adminRepo.UpdatePassword(ctx, adminIDStr, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Password updated successfully",
	})
}

// UpdateNotificationPreferences updates notification preferences for the logged-in admin
func (h *AdminHandler) UpdateNotificationPreferences(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get admin ID from session/JWT
	adminID := c.Locals("user_id")
	if adminID == nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	type UpdatePrefsRequest struct {
		Bookings map[string]bool `json:"bookings"`
		Payments map[string]bool `json:"payments"`
	}

	var req UpdatePrefsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Get current preferences
	prefs, err := h.notificationPrefsRepo.FindByAdminID(ctx, adminID.(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve current preferences",
		})
	}

	// Update bookings if provided
	if req.Bookings != nil {
		if val, ok := req.Bookings["new_booking"]; ok {
			prefs.Bookings.NewBooking = val
		}
		if val, ok := req.Bookings["booking_cancelled"]; ok {
			prefs.Bookings.BookingCancelled = val
		}
		if val, ok := req.Bookings["booking_confirmation"]; ok {
			prefs.Bookings.BookingConfirmation = val
		}
		if val, ok := req.Bookings["booking_rescheduled"]; ok {
			prefs.Bookings.BookingRescheduled = val
		}
		if val, ok := req.Bookings["booking_reminder"]; ok {
			prefs.Bookings.BookingReminder = val
		}
		if val, ok := req.Bookings["booking_updated"]; ok {
			prefs.Bookings.BookingUpdated = val
		}
	}

	// Update payments if provided
	if req.Payments != nil {
		if val, ok := req.Payments["payments_received"]; ok {
			prefs.Payments.PaymentsReceived = val
		}
		if val, ok := req.Payments["payments_confirmed"]; ok {
			prefs.Payments.PaymentsConfirmed = val
		}
		if val, ok := req.Payments["payment_received"]; ok {
			prefs.Payments.PaymentReceived = val
		}
		if val, ok := req.Payments["booking_rescheduled"]; ok {
			prefs.Payments.BookingRescheduled = val
		}
		if val, ok := req.Payments["payment_pending"]; ok {
			prefs.Payments.PaymentPending = val
		}
	}

	// Save preferences
	if err := h.notificationPrefsRepo.Upsert(ctx, prefs); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update notification preferences",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification preferences updated successfully",
		"data":    prefs,
	})
}
