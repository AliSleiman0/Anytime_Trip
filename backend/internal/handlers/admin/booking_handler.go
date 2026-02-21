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
	BookingID       string
	RefID           string
	Type            string
	CustomerName    string
	CustomerEmail   string
	CustomerPhone   string
	ServiceProvider string
	BookingDate     time.Time
	TravelDate      time.Time
	Details         string
	Amount          string
	PaymentStatus   string
	PaymentClass    string
	Status          string
	StatusClass     string
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
		"Title":            "Bookings",
		"ShowExportButton": true,
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

	transferBookings, err := h.transferBookingRepo.FindAll(ctx, 100, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching transfer bookings")
	}

	rows := make([]BookingRow, 0, len(carBookings)+len(flightBookings)+len(hotelBookings)+len(transferBookings))

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
		// derive phone and provider safely
		phone := "-"
		if h.userRepo != nil && b.UserID != "" {
			if u, err := h.userRepo.FindByID(ctx, b.UserID); err == nil && u != nil {
				if u.PhoneNumber != "" {
					phone = u.PhoneNumber
				}
			}
		}
		provider := "-"
		if h.carRepo != nil && b.CarID != "" {
			if p, err := h.carRepo.FindByID(ctx, b.CarID); err == nil && p != nil {
				if p.ProviderName != "" {
					provider = p.ProviderName
				}
			}
		}
		row := BookingRow{
			BookingID:       b.BookingID,
			RefID:           b.ID,
			Type:            "Car",
			CustomerName:    b.Customer.Name,
			CustomerEmail:   b.Customer.Email,
			CustomerPhone:   phone,
			ServiceProvider: provider,
			BookingDate:     b.BookingDate,
			TravelDate:      b.BookingDate,
			Details:         b.Details,
			Amount:          formatAmount(b.Currency, b.Amount),
			PaymentStatus:   string(b.PaymentStatus),
			PaymentClass:    badgeClass("payment", string(b.PaymentStatus)),
			Status:          string(b.Status),
			StatusClass:     badgeClass("status", string(b.Status)),
		}

		if matchesFilters(row) {
			rows = append(rows, row)
		}
	}

	for _, b := range flightBookings {
		phone := "-"
		if h.userRepo != nil && b.UserID != "" {
			if u, err := h.userRepo.FindByID(ctx, b.UserID); err == nil && u != nil {
				if u.PhoneNumber != "" {
					phone = u.PhoneNumber
				}
			}
		}
		provider := "-"
		if h.flightRepo != nil && b.FlightID != "" {
			if p, err := h.flightRepo.FindByID(ctx, b.FlightID); err == nil && p != nil {
				if p.ProviderName != "" {
					provider = p.ProviderName
				}
			}
		}
		row := BookingRow{
			BookingID:       b.BookingID,
			RefID:           b.ID,
			Type:            "Flight",
			CustomerName:    b.Customer.Name,
			CustomerEmail:   b.Customer.Email,
			CustomerPhone:   phone,
			ServiceProvider: provider,
			BookingDate:     b.BookingDate,
			TravelDate:      b.BookingDate,
			Details:         b.Details,
			Amount:          formatAmount(b.Pricing.Currency, b.Pricing.Total),
			PaymentStatus:   string(b.PaymentStatus),
			PaymentClass:    badgeClass("payment", string(b.PaymentStatus)),
			Status:          string(b.Status),
			StatusClass:     badgeClass("status", string(b.Status)),
		}

		if matchesFilters(row) {
			rows = append(rows, row)
		}
	}

	for _, b := range hotelBookings {
		phone := "-"
		if h.userRepo != nil && b.UserID != "" {
			if u, err := h.userRepo.FindByID(ctx, b.UserID); err == nil && u != nil {
				if u.PhoneNumber != "" {
					phone = u.PhoneNumber
				}
			}
		}
		provider := "-"
		if h.hotelRepo != nil && b.HotelID != "" {
			if p, err := h.hotelRepo.FindByID(ctx, b.HotelID); err == nil && p != nil {
				if p.ProviderName != "" {
					provider = p.ProviderName
				}
			}
		}
		row := BookingRow{
			BookingID:       b.BookingID,
			RefID:           b.ID,
			Type:            "Hotel",
			CustomerName:    b.Customer.Name,
			CustomerEmail:   b.Customer.Email,
			CustomerPhone:   phone,
			ServiceProvider: provider,
			BookingDate:     b.BookingDate,
			TravelDate:      b.BookingDate,
			Details:         b.Details,
			Amount:          formatAmount(b.Currency, b.Amount),
			PaymentStatus:   string(b.PaymentStatus),
			PaymentClass:    badgeClass("payment", string(b.PaymentStatus)),
			Status:          string(b.Status),
			StatusClass:     badgeClass("status", string(b.Status)),
		}

		if matchesFilters(row) {
			rows = append(rows, row)
		}
	}

	for _, b := range transferBookings {
		phone := "-"
		if h.userRepo != nil && b.UserID != "" {
			if u, err := h.userRepo.FindByID(ctx, b.UserID); err == nil && u != nil {
				if u.PhoneNumber != "" {
					phone = u.PhoneNumber
				}
			}
		}
		provider := "-"
		if h.transferRepo != nil && b.TransferID != "" {
			if p, err := h.transferRepo.FindByID(ctx, b.TransferID); err == nil && p != nil {
				if p.ProviderName != "" {
					provider = p.ProviderName
				}
			}
		}
		row := BookingRow{
			BookingID:       b.BookingID,
			RefID:           b.ID,
			Type:            "Transfer",
			CustomerName:    b.Customer.Name,
			CustomerEmail:   b.Customer.Email,
			CustomerPhone:   phone,
			ServiceProvider: provider,
			BookingDate:     b.BookingDate,
			TravelDate:      b.BookingDate,
			Details:         b.Details,
			Amount:          formatAmount(b.Currency, b.Amount),
			PaymentStatus:   string(b.PaymentStatus),
			PaymentClass:    badgeClass("payment", string(b.PaymentStatus)),
			Status:          string(b.Status),
			StatusClass:     badgeClass("status", string(b.Status)),
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

// ViewBooking renders the booking details page
func (h *AdminHandler) ViewBooking(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "view-booking.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title":     "Booking Details",
		"BookingID": c.Query("id"),
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetViewBookingFragment serves the booking details fragment for HTMX partial loads
func (h *AdminHandler) GetViewBookingFragment(c *fiber.Ctx) error {
	ctx := c.Context()
	bookingID := c.Query("id")

	if bookingID == "" {
		return c.Status(400).SendString("Missing booking ID")
	}

	funcMap := template.FuncMap{
		"title": func(s string) string {
			if s == "" {
				return ""
			}
			return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
		},
		"printf": fmt.Sprintf,
	}

	// Try to find in flight bookings
	flightBooking, err := h.flightBookingRepo.FindByID(ctx, bookingID)
	if err == nil && flightBooking != nil {
		tmpl, err := template.New("view-booking-frag.html").Funcs(funcMap).ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "view-booking-frag.html")))
		if err != nil {
			return c.Status(500).SendString("Error loading template: " + err.Error())
		}

		data := fiber.Map{
			"BookingID":       flightBooking.BookingID,
			"FlightName":      "Flight Details",
			"Type":            "Flight",
			"Details":         flightBooking.Details,
			"Amount":          flightBooking.Pricing.Total,
			"Currency":        flightBooking.Pricing.Currency,
			"PaymentStatus":   string(flightBooking.PaymentStatus),
			"Email":           flightBooking.Customer.Email,
			"Phone":           "-",
			"Seat":            "-",
			"DepartureDate":   flightBooking.BookingDate.Format("01/02/2006"),
			"ArrivalDate":     flightBooking.BookingDate.Format("01/02/2006"),
			"ReferenceID":     flightBooking.ID,
			"RefundAmount":    flightBooking.RefundAmount,
			"RefundRequested": flightBooking.RefundRequested,
		}

		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c, data)
	}

	// Try to find in car bookings
	carBooking, err := h.carBookingRepo.FindByID(ctx, bookingID)
	if err == nil && carBooking != nil {
		tmpl, err := template.New("view-booking-frag.html").Funcs(funcMap).ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "view-booking-frag.html")))
		if err != nil {
			return c.Status(500).SendString("Error loading template: " + err.Error())
		}

		data := fiber.Map{
			"BookingID":       carBooking.BookingID,
			"FlightName":      "Car Booking",
			"Type":            "Car",
			"Details":         carBooking.Details,
			"Amount":          carBooking.Amount,
			"Currency":        carBooking.Currency,
			"PaymentStatus":   string(carBooking.PaymentStatus),
			"Email":           carBooking.Customer.Email,
			"Phone":           "-",
			"Seat":            "-",
			"DepartureDate":   carBooking.BookingDate.Format("01/02/2006"),
			"ArrivalDate":     carBooking.BookingDate.Format("01/02/2006"),
			"ReferenceID":     carBooking.ID,
			"RefundAmount":    carBooking.RefundAmount,
			"RefundRequested": carBooking.RefundRequested,
		}

		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c, data)
	}

	// Try to find in hotel bookings
	hotelBooking, err := h.hotelBookingRepo.FindByID(ctx, bookingID)
	if err == nil && hotelBooking != nil {
		tmpl, err := template.New("view-booking-frag.html").Funcs(funcMap).ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "view-booking-frag.html")))
		if err != nil {
			return c.Status(500).SendString("Error loading template: " + err.Error())
		}

		data := fiber.Map{
			"BookingID":       hotelBooking.BookingID,
			"FlightName":      "Hotel Booking",
			"Type":            "Hotel",
			"Details":         hotelBooking.Details,
			"Amount":          hotelBooking.Amount,
			"Currency":        hotelBooking.Currency,
			"PaymentStatus":   string(hotelBooking.PaymentStatus),
			"Email":           hotelBooking.Customer.Email,
			"Phone":           "-",
			"Seat":            "-",
			"DepartureDate":   hotelBooking.BookingDate.Format("01/02/2006"),
			"ArrivalDate":     hotelBooking.BookingDate.Format("01/02/2006"),
			"ReferenceID":     hotelBooking.ID,
			"RefundAmount":    hotelBooking.RefundAmount,
			"RefundRequested": hotelBooking.RefundRequested,
		}

		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c, data)
	}

	// Try to find in transfer bookings
	transferBooking, err := h.transferBookingRepo.FindByID(ctx, bookingID)
	if err == nil && transferBooking != nil {
		tmpl, err := template.New("view-booking-frag.html").Funcs(funcMap).ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "view-booking-frag.html")))
		if err != nil {
			return c.Status(500).SendString("Error loading template: " + err.Error())
		}

		data := fiber.Map{
			"BookingID":       transferBooking.BookingID,
			"FlightName":      "Transfer Booking",
			"Type":            "Transfer",
			"Details":         transferBooking.Details,
			"Amount":          transferBooking.Amount,
			"Currency":        transferBooking.Currency,
			"PaymentStatus":   string(transferBooking.PaymentStatus),
			"Email":           transferBooking.Customer.Email,
			"Phone":           "-",
			"Seat":            "-",
			"DepartureDate":   transferBooking.BookingDate.Format("01/02/2006"),
			"ArrivalDate":     transferBooking.BookingDate.Format("01/02/2006"),
			"ReferenceID":     transferBooking.ID,
			"RefundAmount":    transferBooking.RefundAmount,
			"RefundRequested": transferBooking.RefundRequested,
		}

		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c, data)
	}

	return c.Status(404).SendString("Booking not found")
}

