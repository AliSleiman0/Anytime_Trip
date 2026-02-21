package app

import (
	"context"
	"strconv"
	"time"

	"Anytime_Travel/backend/internal/repository/admin"

	"github.com/gofiber/fiber/v2"
)

// SearchCars handles car search requests
func (h *AppHandler) SearchCars(c *fiber.Ctx) error {
	// Get query parameters
	pickupLocation := c.Query("pickup_location")
	dropoffLocation := c.Query("dropoff_location")
	pickupDateStr := c.Query("pickup_date")
	dropoffDateStr := c.Query("dropoff_date")
	carType := c.Query("car_type")
	passengersStr := c.Query("passengers")

	// Validate required parameters
	if pickupLocation == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "pickup_location is required",
		})
	}

	// Parse dates - handle ISO8601 and RFC3339 formats
	var pickupTime, dropoffTime time.Time
	var err error

	// Helper function to parse datetime strings in multiple formats
	parseDateTime := func(dateStr string) (time.Time, error) {
		if dateStr == "" {
			return time.Time{}, nil
		}

		// Try RFC3339 format first (handles ISO8601)
		t, err := time.Parse(time.RFC3339, dateStr)
		if err == nil {
			return t, nil
		}

		// Try RFC3339Nano format (handles milliseconds)
		t, err = time.Parse(time.RFC3339Nano, dateStr)
		if err == nil {
			return t, nil
		}

		// Try basic ISO8601 date format
		t, err = time.Parse("2006-01-02", dateStr)
		if err == nil {
			return t, nil
		}

		// Try datetime format without timezone
		t, err = time.Parse("2006-01-02T15:04:05", dateStr)
		if err == nil {
			return t, nil
		}

		return time.Time{}, err
	}

	if pickupDateStr != "" {
		pickupTime, err = parseDateTime(pickupDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid pickup_date format. Expected ISO8601 (e.g., 2006-01-02T15:04:05Z)",
			})
		}
	}

	if dropoffDateStr != "" {
		dropoffTime, err = parseDateTime(dropoffDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid dropoff_date format. Expected ISO8601 (e.g., 2006-01-02T15:04:05Z)",
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

	// Filter cars by availability if dates are provided
	var availableCars []interface{}
	if !pickupTime.IsZero() && !dropoffTime.IsZero() {
		for _, car := range cars {
			// Check if car is booked during the requested period
			booked, err := h.carBookingRepo.CheckCarBooked(ctx, car.ID, pickupTime, dropoffTime)
			if err != nil {
				continue // Skip this car if we can't check availability
			}

			if !booked {
				// Calculate rental cost
				rentalCost := admin.CalculateRentalCost(car.Cost, pickupTime, dropoffTime)

				// Add calculated cost to car data
				carData := map[string]interface{}{
					"id":                 car.ID,
					"car_id":             car.CarID,
					"car_name":           car.CarName,
					"car_type":           car.CarType,
					"cost":               car.Cost,
					"cost_per_day":       car.Cost,
					"calculated_cost":    rentalCost,
					"currency":           car.Currency,
					"passengers":         car.Passengers,
					"mileage":            car.Mileage,
					"pickup_location":    car.PickupLocation,
					"drop_off_location":  car.DropOffLocation,
					"unlimited_mileage":  car.UnlimitedMileage,
					"shuttle_to_counter": car.ShuttleToCounter,
					"transmission":       car.Transmission,
					"extras":             car.Extras,
					"provider_name":      car.ProviderName,
					"location":           car.Location,
					"rating":             car.Rating,
					"total_bookings":     car.TotalBookings,
				}
				availableCars = append(availableCars, carData)
			}
		}
	} else {
		// If no dates provided, return cars without availability filtering
		for _, car := range cars {
			carData := map[string]interface{}{
				"id":                 car.ID,
				"car_id":             car.CarID,
				"car_name":           car.CarName,
				"car_type":           car.CarType,
				"cost":               car.Cost,
				"cost_per_day":       car.Cost,
				"currency":           car.Currency,
				"passengers":         car.Passengers,
				"mileage":            car.Mileage,
				"pickup_location":    car.PickupLocation,
				"drop_off_location":  car.DropOffLocation,
				"unlimited_mileage":  car.UnlimitedMileage,
				"shuttle_to_counter": car.ShuttleToCounter,
				"transmission":       car.Transmission,
				"extras":             car.Extras,
				"provider_name":      car.ProviderName,
				"location":           car.Location,
				"rating":             car.Rating,
				"total_bookings":     car.TotalBookings,
			}
			availableCars = append(availableCars, carData)
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    availableCars,
		"count":   len(availableCars),
	})
}

// GetAvailableCarLocations returns all available pickup and dropoff locations
func (h *AppHandler) GetAvailableCarLocations(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	locations, err := h.carRepo.GetAvailableLocations(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch available locations",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    locations,
	})
}
