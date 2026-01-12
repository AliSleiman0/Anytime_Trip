package admin

import (
	"html/template"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// GetReports renders the reports page
func (h *AdminHandler) GetReports(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "reports.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Reports and Analytics",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetReportsFragment serves the reports fragment for HTMX partial loads
func (h *AdminHandler) GetReportsFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "reports-frag.html")))
}