// UpdateBookingEmail updates the email address for a booking
func (h *AdminHandler) UpdateBookingEmail(c *fiber.Ctx) error {
	ctx := c.Context()

	type UpdateEmailRequest struct {
		BookingID string `json:"booking_id"`
		Email     string `json:"email"`
	}

	var req UpdateEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if req.BookingID == "" || req.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Booking ID and email are required",
		})
	}

	// Try to update in flight bookings
	flightBooking, err := h.flightBookingRepo.FindByID(ctx, req.BookingID)
	if err == nil && flightBooking != nil {
		if err := h.flightBookingRepo.UpdateCustomerEmail(ctx, req.BookingID, req.Email); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to update email",
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Email updated successfully",
		})
	}

	// Try to update in car bookings
	carBooking, err := h.carBookingRepo.FindByID(ctx, req.BookingID)
	if err == nil && carBooking != nil {
		if err := h.carBookingRepo.UpdateCustomerEmail(ctx, req.BookingID, req.Email); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to update email",
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Email updated successfully",
		})
	}

	// Try to update in hotel bookings
	hotelBooking, err := h.hotelBookingRepo.FindByID(ctx, req.BookingID)
	if err == nil && hotelBooking != nil {
		if err := h.hotelBookingRepo.UpdateCustomerEmail(ctx, req.BookingID, req.Email); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to update email",
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Email updated successfully",
		})
	}

	// Try to update in transfer bookings
	transferBooking, err := h.transferBookingRepo.FindByID(ctx, req.BookingID)
	if err == nil && transferBooking != nil {
		if err := h.transferBookingRepo.UpdateCustomerEmail(ctx, req.BookingID, req.Email); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to update email",
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Email updated successfully",
		})
	}

	return c.Status(404).JSON(fiber.Map{
		"success": false,
		"message": "Booking not found",
	})
}

