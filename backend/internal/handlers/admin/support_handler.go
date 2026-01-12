package admin

import (
	"html/template"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
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

// GetSupportFragment serves the support fragment for HTMX partial loads
func (h *AdminHandler) GetSupportFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "support-frag.html")))
}
