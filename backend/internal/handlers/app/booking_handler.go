package app

import (
	"context"
	"fmt"
	"time"

	appmodels "Anytime_Travel/backend/internal/models/app"

	"github.com/gofiber/fiber/v2"
)

// GetMyBookings returns all bookings for the authenticated user
func (h *AppHandler) GetMyBookings(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIDStr := userID.(string)

	// Debug logging
	fmt.Printf("[DEBUG] Fetching all bookings for user: %s\n", userIDStr)

	// Fetch all booking types
	carBookings, _ := h.carBookingRepo.FindByUserID(ctx, userIDStr)
	flightBookings, _ := h.flightBookingRepo.FindByUserID(ctx, userIDStr)
	hotelBookings, _ := h.hotelBookingRepo.FindByUserID(ctx, userIDStr)
	transferBookings, _ := h.transferBookingRepo.FindByUserID(ctx, userIDStr)

	fmt.Printf("[DEBUG] Found %d car bookings, %d flight bookings, %d hotel bookings, %d transfer bookings\n",
		len(carBookings), len(flightBookings), len(hotelBookings), len(transferBookings))

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"car_bookings":      carBookings,
			"flight_bookings":   flightBookings,
			"hotel_bookings":    hotelBookings,
			"transfer_bookings": transferBookings,
		},
	})
}

// GetMyCarBookings returns car bookings for the authenticated user
func (h *AppHandler) GetMyCarBookings(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bookings, err := h.carBookingRepo.FindByUserID(ctx, userID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch car bookings",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    bookings,
	})
}

// GetMyFlightBookings returns flight bookings for the authenticated user
func (h *AppHandler) GetMyFlightBookings(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bookings, err := h.flightBookingRepo.FindByUserID(ctx, userID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch flight bookings",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    bookings,
	})
}

// GetMyHotelBookings returns hotel bookings for the authenticated user
func (h *AppHandler) GetMyHotelBookings(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bookings, err := h.hotelBookingRepo.FindByUserID(ctx, userID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch hotel bookings",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    bookings,
	})
}

// GetMyTransferBookings returns transfer bookings for the authenticated user
func (h *AppHandler) GetMyTransferBookings(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bookings, err := h.transferBookingRepo.FindByUserID(ctx, userID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch transfer bookings",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    bookings,
	})
}

// GetBookingByID returns a specific booking by ID and type
func (h *AppHandler) GetBookingByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	bookingType := c.Params("type") // car, flight, hotel, transfer
	bookingID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIDStr := userID.(string)

	switch bookingType {
	case "car":
		booking, err := h.carBookingRepo.FindByID(ctx, bookingID)
		if err != nil || booking.UserID != userIDStr {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"data":    booking,
		})

	case "flight":
		booking, err := h.flightBookingRepo.FindByID(ctx, bookingID)
		if err != nil || booking.UserID != userIDStr {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"data":    booking,
		})

	case "hotel":
		booking, err := h.hotelBookingRepo.FindByID(ctx, bookingID)
		if err != nil || booking.UserID != userIDStr {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"data":    booking,
		})

	case "transfer":
		booking, err := h.transferBookingRepo.FindByID(ctx, bookingID)
		if err != nil || booking.UserID != userIDStr {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"data":    booking,
		})

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking type",
		})
	}
}

