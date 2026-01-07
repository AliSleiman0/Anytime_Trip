package admin

import (
	"fmt"
	"html/template"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// UserRow is a view model for the user table.
type UserRow struct {
	ID            string
	Name          string
	Email         string
	Phone         string
	Status        string
	StatusClass   string
	TotalBookings string
	LastBooking   string
	CreatedOn     string
}

// UserBookingRow is a view model for a single booking line in the user details view.
type UserBookingRow struct {
	Provider    string
	Type        string
	Status      string
	StatusClass string
	Location    string
	Rating      string
	Revenue     string
	CreatedAt   time.Time
}

// GetUsersFragment serves the user management fragment for HTMX partial loads.
func (h *AdminHandler) GetUsersFragment(c *fiber.Ctx) error {
	ctx := c.Context()

	users, err := h.userRepo.FindAll(ctx, 200, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching users")
	}

	formatDate := func(t time.Time) string {
		if t.IsZero() {
			return "-"
		}
		return t.Format("02 Jan, 2006")
	}

	rows := make([]UserRow, 0, len(users))
	for _, u := range users {
		status := "Inactive"
		statusClass := "px-3 py-1 bg-[#FFD966] text-[#8A6D00] rounded-full text-xs font-semibold"
		if u.IsActive {
			status = "Active"
			statusClass = "px-3 py-1 bg-[#26CE0033] text-[#26CE00] rounded-full text-xs font-semibold"
		}

		rows = append(rows, UserRow{
			ID:            u.ID,
			Name:          u.Name,
			Email:         u.Email,
			Phone:         u.PhoneNumber,
			Status:        status,
			StatusClass:   statusClass,
			TotalBookings: "-", // booking counts not modeled yet
			LastBooking:   "-", // last booking not modeled yet
			CreatedOn:     formatDate(u.CreatedAt),
		})
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "user-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Users": rows,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetViewUserFragment serves the view user fragment populated with the selected user's data.
func (h *AdminHandler) GetViewUserFragment(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Query("id")
	if strings.TrimSpace(userID) == "" {
		return c.Status(400).SendString("User id is required")
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return c.Status(404).SendString("User not found")
	}

	formatDate := func(t time.Time) string {
		if t.IsZero() {
			return "-"
		}
		return t.Format("02 Jan, 2006")
	}

	titleCase := func(s string) string {
		if s == "" {
			return "-"
		}
		s = strings.ToLower(s)
		return strings.ToUpper(string(s[0])) + s[1:]
	}

	statusBadge := func(status string) string {
		switch strings.ToLower(status) {
		case "confirmed", "completed", "paid", "active":
			return "px-3 py-1 bg-[#26CE0033] text-[#26CE00] rounded-full text-xs font-semibold"
		case "pending":
			return "px-3 py-1 bg-[#FFD966] text-[#8A6D00] rounded-full text-xs font-semibold"
		case "cancelled":
			return "px-3 py-1 bg-[#FEE2E2] text-[#B91C1C] rounded-full text-xs font-semibold"
		default:
			return "px-3 py-1 bg-[#E5E7EB] text-[#374151] rounded-full text-xs font-semibold"
		}
	}

	formatAmount := func(currency string, amount float64) string {
		if amount == 0 {
			return "-"
		}
		cur := strings.ToUpper(strings.TrimSpace(currency))
		if cur == "" {
			cur = "USD"
		}
		return fmt.Sprintf("%s %.2f", cur, amount)
	}

	bookings := make([]UserBookingRow, 0)
	lastBookingAt := time.Time{}

	carBookings, _ := h.carBookingRepo.FindByUserID(ctx, userID)
	for _, b := range carBookings {
		bookings = append(bookings, UserBookingRow{
			Provider:    "Car",
			Type:        "Car Rental",
			Status:      titleCase(string(b.Status)),
			StatusClass: statusBadge(string(b.Status)),
			Location:    "-",
			Rating:      "-",
			Revenue:     formatAmount(b.Currency, b.Amount),
			CreatedAt:   b.CreatedAt,
		})
		if b.CreatedAt.After(lastBookingAt) {
			lastBookingAt = b.CreatedAt
		}
	}

	flightBookings, _ := h.flightBookingRepo.FindByUserID(ctx, userID)
	for _, b := range flightBookings {
		bookings = append(bookings, UserBookingRow{
			Provider:    "Flight",
			Type:        "Flight Booking",
			Status:      titleCase(string(b.Status)),
			StatusClass: statusBadge(string(b.Status)),
			Location:    "-",
			Rating:      "-",
			Revenue:     formatAmount(b.Currency, b.Amount),
			CreatedAt:   b.CreatedAt,
		})
		if b.CreatedAt.After(lastBookingAt) {
			lastBookingAt = b.CreatedAt
		}
	}

	hotelBookings, _ := h.hotelBookingRepo.FindByUserID(ctx, userID)
	for _, b := range hotelBookings {
		bookings = append(bookings, UserBookingRow{
			Provider:    "Hotel",
			Type:        "Hotel Booking",
			Status:      titleCase(string(b.Status)),
			StatusClass: statusBadge(string(b.Status)),
			Location:    "-",
			Rating:      "-",
			Revenue:     formatAmount(b.Currency, b.Amount),
			CreatedAt:   b.CreatedAt,
		})
		if b.CreatedAt.After(lastBookingAt) {
			lastBookingAt = b.CreatedAt
		}
	}

	sort.Slice(bookings, func(i, j int) bool {
		return bookings[i].CreatedAt.After(bookings[j].CreatedAt)
	})
	bookingCount := len(bookings)
	if len(bookings) > 15 {
		bookings = bookings[:15]
	}

	lastBooking := "-"
	if !lastBookingAt.IsZero() {
		lastBooking = formatDate(lastBookingAt)
	}

	statusLabel := "Inactive"
	statusClass := "px-3 py-1 bg-[#FFD966] text-[#8A6D00] rounded-full text-xs font-semibold"
	if user.IsActive {
		statusLabel = "Active"
		statusClass = "px-3 py-1 bg-[#26CE0033] text-[#26CE00] rounded-full text-xs font-semibold"
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "view-user-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Name":          user.Name,
		"Email":         user.Email,
		"Phone":         user.PhoneNumber,
		"Status":        statusLabel,
		"StatusClass":   statusClass,
		"TotalBookings": fmt.Sprintf("%d", bookingCount),
		"LastBooking":   lastBooking,
		"CreatedOn":     formatDate(user.CreatedAt),
		"LastLogin":     formatDate(user.LastLogin),
		"Location":      "-",
		"Role":          "Customer",
		"Bookings":      bookings,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}
