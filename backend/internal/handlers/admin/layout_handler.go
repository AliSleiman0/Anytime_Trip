package admin

import (
	"html/template"
	"path/filepath"
	"strings"

	"Anytime_Travel/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// ServePage returns the admin landing page HTML
func (h *AdminHandler) ServePage(c *fiber.Ctx) error {
	// Serve the shared admin HTML template located under backend/templates
	return c.SendFile(filepath.Clean(filepath.Join("templates", "index.html")))
}

// GetSidebar serves the sidebar HTML fragment for HTMX partial loads
func (h *AdminHandler) GetSidebar(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "sidebar.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"IsSuperAdmin": middleware.IsSuperAdmin(c),
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetHeader serves the header HTML fragment for HTMX partial loads
func (h *AdminHandler) GetHeader(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "header.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	title := c.Query("title", "Dashboard")
	subtitle := c.Query("subtitle", "Here's what's happening today!")

	// Only show export on specific admin pages
	titleKey := strings.ToLower(strings.TrimSpace(title))
	showExport := false
	switch titleKey {
	case "bookings", "user management", "service management", "payments and transactions":
		showExport = true
	}

	data := fiber.Map{
		"Title":            title,
		"Subtitle":         subtitle,
		"ShowExportButton": showExport,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}
