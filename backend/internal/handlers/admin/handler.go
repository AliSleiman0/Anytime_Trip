package admin

import (
	"context"
	"fmt"
	"html/template"
	"math"
	"path/filepath"
	"strings"

	adminmodels "Anytime_Travel/backend/internal/models/admin"
	adminrepo "Anytime_Travel/backend/internal/repository/admin"
	"Anytime_Travel/backend/internal/repository/app"

	"github.com/gofiber/fiber/v2"
)

// ServiceProviderView is a flattened view model for the table
type ServiceProviderView struct {
	ID            string
	Name          string
	Type          string
	Status        string
	Location      string
	ContactEmail  string
	Rating        float64
	TotalBookings int64
	Revenue       float64
	ProviderType  string // flight | car | hotel for filtering
}

// ProviderDetailView is a simplified view model for the detail screen
type ProviderDetailView struct {
	ID                     string
	Name                   string
	Email                  string
	Status                 string
	StatusClass            string
	Type                   string
	Location               string
	CreatedAt              string
	Rating                 string
	TotalBookings          int64
	Revenue                string
	ProviderType           string
	ProfitShareDisplay     string
	PhoneNumberDisplay     string
	RevenueThisMonth       string
	TotalBookingsThisMonth string
}

func statusBadgeClass(status string) string {
	switch strings.ToLower(status) {
	case "active":
		return "bg-[#DCFCE7] text-[#008236]"
	case "inactive":
		return "bg-[#FFE4E6] text-[#DC2626]"
	default:
		return "bg-gray-200 text-gray-700"
	}
}

// formatCurrency formats a float64 amount into a currency string with thousand separators.
func formatCurrency(currency string, amount float64) string {
	if currency == "" {
		currency = "$"
	}

	value := int64(math.Round(amount))
	s := fmt.Sprintf("%d", value)
	n := len(s)
	if n <= 3 {
		return fmt.Sprintf("%s %s", currency, s)
	}

	var b strings.Builder
	rem := n % 3
	if rem == 0 {
		rem = 3
	}
	b.WriteString(s[:rem])
	for i := rem; i < n; i += 3 {
		b.WriteString(",")
		b.WriteString(s[i : i+3])
	}

	return fmt.Sprintf("%s %s", currency, b.String())
}

// formatProfitPercent prefers an explicit percent; if missing, derives from profit/revenue.
func formatProfitPercent(percent float64, profit float64, revenue float64) string {
	if percent > 0 {
		return fmt.Sprintf("%.1f%%", percent)
	}

	if revenue > 0 && profit > 0 {
		derived := (profit / revenue) * 100
		return fmt.Sprintf("%.1f%%", derived)
	}

	return "--"
}

// AdminHandler handles admin-level requests
type AdminHandler struct {
	adminRepo         *adminrepo.AdminRepository
	userRepo          *app.UserRepository
	carBookingRepo    *app.CarBookingRepository
	flightBookingRepo *app.FlightBookingRepository
	hotelBookingRepo  *app.HotelBookingRepository
	flightRepo        *adminrepo.FlightRepository
	carRepo           *adminrepo.CarRepository
	hotelRepo         *adminrepo.HotelRepository
}

