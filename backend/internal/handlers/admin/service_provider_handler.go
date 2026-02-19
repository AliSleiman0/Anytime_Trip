package admin

import (
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	modelsAdmin "Anytime_Travel/backend/internal/models/admin"
	"Anytime_Travel/backend/internal/repository/app"

	"github.com/gofiber/fiber/v2"
)

// ManageServiceProviders renders the service providers management page
func (h *AdminHandler) ManageServiceProviders(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "service-provider.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title":            "Service Providers",
		"ShowExportButton": true,
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
	transfers, _ := h.transferRepo.FindAll(ctx, 0, 0)

	providers := make([]ServiceProviderView, 0, len(flights)+len(cars)+len(hotels)+len(transfers))

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
			IsFreezed:     f.IsFreezed,
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
			IsFreezed:     car.IsFreezed,
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
			IsFreezed:     htl.IsFreezed,
		})
	}

	for _, transfer := range transfers {
		providers = append(providers, ServiceProviderView{
			ID:            transfer.ID,
			Name:          transfer.ProviderName,
			Type:          transfer.ProviderType,
			Status:        string(transfer.Status),
			Location:      transfer.Location,
			ContactEmail:  transfer.ContactEmail,
			Rating:        transfer.Rating,
			TotalBookings: transfer.TotalBookings,
			Revenue:       transfer.Revenue,
			ProviderType:  "transfer",
			IsFreezed:     transfer.IsFreezed,
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

// ViewService renders the service provider detail page
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
			IsFreezed: flight.IsFreezed,
			FreezeActionText: func() string {
				if flight.IsFreezed {
					return "Unfreeze Account"
				} else {
					return "Freeze Account"
				}
			}(),
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
			Bookings:               getFlightBookingsForProvider(ctx, h.flightBookingRepo, h.userRepo, flight.ID),
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
			IsFreezed: car.IsFreezed,
			FreezeActionText: func() string {
				if car.IsFreezed {
					return "Unfreeze Account"
				} else {
					return "Freeze Account"
				}
			}(),
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
			Bookings:               getCarBookingsForProvider(ctx, h.carBookingRepo, h.userRepo, car.ID),
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
			IsFreezed: hotel.IsFreezed,
			FreezeActionText: func() string {
				if hotel.IsFreezed {
					return "Unfreeze Account"
				} else {
					return "Freeze Account"
				}
			}(),
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
			Bookings:               getHotelBookingsForProvider(ctx, h.hotelBookingRepo, h.userRepo, hotel.ID),
		}
	case "transfer":
		transfer, err := h.transferRepo.FindByID(ctx, id)
		if err != nil {
			return c.Status(404).SendString("transfer provider not found")
		}

		currency := transfer.Currency
		if currency == "" {
			currency = "$"
		}

		detail = ProviderDetailView{
			IsFreezed: transfer.IsFreezed,
			FreezeActionText: func() string {
				if transfer.IsFreezed {
					return "Unfreeze Account"
				} else {
					return "Freeze Account"
				}
			}(),
			ID:                     transfer.ID,
			Name:                   transfer.ProviderName,
			Email:                  transfer.ContactEmail,
			Status:                 string(transfer.Status),
			StatusClass:            statusBadgeClass(string(transfer.Status)),
			Type:                   transfer.ProviderType,
			Location:               transfer.Location,
			CreatedAt:              transfer.CreatedAt.Format("02 Jan, 2006"),
			Rating:                 fmt.Sprintf("%.1f", transfer.Rating),
			TotalBookings:          transfer.TotalBookings,
			Revenue:                formatCurrency(currency, transfer.Revenue),
			ProviderType:           providerType,
			ProfitShareDisplay:     formatProfitPercent(transfer.ProfitPercent, transfer.Profit, transfer.Revenue),
			PhoneNumberDisplay:     "N/A",
			RevenueThisMonth:       "--",
			TotalBookingsThisMonth: "--",
			Bookings:               getTransferBookingsForProvider(ctx, h.userRepo, transfer.ID),
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

// UpdateServiceProviderProfitPercent updates the profit percent for a service provider
func (h *AdminHandler) UpdateServiceProviderProfitPercent(c *fiber.Ctx) error {
	ctx := c.Context()

	type UpdateProfitRequest struct {
		ProviderID    string  `json:"provider_id"`
		ProviderType  string  `json:"provider_type"`
		ProfitPercent float64 `json:"profit_percent"`
	}

	var req UpdateProfitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if req.ProviderID == "" || req.ProviderType == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Provider ID and type are required",
		})
	}

	// Validate profit percent range
	if req.ProfitPercent < 0 || req.ProfitPercent > 1000 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Profit percent must be between 0 and 1000",
		})
	}

	requestType := strings.ToLower(req.ProviderType)

	switch requestType {
	case "flight":
		flight, err := h.flightRepo.FindByID(ctx, req.ProviderID)
		if err != nil || flight == nil {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Flight provider not found",
			})
		}
		if err := h.flightRepo.UpdateProfitPercent(ctx, req.ProviderID, req.ProfitPercent); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to update profit percent",
			})
		}
	case "car":
		car, err := h.carRepo.FindByID(ctx, req.ProviderID)
		if err != nil || car == nil {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Car provider not found",
			})
		}
		if err := h.carRepo.UpdateProfitPercent(ctx, req.ProviderID, req.ProfitPercent); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to update profit percent",
			})
		}
	case "hotel":
		hotel, err := h.hotelRepo.FindByID(ctx, req.ProviderID)
		if err != nil || hotel == nil {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Hotel provider not found",
			})
		}
		if err := h.hotelRepo.UpdateProfitPercent(ctx, req.ProviderID, req.ProfitPercent); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"success": false,
				"message": "Failed to update profit percent",
			})
		}
	default:
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid provider type",
		})
	}

	return c.JSON(fiber.Map{
		"success":        true,
		"message":        "Profit percent updated successfully",
		"profit_percent": req.ProfitPercent,
	})
}

