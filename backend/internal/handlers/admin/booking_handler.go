package admin

import (
	"Anytime_Travel/backend/internal/repository/admin"
	"Anytime_Travel/backend/internal/repository/app"
	"fmt"
	"html/template"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// BookingRow is a flattened view model combining all booking types for the admin table.
type BookingRow struct {
	BookingID     string
	RefID         string
	Type          string
	CustomerName  string
	CustomerEmail string
	BookingDate   time.Time
	TravelDate    time.Time
	Details       string
	Amount        string
	PaymentStatus string
	PaymentClass  string
	Status        string
	StatusClass   string
}

// NewBookingHandler returns the existing AdminHandler with repos injected.
// NOTE: AdminHandler already holds the repos; this helper preserves constructor parity if needed.
func NewBookingHandler(adminRepo *admin.AdminRepository, userRepo *app.UserRepository, carBookingRepo *app.CarBookingRepository, flightBookingRepo *app.FlightBookingRepository, hotelBookingRepo *app.HotelBookingRepository) *AdminHandler {
	return &AdminHandler{
		adminRepo:         adminRepo,
		userRepo:          userRepo,
		carBookingRepo:    carBookingRepo,
		flightBookingRepo: flightBookingRepo,
		hotelBookingRepo:  hotelBookingRepo,
	}
}

// GetBookings renders the bookings page shell.
func (h *AdminHandler) GetBookings(c *fiber.Ctx) error {
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

// GetBookingsFragment serves the bookings fragment for HTMX partial loads.
func (h *AdminHandler) GetBookingsFragment(c *fiber.Ctx) error {
	ctx := c.Context()

	carBookings, err := h.carBookingRepo.FindAll(ctx, 100, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching car bookings")
	}

	flightBookings, err := h.flightBookingRepo.FindAll(ctx, 100, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching flight bookings")
	}

	hotelBookings, err := h.hotelBookingRepo.FindAll(ctx, 100, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching hotel bookings")
	}

	rows := make([]BookingRow, 0, len(carBookings)+len(flightBookings)+len(hotelBookings))

	typeFilter := strings.ToLower(c.Query("bookings", "all"))
	statusFilter := strings.ToLower(c.Query("status", "all"))
	dateFilter := c.Query("date", "")

	formatAmount := func(currency string, amount float64) string {
		return fmt.Sprintf("%s %.2f", currency, amount)
	}

	badgeClass := func(kind string, value string) string {
		switch kind {
		case "payment":
			switch value {
			case "paid":
				return "px-2 py-1 bg-[#26CE0033] text-[#26CE00] rounded-full text-xs font-semibold"
			case "pending", "by_whish":
				return "px-2 py-1 bg-[#FFD966] text-[#8A6D00] rounded-full text-xs font-semibold"
			case "refunded", "cancelled":
				return "px-2 py-1 bg-[#D2412433] text-[#D24124] rounded-full text-xs font-semibold"
			default:
				return "px-2 py-1 bg-gray-200 text-gray-700 rounded-full text-xs font-semibold"
			}
		case "status":
			switch value {
			case "confirmed", "completed":
				return "px-2 py-1 bg-[#26CE0033] text-[#26CE00] rounded-full text-xs font-semibold"
			case "pending":
				return "px-2 py-1 bg-[#FFD966] text-[#8A6D00] rounded-full text-xs font-semibold"
			case "cancelled":
				return "px-2 py-1 bg-[#D2412433] text-[#D24124] rounded-full text-xs font-semibold"
			default:
				return "px-2 py-1 bg-gray-200 text-gray-700 rounded-full text-xs font-semibold"
			}
		default:
			return "px-2 py-1 bg-gray-200 text-gray-700 rounded-full text-xs font-semibold"
		}
	}

	matchesFilters := func(row BookingRow) bool {
		if typeFilter != "" && typeFilter != "all" && strings.ToLower(row.Type) != typeFilter {
			return false
		}

		if statusFilter != "" && statusFilter != "all" && strings.ToLower(row.Status) != statusFilter {
			return false
		}

		if dateFilter != "" {
			if row.BookingDate.IsZero() || row.BookingDate.Format("2006-01-02") != dateFilter {
				return false
			}
		}

		return true
	}

	for _, b := range carBookings {
		row := BookingRow{
			BookingID:     b.BookingID,
			RefID:         b.ID,
			Type:          "Car",
			CustomerName:  b.Customer.Name,
			CustomerEmail: b.Customer.Email,
			BookingDate:   b.BookingDate,
			TravelDate:    b.BookingDate,
			Details:       b.Details,
			Amount:        formatAmount(b.Currency, b.Amount),
			PaymentStatus: string(b.PaymentStatus),
			PaymentClass:  badgeClass("payment", string(b.PaymentStatus)),
			Status:        string(b.Status),
			StatusClass:   badgeClass("status", string(b.Status)),
		}

		if matchesFilters(row) {
			rows = append(rows, row)
		}
	}

	for _, b := range flightBookings {
		row := BookingRow{
			BookingID:     b.BookingID,
			RefID:         b.ID,
			Type:          "Flight",
			CustomerName:  b.Customer.Name,
			CustomerEmail: b.Customer.Email,
			BookingDate:   b.BookingDate,
			TravelDate:    b.BookingDate,
			Details:       b.Details,
			Amount:        formatAmount(b.Currency, b.Amount),
			PaymentStatus: string(b.PaymentStatus),
			PaymentClass:  badgeClass("payment", string(b.PaymentStatus)),
			Status:        string(b.Status),
			StatusClass:   badgeClass("status", string(b.Status)),
		}

		if matchesFilters(row) {
			rows = append(rows, row)
		}
	}

	for _, b := range hotelBookings {
		row := BookingRow{
			BookingID:     b.BookingID,
			RefID:         b.ID,
			Type:          "Hotel",
			CustomerName:  b.Customer.Name,
			CustomerEmail: b.Customer.Email,
			BookingDate:   b.BookingDate,
			TravelDate:    b.BookingDate,
			Details:       b.Details,
			Amount:        formatAmount(b.Currency, b.Amount),
			PaymentStatus: string(b.PaymentStatus),
			PaymentClass:  badgeClass("payment", string(b.PaymentStatus)),
			Status:        string(b.Status),
			StatusClass:   badgeClass("status", string(b.Status)),
		}

		if matchesFilters(row) {
			rows = append(rows, row)
		}
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].BookingDate.After(rows[j].BookingDate)
	})

	filters := struct {
		Date   string
		Type   string
		Status string
	}{
		Date:   dateFilter,
		Type:   typeFilter,
		Status: statusFilter,
	}

	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			if t.IsZero() {
				return "-"
			}
			return t.Format("01/02/2006")
		},
	}

	tmpl, err := template.New("bookings-frag.html").Funcs(funcMap).ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "bookings-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Bookings": rows,
		"Filters":  filters,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}
