package admin

import (
	"context"
	"html/template"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// CMSHotels renders the hotels management page
func (h *AdminHandler) CMSHotels(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "cms-hotels.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Hotels Management",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetCMSHotelsFragment serves the hotels management fragment for HTMX partial loads
func (h *AdminHandler) GetCMSHotelsFragment(c *fiber.Ctx) error {
	ctx := context.Background()

	// Fetch all hotels from database
	hotels, err := h.hotelRepo.FindAll(ctx, 100, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching hotels")
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-hotels-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Hotels": hotels,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}
