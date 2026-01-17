package admin

import (
	"html/template"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// GetDashboard renders the dashboard page shell.
func (h *AdminHandler) GetDashboard(c *fiber.Ctx) error {
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

// GetDashboardFragment serves the dashboard fragment for HTMX partial loads.
func (h *AdminHandler) GetDashboardFragment(c *fiber.Ctx) error {
	totalUsers, err := h.userRepo.CountTotalUsers(c.Context())
	if err != nil {
		return c.Status(500).SendString("Error fetching total users")
	}

	newUsersToday, err := h.userRepo.CountNewUsersToday(c.Context())
	if err != nil {
		return c.Status(500).SendString("Error fetching new users today")
	}

	activeUsers, err := h.userRepo.CountActiveUsers(c.Context())
	if err != nil {
		return c.Status(500).SendString("Error fetching active users")
	}

	unactiveUsers := totalUsers - activeUsers

	data := fiber.Map{
		"TotalUsers":           totalUsers,
		"NewUsersToday":        newUsersToday,
		"ActiveUsers":          activeUsers,
		"UnactiveUsers":        unactiveUsers,
		"Visits":               71000, // placeholder
		"NewUsersPercent":      "+11.01%",
		"ActiveUsersPercent":   "+11.01%",
		"UnactiveUsersPercent": "+11.01%",
		"VisitsPercent":        "-11.01%",
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "dashboard-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}
