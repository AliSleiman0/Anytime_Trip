package admin

import (
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ManageServiceProviders renders the service providers management page
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