// ToggleFreezeProvider toggles the is_freezed flag for a service provider
func (h *AdminHandler) ToggleFreezeProvider(c *fiber.Ctx) error {
	ctx := c.Context()

	type Req struct {
		ProviderID   string `json:"provider_id"`
		ProviderType string `json:"provider_type"`
	}

	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}

	if req.ProviderID == "" || req.ProviderType == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Provider ID and type are required"})
	}

	switch strings.ToLower(req.ProviderType) {
	case "flight":
		flight, err := h.flightRepo.FindByID(ctx, req.ProviderID)
		if err != nil || flight == nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Flight provider not found"})
		}
		newState := !flight.IsFreezed
		if err := h.flightRepo.UpdateFreeze(ctx, req.ProviderID, newState); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update provider state"})
		}
		return c.JSON(fiber.Map{"success": true, "is_freezed": newState})
	case "car":
		car, err := h.carRepo.FindByID(ctx, req.ProviderID)
		if err != nil || car == nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Car provider not found"})
		}
		newState := !car.IsFreezed
		if err := h.carRepo.UpdateFreeze(ctx, req.ProviderID, newState); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update provider state"})
		}
		return c.JSON(fiber.Map{"success": true, "is_freezed": newState})
	case "hotel":
		hotel, err := h.hotelRepo.FindByID(ctx, req.ProviderID)
		if err != nil || hotel == nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Hotel provider not found"})
		}
		newState := !hotel.IsFreezed
		if err := h.hotelRepo.UpdateFreeze(ctx, req.ProviderID, newState); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update provider state"})
		}
		return c.JSON(fiber.Map{"success": true, "is_freezed": newState})
	default:
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid provider type"})
	}
}

