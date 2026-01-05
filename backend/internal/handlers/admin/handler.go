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

// GetBookingsFragment serves the bookings fragment for HTMX partial loads
func (h *AdminHandler) GetBookingsFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "bookings-frag.html")))
}

func (h *AdminHandler) ManageUsers(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "user.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Users",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetUsersFragment serves the user management fragment for HTMX partial loads
func (h *AdminHandler) GetUsersFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "user-frag.html")))
}

func (h *AdminHandler) ViewUser(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "view-user.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "View User",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetViewUserFragment serves the view user fragment for HTMX partial loads
func (h *AdminHandler) GetViewUserFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "view-user-frag.html")))
}

func (h *AdminHandler) ManageServiceProviders(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "service-provider.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Service Providers",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetServiceProvidersFragment serves the service provider management fragment for HTMX partial loads
func (h *AdminHandler) GetServiceProvidersFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "service-provider-frag.html")))
}

func (h *AdminHandler) ViewService(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "view-service.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "View Service Provider",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetViewServiceFragment serves the view service provider fragment for HTMX partial loads
func (h *AdminHandler) GetViewServiceFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "view-service-frag.html")))
}

func (h *AdminHandler) CMSFlights(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "cms-flights.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Flights Management",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetCMSFlightsFragment serves the flights management fragment for HTMX partial loads
func (h *AdminHandler) GetCMSFlightsFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-flights-frag.html")))
}

func (h *AdminHandler) CMSCars(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "cms-cars.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Cars Management",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetCMSCarsFragment serves the cars management fragment for HTMX partial loads
func (h *AdminHandler) GetCMSCarsFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-cars-frag.html")))
}

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
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-hotels-frag.html")))
}

func (h *AdminHandler) CMSTravelExperience(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "cms-home-travel.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Homepage - Travel Experience",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetCMSTravelExperienceFragment serves the travel experience fragment for HTMX partial loads
func (h *AdminHandler) GetCMSTravelExperienceFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-travel-frag.html")))
}

func (h *AdminHandler) CMSBanner(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "cms-home-banner.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Homepage - Banner",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetCMSBannerFragment serves the banner fragment for HTMX partial loads
func (h *AdminHandler) GetCMSBannerFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-banner-frag.html")))
}

func (h *AdminHandler) CMSPopularLocations(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "cms-home-popular.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Homepage - Popular Locations",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetCMSPopularLocationsFragment serves the popular locations fragment for HTMX partial loads
func (h *AdminHandler) GetCMSPopularLocationsFragment(c *fiber.Ctx) error {
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-home-popular-frag.html")))
}

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
	return c.SendFile(filepath.Clean(filepath.Join("templates", "admin", "fragments", "settings-frag.html")))
}

// ServePage returns the admin landing page HTML
func (h *AdminHandler) ServePage(c *fiber.Ctx) error {
	// Serve the shared admin HTML template located under backend/templates
	return c.SendFile(filepath.Clean(filepath.Join("templates", "index.html")))
}
