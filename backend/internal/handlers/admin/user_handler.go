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
	Initials      string
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

	// Get filter parameters
	createdDate := c.Query("created_date")
	lastBooking := c.Query("last_booking")
	statusFilter := c.Query("status")

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

		// Apply status filter
		if statusFilter != "" && statusFilter != "all" {
			if statusFilter == "active" && !u.IsActive {
				continue
			}
			if statusFilter == "inactive" && u.IsActive {
				continue
			}
		}

		// Apply created date filter
		if createdDate != "" {
			filterDate, err := time.Parse("2006-01-02", createdDate)
			if err == nil {
				// Filter users created on the selected date
				if !isSameDay(u.CreatedAt, filterDate) {
					continue
				}
			}
		}

		// Apply last booking filter
		if lastBooking != "" && lastBooking != "all" {
			now := time.Now()
			switch lastBooking {
			case "today":
				if !isSameDay(u.LastBooking, now) {
					continue
				}
			case "this_week":
				weekAgo := now.AddDate(0, 0, -7)
				if u.LastBooking.IsZero() || u.LastBooking.Before(weekAgo) {
					continue
				}
			case "this_month":
				monthAgo := now.AddDate(0, -1, 0)
				if u.LastBooking.IsZero() || u.LastBooking.Before(monthAgo) {
					continue
				}
			case "older":
				monthAgo := now.AddDate(0, -1, 0)
				if u.LastBooking.IsZero() || u.LastBooking.After(monthAgo) {
					continue
				}
			}
		}

		// Extract initials from name
		initials := "U"
		if u.Name != "" {
			parts := strings.Fields(u.Name)
			if len(parts) == 1 {
				initials = strings.ToUpper(string(parts[0][0]))
			} else if len(parts) >= 2 {
				initials = strings.ToUpper(string(parts[0][0])) + strings.ToUpper(string(parts[len(parts)-1][0]))
			}
		}

		// Format total bookings
		totalBookingsStr := "0"
		if u.TotalBookings > 0 {
			totalBookingsStr = fmt.Sprintf("%d", u.TotalBookings)
		}

		// Format last booking
		lastBookingStr := "-"
		if !u.LastBooking.IsZero() {
			lastBookingStr = formatDate(u.LastBooking)
		}

		rows = append(rows, UserRow{
			ID:            u.ID,
			Name:          u.Name,
			Initials:      initials,
			Email:         u.Email,
			Phone:         u.PhoneNumber,
			Status:        status,
			StatusClass:   statusClass,
			TotalBookings: totalBookingsStr,
			LastBooking:   lastBookingStr,
			CreatedOn:     formatDate(u.CreatedAt),
		})
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "user-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Users": rows,
		"Filters": fiber.Map{
			"CreatedDate": createdDate,
			"LastBooking": lastBooking,
			"Status":      statusFilter,
		},
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// isSameDay checks if two times are on the same day
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
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
