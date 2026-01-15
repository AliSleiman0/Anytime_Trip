package admin

import (
	"fmt"
	"math"
	"strings"

	adminrepo "Anytime_Travel/backend/internal/repository/admin"
	"Anytime_Travel/backend/internal/repository/app"
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
	adminRepo             *adminrepo.AdminRepository
	notificationPrefsRepo *adminrepo.NotificationPreferencesRepository
	userRepo              *app.UserRepository
	carBookingRepo        *app.CarBookingRepository
	flightBookingRepo     *app.FlightBookingRepository
	hotelBookingRepo      *app.HotelBookingRepository
	supportTicketRepo     *app.SupportTicketRepository
	flightRepo            *adminrepo.FlightRepository
	carRepo               *adminrepo.CarRepository
	hotelRepo             *adminrepo.HotelRepository
	jwtSecret             string
}

func NewAdminHandler(adminRepo *adminrepo.AdminRepository, notificationPrefsRepo *adminrepo.NotificationPreferencesRepository, userRepo *app.UserRepository, carBookingRepo *app.CarBookingRepository, flightBookingRepo *app.FlightBookingRepository, hotelBookingRepo *app.HotelBookingRepository, supportTicketRepo *app.SupportTicketRepository, flightRepo *adminrepo.FlightRepository, carRepo *adminrepo.CarRepository, hotelRepo *adminrepo.HotelRepository, jwtSecret string) *AdminHandler {
	return &AdminHandler{
		adminRepo:             adminRepo,
		notificationPrefsRepo: notificationPrefsRepo,
		userRepo:              userRepo,
		carBookingRepo:        carBookingRepo,
		flightBookingRepo:     flightBookingRepo,
		hotelBookingRepo:      hotelBookingRepo,
		supportTicketRepo:     supportTicketRepo,
		flightRepo:            flightRepo,
		carRepo:               carRepo,
		hotelRepo:             hotelRepo,
		jwtSecret:             jwtSecret,
	}
}