// CancelBooking cancels a user's booking
func (h *AppHandler) CancelBooking(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	bookingType := c.Params("type")
	bookingID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIDStr := userID.(string)

	fmt.Printf("[CANCEL] Type: %s, BookingID: %s, UserID: %s\n", bookingType, bookingID, userIDStr)

	switch bookingType {
	case "car":
		// Find by booking_id field, not MongoDB ObjectID
		booking, err := h.carBookingRepo.FindByBookingID(ctx, bookingID)
		if err != nil {
			fmt.Printf("[CANCEL-ERROR] Car booking not found: %v\n", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if booking.UserID != userIDStr {
			fmt.Printf("[CANCEL-ERROR] User mismatch: booking.UserID=%s, userID=%s\n", booking.UserID, userIDStr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if err := h.carBookingRepo.UpdateStatus(ctx, booking.ID, appmodels.CarBookingStatusCancelled); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to cancel booking",
			})
		}

		// Send cancellation email with preference checking (use background context for async operation)
		go func() {
			notifyCtx := context.Background()
			h.notificationHelper.SendCarBookingCancellation(
				notifyCtx,
				booking.Customer.Email,
				booking.Customer.Name,
				booking.BookingID,
				booking.CarType,
				booking.Pickup.Location,
				booking.Pickup.Date.Format("Monday, January 2, 2006"),
				booking.Pickup.Time,
				booking.Pricing.Total,
			)
		}()

	case "flight":
		booking, err := h.flightBookingRepo.FindByBookingID(ctx, bookingID)
		if err != nil {
			fmt.Printf("[CANCEL-ERROR] Flight booking not found: %v\n", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if booking.UserID != userIDStr {
			fmt.Printf("[CANCEL-ERROR] User mismatch: booking.UserID=%s, userID=%s\n", booking.UserID, userIDStr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if err := h.flightBookingRepo.UpdateStatus(ctx, booking.ID, appmodels.FlightBookingStatusCancelled); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to cancel booking",
			})
		}

		// Send cancellation email with preference checking (use background context for async operation)
		if len(booking.OutboundFlights) > 0 {
			firstSegment := booking.OutboundFlights[0]
			go func() {
				notifyCtx := context.Background()
				h.notificationHelper.SendFlightBookingCancellation(
					notifyCtx,
					booking.Customer.Email,
					booking.Customer.Name,
					booking.BookingID,
					firstSegment.Airline,
					firstSegment.FlightNumber,
					firstSegment.DepartureCity,
					firstSegment.ArrivalCity,
					firstSegment.DepartureDate,
					firstSegment.DepartureTime,
					booking.Pricing.Total,
				)
			}()
		}

	case "hotel":
		booking, err := h.hotelBookingRepo.FindByBookingID(ctx, bookingID)
		if err != nil {
			fmt.Printf("[CANCEL-ERROR] Hotel booking not found: %v\n", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if booking.UserID != userIDStr {
			fmt.Printf("[CANCEL-ERROR] User mismatch: booking.UserID=%s, userID=%s\n", booking.UserID, userIDStr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if err := h.hotelBookingRepo.UpdateStatus(ctx, booking.ID, appmodels.HotelBookingStatusCancelled); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to cancel booking",
			})
		}

		// Send cancellation email with preference checking (use background context for async operation)
		go func() {
			notifyCtx := context.Background()
			h.notificationHelper.SendHotelBookingCancellation(
				notifyCtx,
				booking.Customer.Email,
				booking.Customer.Name,
				booking.BookingID,
				booking.HotelName,
				booking.CheckInDate.Format("Monday, January 2, 2006"),
				booking.Nights,
				booking.Pricing.Total,
			)
		}()

	case "transfer":
		booking, err := h.transferBookingRepo.FindByBookingID(ctx, bookingID)
		if err != nil {
			fmt.Printf("[CANCEL-ERROR] Transfer booking not found: %v\n", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if booking.UserID != userIDStr {
			fmt.Printf("[CANCEL-ERROR] User mismatch: booking.UserID=%s, userID=%s\n", booking.UserID, userIDStr)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if err := h.transferBookingRepo.UpdateStatus(ctx, booking.ID, appmodels.TransferBookingStatusCancelled); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to cancel booking",
			})
		}

		// Send cancellation email with preference checking (use background context for async operation)
		go func() {
			notifyCtx := context.Background()
			h.notificationHelper.SendTransferBookingCancellation(
				notifyCtx,
				booking.Customer.Email,
				booking.Customer.Name,
				booking.BookingID,
				booking.Vehicle.Type,
				booking.Pickup.Name,
				booking.Pickup.Address,
				booking.Pickup.Date.Format("Monday, January 2, 2006"),
				booking.Pickup.Time,
				booking.Pricing.Total,
			)
		}()

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking type",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Booking cancelled successfully",
	})
}

// RequestRefund marks a booking as having a refund requested
func (h *AppHandler) RequestRefund(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	bookingType := c.Params("type")
	bookingID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIDStr := userID.(string)

	fmt.Printf("[REFUND-REQUEST] Type: %s, BookingID: %s, UserID: %s\n", bookingType, bookingID, userIDStr)

	switch bookingType {
	case "car":
		booking, err := h.carBookingRepo.FindByBookingID(ctx, bookingID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if booking.UserID != userIDStr {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if err := h.carBookingRepo.UpdateRefundRequested(ctx, booking.ID, true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to request refund",
			})
		}

	case "flight":
		booking, err := h.flightBookingRepo.FindByBookingID(ctx, bookingID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if booking.UserID != userIDStr {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if err := h.flightBookingRepo.UpdateRefundRequested(ctx, booking.ID, true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to request refund",
			})
		}

	case "hotel":
		booking, err := h.hotelBookingRepo.FindByBookingID(ctx, bookingID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if booking.UserID != userIDStr {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if err := h.hotelBookingRepo.UpdateRefundRequested(ctx, booking.ID, true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to request refund",
			})
		}

	case "transfer":
		booking, err := h.transferBookingRepo.FindByBookingID(ctx, bookingID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if booking.UserID != userIDStr {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Booking not found",
			})
		}

		if err := h.transferBookingRepo.UpdateRefundRequested(ctx, booking.ID, true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to request refund",
			})
		}

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid booking type",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Refund request submitted successfully",
	})
}
