package admin

import (
	"Anytime_Travel/backend/internal/repository/admin"
	"html/template"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// AdminHandler handles admin-level requests
type AdminHandler struct {
	adminRepo *admin.AdminRepository
}

func NewAdminHandler(adminRepo *admin.AdminRepository) *AdminHandler {
	return &AdminHandler{
		adminRepo: adminRepo,
	}
}

func (h *AdminHandler) GetDashboard(c *fiber.Ctx) error {
	// Parse and render dashboard template with data
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "dashboard.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Dashboard",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetDashboardFragment serves the dashboard fragment for HTMX partial loads
func (h *AdminHandler) GetDashboardFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "dashboard-frag.html")))
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
func (h *AdminHandler) GetBookings(c *fiber.Ctx) error {
	// Parse and render bookings template with data
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "bookings.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Bookings",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetBookingsFragment serves the bookings fragment for HTMX partial loads
func (h *AdminHandler) GetBookingsFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "bookings-frag.html")))
}

func (h *AdminHandler) ManageUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Admin User Management",
	})
}

// ServePage returns the admin landing page HTML
func (h *AdminHandler) ServePage(c *fiber.Ctx) error {
	// Serve the shared admin HTML template located under backend/templates
	return c.SendFile(filepath.Clean(filepath.Join("templates", "index.html")))
}
