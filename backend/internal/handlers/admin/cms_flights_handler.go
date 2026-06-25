package admin

import (
	"context"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"travel/backend/core/utils"
	adminmodels "travel/backend/internal/models/admin"

	"github.com/gofiber/fiber/v2"
)

// CMSFlights renders the flights management page
func (h *AdminHandler) CMSFlights(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "cms-flights.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Flights Management",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetCMSFlightsFragment serves the flights management fragment for HTMX partial loads
func (h *AdminHandler) GetCMSFlightsFragment(c *fiber.Ctx) error {
	ctx := context.Background()

	// Fetch all flights from database
	flights, err := h.flightRepo.FindAll(ctx, 100, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching flights")
	}

	// Create template with custom functions
	funcMap := template.FuncMap{
		"formatDateTime": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02T15:04")
		},
		"formatTime": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("3:04 PM")
		},
	}

	tmpl, err := template.New("cms-flights-frag.html").Funcs(funcMap).ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "cms-flights-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Flights": flights,
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// UpdateFlight handles flight update requests from the CMS modal
func (h *AdminHandler) UpdateFlight(c *fiber.Ctx) error {
	ctx := context.Background()

	type FlightUpdateRequest struct {
		ID                string  `json:"id"`
		FlightNumber      string  `json:"flight_number"`
		FlightID          string  `json:"flight_id"`
		Status            string  `json:"status"`
		Cost              float64 `json:"cost"`
		ProfitPercent     float64 `json:"profit_percent"`
		DepartureTime     string  `json:"departure_time"`
		ArrivalTime       string  `json:"arrival_time"`
		DepartureLocation string  `json:"departure_location"`
		ArrivalLocation   string  `json:"arrival_location"`
		Refundable        bool    `json:"refundable"`
		AllowChanges      bool    `json:"allow_changes"`
		AllowSeatChoice   bool    `json:"allow_seat_choice"`
		AllowCarryOn      bool    `json:"allow_carry_on"`
	}

	var req FlightUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate inputs
	var errs []string
	if msg := utils.ValidateRequired("ID", req.ID); msg != "" {
		errs = append(errs, msg)
	}
	if msg := utils.ValidateLengthBetween("Flight number", req.FlightNumber, 2, 20); msg != "" {
		errs = append(errs, msg)
	}
	if msg := utils.ValidateRequired("Flight ID", req.FlightID); msg != "" {
		errs = append(errs, msg)
	}
	if msg := utils.ValidateEnum("Status", req.Status, []string{string(adminmodels.FlightStatusActive), string(adminmodels.FlightStatusInactive), string(adminmodels.FlightStatusCancelled), string(adminmodels.FlightStatusDelayed)}); msg != "" {
		errs = append(errs, msg)
	}
	if msg := utils.ValidatePositiveFloat("Cost", req.Cost); msg != "" {
		errs = append(errs, msg)
	}
	if msg := utils.ValidateFloatRange("Profit percent", req.ProfitPercent, 0, 100); msg != "" {
		errs = append(errs, msg)
	}
	if msg := utils.ValidateRequired("Departure location", req.DepartureLocation); msg != "" {
		errs = append(errs, msg)
	}
	if msg := utils.ValidateRequired("Arrival location", req.ArrivalLocation); msg != "" {
		errs = append(errs, msg)
	}

	var parsedDeparture time.Time
	var parsedArrival time.Time

	if strings.TrimSpace(req.DepartureTime) != "" {
		dt, err := time.Parse(time.RFC3339, req.DepartureTime)
		if err != nil {
			errs = append(errs, "Departure time must be an RFC3339 timestamp")
		} else {
			parsedDeparture = dt
		}
	}

	if strings.TrimSpace(req.ArrivalTime) != "" {
		at, err := time.Parse(time.RFC3339, req.ArrivalTime)
		if err != nil {
			errs = append(errs, "Arrival time must be an RFC3339 timestamp")
		} else {
			parsedArrival = at
		}
	}

	if !parsedDeparture.IsZero() && !parsedArrival.IsZero() && parsedArrival.Before(parsedDeparture) {
		errs = append(errs, "Arrival time must be after departure time")
	}

	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	// Fetch existing flight
	flight, err := h.flightRepo.FindByID(ctx, req.ID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Flight not found",
		})
	}

	// Update flight fields
	flight.FlightNumber = req.FlightNumber
	flight.FlightID = req.FlightID
	flight.Status = adminmodels.FlightStatus(req.Status)
	flight.Cost = req.Cost
	flight.ProfitPercent = req.ProfitPercent
	flight.DepartureLocation = req.DepartureLocation
	flight.ArrivalLocation = req.ArrivalLocation
	flight.Refundable = req.Refundable
	flight.AllowChanges = req.AllowChanges
	flight.AllowSeatChoice = req.AllowSeatChoice
	flight.AllowCarryOn = req.AllowCarryOn

	// Parse and update departure/arrival times if provided
	if !parsedDeparture.IsZero() {
		flight.DepartureTime = parsedDeparture
	}
	if !parsedArrival.IsZero() {
		flight.ArrivalTime = parsedArrival
	}

	// Update in database
	if err := h.flightRepo.Update(ctx, req.ID, flight); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update flight",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Flight updated successfully",
	})
}

// DeleteFlight handles flight deletion from CMS
func (h *AdminHandler) DeleteFlight(c *fiber.Ctx) error {
	ctx := context.Background()
	id := c.Query("id")

	if msg := utils.ValidateRequired("Flight ID", id); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": []string{msg}})
	}

	if err := h.flightRepo.Delete(ctx, id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete flight",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Flight deleted successfully",
	})
}
