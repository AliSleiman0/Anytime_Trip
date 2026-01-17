package admin

import (
	"context"
	"html/template"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CMSCars renders the cars management page
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

	// Create template with custom functions
	funcMap := template.FuncMap{
		"formatDateTime": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02T15:04")
		},
		"formatDate": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("Mon, Jan 2, 3:04 PM")
		},
	}

	tmpl, err := template.New("cms-cars-frag.html").Funcs(funcMap).ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-cars-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Cars": cars,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// UpdateCar handles car update requests from the CMS modal
func (h *AdminHandler) UpdateCar(c *fiber.Ctx) error {
	ctx := context.Background()

	type CarUpdateRequest struct {
		ID               string   `json:"id"`
		CarName          string   `json:"car_name"`
		CarType          string   `json:"car_type"`
		Passengers       int      `json:"passengers"`
		Mileage          string   `json:"mileage"`
		PickupLocation   string   `json:"pickup_location"`
		PickupTime       string   `json:"pickup_time"`
		DropOffLocation  string   `json:"drop_off_location"`
		DropOffTime      string   `json:"drop_off_time"`
		ProfitPercent    float64  `json:"profit_percent"`
		UnlimitedMileage bool     `json:"unlimited_mileage"`
		ShuttleToCounter bool     `json:"shuttle_to_counter"`
		Extras           []string `json:"extras"`
	}

	var req CarUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Fetch existing car
	car, err := h.carRepo.FindByID(ctx, req.ID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Car not found",
		})
	}

	// Update car fields
	car.CarName = req.CarName
	car.CarType = req.CarType
	car.Passengers = req.Passengers
	car.Mileage = req.Mileage
	car.PickupLocation = req.PickupLocation
	car.DropOffLocation = req.DropOffLocation
	car.ProfitPercent = req.ProfitPercent
	car.UnlimitedMileage = req.UnlimitedMileage
	car.ShuttleToCounter = req.ShuttleToCounter
	car.Extras = req.Extras

	// Parse and update pickup/dropoff times if provided
	if req.PickupTime != "" {
		if pickupTime, err := time.Parse(time.RFC3339, req.PickupTime); err == nil {
			car.PickupTime = pickupTime
		}
	}
	if req.DropOffTime != "" {
		if dropOffTime, err := time.Parse(time.RFC3339, req.DropOffTime); err == nil {
			car.DropOffTime = dropOffTime
		}
	}

	// Update in database
	if err := h.carRepo.Update(ctx, req.ID, car); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update car",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Car updated successfully",
	})
}

// DeleteCar handles car deletion from CMS
func (h *AdminHandler) DeleteCar(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Query("id")

	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Car ID is required",
		})
	}

	if err := h.carRepo.Delete(ctx, id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete car",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Car deleted successfully",
	})
}
