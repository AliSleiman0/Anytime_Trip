package admin

import (
	"html/template"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// GetPayments renders the payments page
func (h *AdminHandler) GetPayments(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "payments.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Payments and Transactions",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetPaymentsFragment serves the payments fragment for HTMX partial loads
func (h *AdminHandler) GetPaymentsFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "payments-frag.html")))
}