// UpdateBookingRefund updates the refund amount for a booking (flight/car/hotel)
func (h *AdminHandler) UpdateBookingRefund(c *fiber.Ctx) error {
	ctx := c.Context()

	type UpdateRefundRequest struct {
		BookingID    string  `json:"booking_id"`
		RefundAmount float64 `json:"refund_amount"`
	}

	var req UpdateRefundRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if req.BookingID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Booking ID is required",
		})
	}

	// Try to update in flight bookings
	if fb, err := h.flightBookingRepo.FindByID(ctx, req.BookingID); err == nil && fb != nil {
		if err := h.flightBookingRepo.UpdateRefundAmount(ctx, req.BookingID, req.RefundAmount); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update refund"})
		}
		return c.JSON(fiber.Map{"success": true, "message": "Refund updated successfully"})
	}

	// Try car bookings
	if cb, err := h.carBookingRepo.FindByID(ctx, req.BookingID); err == nil && cb != nil {
		if err := h.carBookingRepo.UpdateRefundAmount(ctx, req.BookingID, req.RefundAmount); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update refund"})
		}
		return c.JSON(fiber.Map{"success": true, "message": "Refund updated successfully"})
	}

	// Try hotel bookings
	if hb, err := h.hotelBookingRepo.FindByID(ctx, req.BookingID); err == nil && hb != nil {
		if err := h.hotelBookingRepo.UpdateRefundAmount(ctx, req.BookingID, req.RefundAmount); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update refund"})
		}
		return c.JSON(fiber.Map{"success": true, "message": "Refund updated successfully"})
	}

	return c.Status(404).JSON(fiber.Map{"success": false, "message": "Booking not found"})
}
