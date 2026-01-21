package admin

import (
	"context"
	"html/template"
	"path/filepath"

	adminmodels "Anytime_Travel/backend/internal/models/admin"

	"github.com/gofiber/fiber/v2"
)

// CMSHotels renders the hotels management page
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

// HotelUpdateRequest represents the request body for updating a hotel
type HotelUpdateRequest struct {
	ID              string  `json:"id"`
	HotelID         string  `json:"hotel_id"`
	HotelName       string  `json:"hotel_name"`
	Status          string  `json:"status"`
	Cost            float64 `json:"cost"`
	ProfitPercent   float64 `json:"profit_percent"`
	Location        string  `json:"location"`
	PropertyDetails string  `json:"property_details"`
}

// UpdateHotel handles updating hotel information
func (h *AdminHandler) UpdateHotel(c *fiber.Ctx) error {
	var req HotelUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Validate required fields
	if req.ID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Hotel ID is required",
		})
	}

	if req.HotelName == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Hotel name is required",
		})
	}

	if req.Cost <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Cost must be greater than 0",
		})
	}

	if req.ProfitPercent < 0 || req.ProfitPercent > 100 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Profit percent must be between 0 and 100",
		})
	}

	ctx := context.Background()

	// Find the hotel
	hotel, err := h.hotelRepo.FindByID(ctx, req.ID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Hotel not found",
		})
	}

	// Update hotel fields
	hotel.HotelID = req.HotelID
	hotel.HotelName = req.HotelName
	hotel.Status = adminmodels.HotelStatus(req.Status)
	hotel.Cost = req.Cost
	hotel.ProfitPercent = req.ProfitPercent
	hotel.Location = req.Location
	hotel.PropertyDetails = req.PropertyDetails

	// Update hotel in database
	err = h.hotelRepo.Update(ctx, req.ID, hotel)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update hotel",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Hotel updated successfully",
	})
}

// DeleteHotel handles deleting a hotel
func (h *AdminHandler) DeleteHotel(c *fiber.Ctx) error {
	hotelID := c.Query("id")
	if hotelID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "Hotel ID is required",
		})
	}

	ctx := context.Background()

	// Delete the hotel
	err := h.hotelRepo.Delete(ctx, hotelID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to delete hotel",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Hotel deleted successfully",
	})
}