// ToggleArchiveProvider toggles archived state for a service provider by updating its status to "archived" or "active"
func (h *AdminHandler) ToggleArchiveProvider(c *fiber.Ctx) error {
	ctx := c.Context()

	type Req struct {
		ProviderID   string `json:"provider_id"`
		ProviderType string `json:"provider_type"`
	}

	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}

	if req.ProviderID == "" || req.ProviderType == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Provider ID and type are required"})
	}

	switch strings.ToLower(req.ProviderType) {
	case "flight":
		flight, err := h.flightRepo.FindByID(ctx, req.ProviderID)
		if err != nil || flight == nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Flight provider not found"})
		}
		// toggle between archived and active
		var newStatus string
		if strings.EqualFold(string(flight.Status), "archived") {
			newStatus = "active"
		} else {
			newStatus = "archived"
		}
		if err := h.flightRepo.UpdateStatus(ctx, req.ProviderID, modelsAdmin.FlightStatus(newStatus)); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update provider status"})
		}
		return c.JSON(fiber.Map{"success": true, "status": newStatus})
	case "car":
		car, err := h.carRepo.FindByID(ctx, req.ProviderID)
		if err != nil || car == nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Car provider not found"})
		}
		var newStatus string
		if strings.EqualFold(string(car.Status), "archived") {
			newStatus = "active"
		} else {
			newStatus = "archived"
		}
		if err := h.carRepo.UpdateStatus(ctx, req.ProviderID, modelsAdmin.CarStatus(newStatus)); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update provider status"})
		}
		return c.JSON(fiber.Map{"success": true, "status": newStatus})
	case "hotel":
		hotel, err := h.hotelRepo.FindByID(ctx, req.ProviderID)
		if err != nil || hotel == nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Hotel provider not found"})
		}
		var newStatus string
		if strings.EqualFold(string(hotel.Status), "archived") {
			newStatus = "active"
		} else {
			newStatus = "archived"
		}
		if err := h.hotelRepo.UpdateStatus(ctx, req.ProviderID, modelsAdmin.HotelStatus(newStatus)); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update provider status"})
		}
		return c.JSON(fiber.Map{"success": true, "status": newStatus})
	default:
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid provider type"})
	}
}

// ToggleActivateProvider toggles between active and inactive status for a service provider
func (h *AdminHandler) ToggleActivateProvider(c *fiber.Ctx) error {
	ctx := c.Context()

	type Req struct {
		ProviderID   string `json:"provider_id"`
		ProviderType string `json:"provider_type"`
	}

	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}

	if req.ProviderID == "" || req.ProviderType == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Provider ID and type are required"})
	}

	switch strings.ToLower(req.ProviderType) {
	case "flight":
		flight, err := h.flightRepo.FindByID(ctx, req.ProviderID)
		if err != nil || flight == nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Flight provider not found"})
		}
		var newStatus string
		if strings.EqualFold(string(flight.Status), "inactive") {
			newStatus = "active"
		} else {
			newStatus = "inactive"
		}
		if err := h.flightRepo.UpdateStatus(ctx, req.ProviderID, modelsAdmin.FlightStatus(newStatus)); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update provider status"})
		}
		return c.JSON(fiber.Map{"success": true, "status": newStatus})
	case "car":
		car, err := h.carRepo.FindByID(ctx, req.ProviderID)
		if err != nil || car == nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Car provider not found"})
		}
		var newStatus string
		if strings.EqualFold(string(car.Status), "inactive") {
			newStatus = "active"
		} else {
			newStatus = "inactive"
		}
		if err := h.carRepo.UpdateStatus(ctx, req.ProviderID, modelsAdmin.CarStatus(newStatus)); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update provider status"})
		}
		return c.JSON(fiber.Map{"success": true, "status": newStatus})
	case "hotel":
		hotel, err := h.hotelRepo.FindByID(ctx, req.ProviderID)
		if err != nil || hotel == nil {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Hotel provider not found"})
		}
		var newStatus string
		if strings.EqualFold(string(hotel.Status), "inactive") {
			newStatus = "active"
		} else {
			newStatus = "inactive"
		}
		if err := h.hotelRepo.UpdateStatus(ctx, req.ProviderID, modelsAdmin.HotelStatus(newStatus)); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Failed to update provider status"})
		}
		return c.JSON(fiber.Map{"success": true, "status": newStatus})
	default:
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid provider type"})
	}
}

