package app

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SearchCars handles car search requests
func (h *AppHandler) SearchCars(c *fiber.Ctx) error {
	// Get query parameters
	pickupLocation := c.Query("pickup_location")
	dropoffLocation := c.Query("dropoff_location")
	pickupTimeStr := c.Query("pickup_time")
	dropoffTimeStr := c.Query("dropoff_time")
	carType := c.Query("car_type")
	passengersStr := c.Query("passengers")

	// Validate required parameters
	if pickupLocation == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "pickup_location is required",
		})
	}

	// Parse time values (use zero time if not provided)
	var pickupTime, dropoffTime time.Time
	var err error

	if pickupTimeStr != "" {
		pickupTime, err = time.Parse(time.RFC3339, pickupTimeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid pickup_time format. Use RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)",
			})
		}
	}

	if dropoffTimeStr != "" {
		dropoffTime, err = time.Parse(time.RFC3339, dropoffTimeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid dropoff_time format. Use RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00)",
			})
		}
	}

	// Parse passengers (default to 1 if not provided)
	passengers := 1
	if passengersStr != "" {
		p, err := strconv.Atoi(passengersStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid passengers value",
			})
		}
		passengers = p
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Search for available cars
	cars, err := h.carRepo.SearchCars(ctx, pickupLocation, dropoffLocation, pickupTime, dropoffTime, carType, passengers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search cars",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    cars,
		"count":   len(cars),
	})
}
