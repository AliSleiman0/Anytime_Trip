package admin

import (
	"html/template"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// ServePage returns the admin landing page HTML
func (h *AdminHandler) ServePage(c *fiber.Ctx) error {
	// Serve the shared admin HTML template located under backend/templates
	return c.SendFile(filepath.Clean(filepath.Join("templates", "index.html")))
}

// GetSidebar serves the sidebar HTML fragment for HTMX partial loads
func (h *AdminHandler) GetSidebar(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "sidebar.html")))
}

// GetHeader serves the header HTML fragment for HTMX partial loads
func (h *AdminHandler) GetHeader(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "header.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	title := c.Query("title", "Dashboard")
	subtitle := c.Query("subtitle", "Here's what's happening today!")

	data := fiber.Map{
		"Title":    title,
		"Subtitle": subtitle,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}