// getFlightBookingsForProvider fetches all flight bookings for a specific flight provider
func getFlightBookingsForProvider(ctx context.Context, flightBookingRepo *app.FlightBookingRepository, userRepo *app.UserRepository, flightProviderID string) []BookingView {
	bookings := make([]BookingView, 0)

	allBookings, err := flightBookingRepo.FindAll(ctx, 1000, 0)
	if err != nil {
		return bookings
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

	for _, b := range allBookings {
		if b.FlightID != flightProviderID {
			continue
		}

		user, _ := userRepo.FindByID(ctx, b.UserID)
		userName := "Unknown"
		if user != nil {
			userName = user.Name
		}

		bookings = append(bookings, BookingView{
			UserName:    userName,
			BookingType: "Flight Booking",
			Status:      titleCase(string(b.Status)),
			StatusClass: statusBadge(string(b.Status)),
			Destination: "-",
			BookingDate: b.CreatedAt.Format("02 Jan, 2006"),
			Amount:      fmt.Sprintf("%s %.2f", strings.ToUpper(b.Pricing.Currency), b.Pricing.Total),
		})
	}

	return bookings
}

// getCarBookingsForProvider fetches all car bookings for a specific car provider
func getCarBookingsForProvider(ctx context.Context, carBookingRepo *app.CarBookingRepository, userRepo *app.UserRepository, carProviderID string) []BookingView {
	bookings := make([]BookingView, 0)

	allBookings, err := carBookingRepo.FindAll(ctx, 1000, 0)
	if err != nil {
		return bookings
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

	for _, b := range allBookings {
		if b.CarID != carProviderID {
			continue
		}

		user, _ := userRepo.FindByID(ctx, b.UserID)
		userName := "Unknown"
		if user != nil {
			userName = user.Name
		}

		bookings = append(bookings, BookingView{
			UserName:    userName,
			BookingType: "Car Rental",
			Status:      titleCase(string(b.Status)),
			StatusClass: statusBadge(string(b.Status)),
			Destination: "-",
			BookingDate: b.CreatedAt.Format("02 Jan, 2006"),
			Amount:      fmt.Sprintf("%s %.2f", strings.ToUpper(b.Currency), b.Amount),
		})
	}

	return bookings
}

// getHotelBookingsForProvider fetches all hotel bookings for a specific hotel provider
func getHotelBookingsForProvider(ctx context.Context, hotelBookingRepo *app.HotelBookingRepository, userRepo *app.UserRepository, hotelProviderID string) []BookingView {
	bookings := make([]BookingView, 0)

	allBookings, err := hotelBookingRepo.FindAll(ctx, 1000, 0)
	if err != nil {
		return bookings
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

	for _, b := range allBookings {
		if b.HotelID != hotelProviderID {
			continue
		}

		user, _ := userRepo.FindByID(ctx, b.UserID)
		userName := "Unknown"
		if user != nil {
			userName = user.Name
		}

		bookings = append(bookings, BookingView{
			UserName:    userName,
			BookingType: "Hotel Booking",
			Status:      titleCase(string(b.Status)),
			StatusClass: statusBadge(string(b.Status)),
			Destination: "-",
			BookingDate: b.CreatedAt.Format("02 Jan, 2006"),
			Amount:      fmt.Sprintf("%s %.2f", strings.ToUpper(b.Currency), b.Amount),
		})
	}

	return bookings
}
func getTransferBookingsForProvider(ctx context.Context, userRepo *app.UserRepository, transferProviderID string) []BookingView {
	bookings := make([]BookingView, 0)

	// TODO: Implement transfer booking repository and fetch bookings
	// For now, return empty bookings list until transfer booking repository is created
	return bookings
}
