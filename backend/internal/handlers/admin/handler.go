package admin

import (
	"fmt"
	"math"
	"strings"

	adminrepo "travel/backend/internal/repository/admin"
	"travel/backend/internal/repository/app"
	superadminrepo "travel/backend/internal/repository/superadmin"
	"travel/backend/internal/ws"
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
	IsFreezed     bool
}

// BookingView represents a booking entry for the provider detail page
type BookingView struct {
	UserName    string
	BookingType string
	Status      string
	StatusClass string
	Destination string
	Rating      string
	Amount      string
	BookingDate string
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
	IsFreezed              bool
	FreezeActionText       string
	Bookings               []BookingView
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
	passwordResetRepo     *adminrepo.PasswordResetRepository
	userRepo              *app.UserRepository
	carBookingRepo        *app.CarBookingRepository
	flightBookingRepo     *app.FlightBookingRepository
	hotelBookingRepo      *app.HotelBookingRepository
	transferBookingRepo   *app.TransferBookingRepository
	supportTicketRepo     *app.SupportTicketRepository
	paymentRepo           *app.PaymentRepository
	flightRepo            *adminrepo.FlightRepository
	carRepo               *adminrepo.CarRepository
	hotelRepo             *adminrepo.HotelRepository
	transferRepo          *adminrepo.TransferRepository
	travelRepo            *adminrepo.TravelRepository
	bannerRepo            *adminrepo.BannerRepository
	popularRepo           *adminrepo.PopularRepository
	predefinedAnswerRepo  *superadminrepo.PredefinedAnswerRepository
	jwtSecret             string
	chatHub               *ws.Hub
	notifyHub             *ws.AdminHub
	loginAttemptRepo      *adminrepo.LoginAttemptRepository
}

func NewAdminHandler(adminRepo *adminrepo.AdminRepository, notificationPrefsRepo *adminrepo.NotificationPreferencesRepository, passwordResetRepo *adminrepo.PasswordResetRepository, userRepo *app.UserRepository, carBookingRepo *app.CarBookingRepository, flightBookingRepo *app.FlightBookingRepository, hotelBookingRepo *app.HotelBookingRepository, transferBookingRepo *app.TransferBookingRepository, supportTicketRepo *app.SupportTicketRepository, paymentRepo *app.PaymentRepository, flightRepo *adminrepo.FlightRepository, carRepo *adminrepo.CarRepository, hotelRepo *adminrepo.HotelRepository, transferRepo *adminrepo.TransferRepository, bannerRepo *adminrepo.BannerRepository, travelRepo *adminrepo.TravelRepository, popularRepo *adminrepo.PopularRepository, predefinedAnswerRepo *superadminrepo.PredefinedAnswerRepository, loginAttemptRepo *adminrepo.LoginAttemptRepository, jwtSecret string, chatHub *ws.Hub, notifyHub *ws.AdminHub) *AdminHandler {
	return &AdminHandler{
		adminRepo:             adminRepo,
		notificationPrefsRepo: notificationPrefsRepo,
		passwordResetRepo:     passwordResetRepo,
		userRepo:              userRepo,
		carBookingRepo:        carBookingRepo,
		flightBookingRepo:     flightBookingRepo,
		hotelBookingRepo:      hotelBookingRepo,
		transferBookingRepo:   transferBookingRepo,
		supportTicketRepo:     supportTicketRepo,
		paymentRepo:           paymentRepo,
		flightRepo:            flightRepo,
		carRepo:               carRepo,
		hotelRepo:             hotelRepo,
		transferRepo:          transferRepo,
		bannerRepo:            bannerRepo,
		travelRepo:            travelRepo,
		popularRepo:           popularRepo,
		predefinedAnswerRepo:  predefinedAnswerRepo,
		loginAttemptRepo:      loginAttemptRepo,
		jwtSecret:             jwtSecret,
		chatHub:               chatHub,
		notifyHub:             notifyHub,
	}
}