func NewAdminHandler(adminRepo *adminrepo.AdminRepository, userRepo *app.UserRepository, carBookingRepo *app.CarBookingRepository, flightBookingRepo *app.FlightBookingRepository, hotelBookingRepo *app.HotelBookingRepository, flightRepo *adminrepo.FlightRepository, carRepo *adminrepo.CarRepository, hotelRepo *adminrepo.HotelRepository) *AdminHandler {
	return &AdminHandler{
		adminRepo:         adminRepo,
		userRepo:          userRepo,
		carBookingRepo:    carBookingRepo,
		flightBookingRepo: flightBookingRepo,
		hotelBookingRepo:  hotelBookingRepo,
		flightRepo:        flightRepo,
		carRepo:           carRepo,
		hotelRepo:         hotelRepo,
	}
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

func (h *AdminHandler) ViewUser(c *fiber.Ctx) error {
	userID := c.Query("id")
	if strings.TrimSpace(userID) == "" {
		return c.Status(400).SendString("User id is required")
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "view-user.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title":  "View User",
		"UserID": userID,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
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
	ctx := context.Background()

	flights, _ := h.flightRepo.FindAll(ctx, 0, 0)
	cars, _ := h.carRepo.FindAll(ctx, 0, 0)
	hotels, _ := h.hotelRepo.FindAll(ctx, 0, 0)

	providers := make([]ServiceProviderView, 0, len(flights)+len(cars)+len(hotels))

	for _, f := range flights {
		providers = append(providers, ServiceProviderView{
			ID:            f.ID,
			Name:          f.ProviderName,
			Type:          f.ProviderType,
			Status:        string(f.Status),
			Location:      f.Location,
			ContactEmail:  f.ContactEmail,
			Rating:        f.Rating,
			TotalBookings: f.TotalBookings,
			Revenue:       f.Revenue,
			ProviderType:  "flight",
		})
	}

	for _, car := range cars {
		providers = append(providers, ServiceProviderView{
			ID:            car.ID,
			Name:          car.ProviderName,
			Type:          car.ProviderType,
			Status:        string(car.Status),
			Location:      car.Location,
			ContactEmail:  car.ContactEmail,
			Rating:        car.Rating,
			TotalBookings: car.TotalBookings,
			Revenue:       car.Revenue,
			ProviderType:  "car",
		})
	}

	for _, htl := range hotels {
		providers = append(providers, ServiceProviderView{
			ID:            htl.ID,
			Name:          htl.ProviderName,
			Type:          htl.ProviderType,
			Status:        string(htl.Status),
			Location:      htl.Location,
			ContactEmail:  htl.ContactEmail,
			Rating:        htl.Rating,
			TotalBookings: htl.TotalBookings,
			Revenue:       htl.Revenue,
			ProviderType:  "hotel",
		})
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "service-provider-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Providers": providers,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

func (h *AdminHandler) ViewService(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Query("id"))
	providerType := strings.ToLower(strings.TrimSpace(c.Query("type")))

	if id == "" || providerType == "" {
		return c.Status(400).SendString("provider id and type are required")
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "view-service.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title":        "View Service Provider",
		"ProviderID":   id,
		"ProviderType": providerType,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetViewServiceFragment serves the view service provider fragment for HTMX partial loads
func (h *AdminHandler) GetViewServiceFragment(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Query("id"))
	providerType := strings.ToLower(strings.TrimSpace(c.Query("type")))

	if id == "" || providerType == "" {
		return c.Status(400).SendString("provider id and type are required")
	}

	ctx := context.Background()

	var detail ProviderDetailView

	switch providerType {
	case "flight":
		flight, err := h.flightRepo.FindByID(ctx, id)
		if err != nil {
			return c.Status(404).SendString("flight provider not found")
		}

		currency := flight.Currency
		if currency == "" {
			currency = "$"
		}

		detail = ProviderDetailView{
			ID:                     flight.ID,
			Name:                   flight.ProviderName,
			Email:                  flight.ContactEmail,
			Status:                 string(flight.Status),
			StatusClass:            statusBadgeClass(string(flight.Status)),
			Type:                   flight.ProviderType,
			Location:               flight.Location,
			CreatedAt:              flight.CreatedAt.Format("02 Jan, 2006"),
			Rating:                 fmt.Sprintf("%.1f", flight.Rating),
			TotalBookings:          flight.TotalBookings,
			Revenue:                formatCurrency(currency, flight.Revenue),
			ProviderType:           providerType,
			ProfitShareDisplay:     formatProfitPercent(flight.ProfitPercent, flight.Profit, flight.Revenue),
			PhoneNumberDisplay:     "N/A",
			RevenueThisMonth:       "--",
			TotalBookingsThisMonth: "--",
		}
	case "car":
		car, err := h.carRepo.FindByID(ctx, id)
		if err != nil {
			return c.Status(404).SendString("car provider not found")
		}

		currency := car.Currency
		if currency == "" {
			currency = "$"
		}

		detail = ProviderDetailView{
			ID:                     car.ID,
			Name:                   car.ProviderName,
			Email:                  car.ContactEmail,
			Status:                 string(car.Status),
			StatusClass:            statusBadgeClass(string(car.Status)),
			Type:                   car.ProviderType,
			Location:               car.Location,
			CreatedAt:              car.CreatedAt.Format("02 Jan, 2006"),
			Rating:                 fmt.Sprintf("%.1f", car.Rating),
			TotalBookings:          car.TotalBookings,
			Revenue:                formatCurrency(currency, car.Revenue),
			ProviderType:           providerType,
			ProfitShareDisplay:     formatProfitPercent(car.ProfitPercent, car.Profit, car.Revenue),
			PhoneNumberDisplay:     "N/A",
			RevenueThisMonth:       "--",
			TotalBookingsThisMonth: "--",
		}
	case "hotel":
		hotel, err := h.hotelRepo.FindByID(ctx, id)
		if err != nil {
			return c.Status(404).SendString("hotel provider not found")
		}

		currency := hotel.Currency
		if currency == "" {
			currency = "$"
		}

		detail = ProviderDetailView{
			ID:                     hotel.ID,
			Name:                   hotel.ProviderName,
			Email:                  hotel.ContactEmail,
			Status:                 string(hotel.Status),
			StatusClass:            statusBadgeClass(string(hotel.Status)),
			Type:                   hotel.ProviderType,
			Location:               hotel.Location,
			CreatedAt:              hotel.CreatedAt.Format("02 Jan, 2006"),
			Rating:                 fmt.Sprintf("%.1f", hotel.Rating),
			TotalBookings:          hotel.TotalBookings,
			Revenue:                formatCurrency(currency, hotel.Revenue),
			ProviderType:           providerType,
			ProfitShareDisplay:     formatProfitPercent(hotel.ProfitPercent, hotel.Profit, hotel.Revenue),
			PhoneNumberDisplay:     "N/A",
			RevenueThisMonth:       "--",
			TotalBookingsThisMonth: "--",
		}
	default:
		return c.Status(400).SendString("invalid provider type")
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "view-service-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Provider": detail,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
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
	ctx := context.Background()

	// Fetch all flights from database
	flights, err := h.flightRepo.FindAll(ctx, 100, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching flights")
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-flights-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Flights": flights,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// UpdateFlight handles flight update requests from the CMS modal
func (h *AdminHandler) UpdateFlight(c *fiber.Ctx) error {
	ctx := context.Background()

	type FlightUpdateRequest struct {
		ID              string  `json:"id"`
		FlightNumber    string  `json:"flight_number"`
		FlightID        string  `json:"flight_id"`
		Status          string  `json:"status"`
		Cost            float64 `json:"cost"`
		ProfitPercent   float64 `json:"profit_percent"`
		Refundable      bool    `json:"refundable"`
		AllowChanges    bool    `json:"allow_changes"`
		AllowSeatChoice bool    `json:"allow_seat_choice"`
		AllowCarryOn    bool    `json:"allow_carry_on"`
	}

	var req FlightUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Fetch existing flight
	flight, err := h.flightRepo.FindByID(ctx, req.ID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Flight not found",
		})
	}

	// Update flight fields
	flight.FlightNumber = req.FlightNumber
	flight.FlightID = req.FlightID
	flight.Status = adminmodels.FlightStatus(req.Status)
	flight.Cost = req.Cost
	flight.ProfitPercent = req.ProfitPercent

	// Update in database
	if err := h.flightRepo.Update(ctx, req.ID, flight); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update flight",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Flight updated successfully",
	})
}

// DeleteFlight handles flight deletion from CMS
func (h *AdminHandler) DeleteFlight(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Query("id")

	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Flight ID is required",
		})
	}

	if err := h.flightRepo.Delete(ctx, id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete flight",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Flight deleted successfully",
	})
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
	ctx := context.Background()

	// Fetch all cars from database
	cars, err := h.carRepo.FindAll(ctx, 100, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching cars")
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-cars-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Cars": cars,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
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
