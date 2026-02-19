package utils

import (
	"context"
	"fmt"
	"time"

	appmodels "Anytime_Travel/backend/internal/models/app"
	adminrepo "Anytime_Travel/backend/internal/repository/admin"
	apprepo "Anytime_Travel/backend/internal/repository/app"
)

// NotificationHelper handles notification logic with preference checking
type NotificationHelper struct {
	emailService               *EmailService
	adminNotificationPrefsRepo *adminrepo.NotificationPreferencesRepository
	userRepo                   *apprepo.UserRepository
}

// NewNotificationHelper creates a new notification helper
func NewNotificationHelper(
	emailService *EmailService,
	adminNotificationPrefsRepo *adminrepo.NotificationPreferencesRepository,
	userRepo *apprepo.UserRepository,
) *NotificationHelper {
	return &NotificationHelper{
		emailService:               emailService,
		adminNotificationPrefsRepo: adminNotificationPrefsRepo,
		userRepo:                   userRepo,
	}
}

// SendNewBookingEmail sends a new booking notification email after checking preferences
// Priority 1: Check admin notification preferences
// Priority 2: Check user notification preferences
func (nh *NotificationHelper) SendNewBookingEmail(
	ctx context.Context,
	userID string,
	bookingID string,
	bookingType string,
	amount float64,
	bookingDate time.Time,
) error {
	fmt.Printf("[NOTIFICATION] ====== Starting SendNewBookingEmail ======\n")
	fmt.Printf("[NOTIFICATION] UserID: %s, BookingID: %s, Type: %s, Amount: %.2f\n", userID, bookingID, bookingType, amount)

	// Get user information
	user, err := nh.userRepo.FindByID(ctx, userID)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: User %s not found, skipping email notification: %v\n", userID, err)
		return nil // Don't fail the notification if user doesn't exist
	}
	fmt.Printf("[NOTIFICATION] Found user: %s (%s)\n", user.Email, user.Name)

	// Priority 1: Check admin notification preferences
	// Get all admin preferences and check if ANY admin has new booking notifications enabled
	adminPrefs, err := nh.adminNotificationPrefsRepo.FindAll(ctx)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: Could not fetch admin preferences: %v\n", err)
		// Continue anyway - don't block on admin prefs error
	}

	adminAllowsEmail := false
	if len(adminPrefs) > 0 {
		fmt.Printf("[NOTIFICATION] Checking %d admin preference(s)\n", len(adminPrefs))
		// Check if any admin has enabled new booking notifications
		for _, pref := range adminPrefs {
			fmt.Printf("[NOTIFICATION] Admin pref - NewBooking enabled: %v\n", pref.Bookings.NewBooking)
			if pref.Bookings.NewBooking {
				adminAllowsEmail = true
				break
			}
		}
	} else {
		// If no admin preferences found, default to allowing emails
		// (This handles the case where admin preferences haven't been set up yet)
		adminAllowsEmail = true
		fmt.Printf("[NOTIFICATION] No admin preferences found, defaulting to allow email\n")
	}
	fmt.Printf("[NOTIFICATION] Admin allows email: %v\n", adminAllowsEmail)

	if !adminAllowsEmail {
		fmt.Printf("[NOTIFICATION] Email blocked by admin preferences for booking %s\n", bookingID)
		return nil // Don't send email if admin has disabled new booking notifications
	}

	// Priority 2: Check user notification preferences
	fmt.Printf("[NOTIFICATION] User email preference enabled: %v\n", user.NotificationPreferences.Email)
	if !user.NotificationPreferences.Email {
		fmt.Printf("[NOTIFICATION] Email blocked by user preferences for user %s\n", userID)
		return nil // Don't send email if user has disabled email notifications
	}

	// Both admin and user allow email - send it
	fmt.Printf("[NOTIFICATION] Sending new booking email to %s for booking %s\n", user.Email, bookingID)
	err = nh.emailService.SendNewBookingNotification(
		user.Email,
		user.Name,
		bookingID,
		bookingType,
		amount,
		bookingDate,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("[NOTIFICATION] Successfully sent new booking email for %s\n", bookingID)
	return nil
}

// Helper functions for specific booking types

// SendHotelBookingEmail sends notification for hotel booking
func (nh *NotificationHelper) SendHotelBookingEmail(
	ctx context.Context,
	booking *appmodels.HotelBooking,
) error {
	// Convert guests to string array
	guestNames := make([]string, len(booking.Guests))
	for i, g := range booking.Guests {
		if g.Type != "" && g.Type != "Adult" {
			guestNames[i] = fmt.Sprintf("%s (%s)", g.Name, g.Type)
		} else {
			guestNames[i] = g.Name
		}
	}

	// Use the detailed hotel booking email template
	return nh.emailService.SendHotelBookingEmail(
		booking.Customer.Email,
		booking.Customer.Name,
		booking.BookingID,
		booking.BookingDate,
		guestNames,
		booking.HotelName,
		booking.HotelImage,
		booking.CheckInDate,
		booking.CheckOutDate,
		booking.Nights,
		booking.Rooms,
		booking.Pricing.NightPrice,
		booking.Pricing.Taxes,
		booking.Pricing.DestinationFee,
		booking.Pricing.ServiceFee,
		booking.Pricing.Total,
	)
}

// SendFlightBookingEmail sends notification for flight booking
func (nh *NotificationHelper) SendFlightBookingEmail(
	ctx context.Context,
	booking *appmodels.FlightBooking,
) error {
	// Convert travelers to string array
	travelers := make([]string, len(booking.Travelers))
	for i, t := range booking.Travelers {
		if t.Type != "" && t.Type != "Adult" {
			travelers[i] = fmt.Sprintf("%s (%s)", t.Name, t.Type)
		} else {
			travelers[i] = t.Name
		}
	}

	// Convert flight segments to email format
	outboundSegments := make([]FlightSegment, len(booking.OutboundFlights))
	for i, seg := range booking.OutboundFlights {
		outboundSegments[i] = FlightSegment{
			Airline:          seg.Airline,
			FlightNumber:     seg.FlightNumber,
			Class:            seg.Class,
			DepartureCity:    seg.DepartureCity,
			DepartureCode:    seg.DepartureCode,
			DepartureAirport: seg.DepartureAirport,
			DepartureDate:    seg.DepartureDate,
			DepartureTime:    seg.DepartureTime,
			ArrivalCity:      seg.ArrivalCity,
			ArrivalCode:      seg.ArrivalCode,
			ArrivalAirport:   seg.ArrivalAirport,
			ArrivalDate:      seg.ArrivalDate,
			ArrivalTime:      seg.ArrivalTime,
			Duration:         seg.Duration,
		}
	}

	returnSegments := make([]FlightSegment, len(booking.ReturnFlights))
	for i, seg := range booking.ReturnFlights {
		returnSegments[i] = FlightSegment{
			Airline:          seg.Airline,
			FlightNumber:     seg.FlightNumber,
			Class:            seg.Class,
			DepartureCity:    seg.DepartureCity,
			DepartureCode:    seg.DepartureCode,
			DepartureAirport: seg.DepartureAirport,
			DepartureDate:    seg.DepartureDate,
			DepartureTime:    seg.DepartureTime,
			ArrivalCity:      seg.ArrivalCity,
			ArrivalCode:      seg.ArrivalCode,
			ArrivalAirport:   seg.ArrivalAirport,
			ArrivalDate:      seg.ArrivalDate,
			ArrivalTime:      seg.ArrivalTime,
			Duration:         seg.Duration,
		}
	}

	// Route to appropriate email template based on trip type
	switch booking.TripType {
	case appmodels.FlightTripTypeOneWay:
		return nh.emailService.SendOneWayFlightBookingEmail(
			booking.Customer.Email,
			booking.Customer.Name,
			booking.BookingID,
			booking.BookingDate,
			travelers,
			outboundSegments,
			booking.Pricing.BasePrice,
			booking.Pricing.Taxes,
			booking.Pricing.ServiceFee,
			booking.Pricing.Total,
		)

	case appmodels.FlightTripTypeRoundTrip:
		return nh.emailService.SendRoundTripFlightBookingEmail(
			booking.Customer.Email,
			booking.Customer.Name,
			booking.BookingID,
			booking.BookingDate,
			travelers,
			outboundSegments,
			returnSegments,
			booking.Pricing.BasePrice,
			booking.Pricing.Taxes,
			booking.Pricing.ServiceFee,
			booking.Pricing.Total,
		)

	case appmodels.FlightTripTypeMultiCity:
		// For multi-city, combine all segments
		allSegments := append(outboundSegments, returnSegments...)
		return nh.emailService.SendMultiCityFlightBookingEmail(
			booking.Customer.Email,
			booking.Customer.Name,
			booking.BookingID,
			booking.BookingDate,
			travelers,
			allSegments,
			booking.Pricing.BasePrice,
			booking.Pricing.Taxes,
			booking.Pricing.ServiceFee,
			booking.Pricing.Total,
		)

	default:
		// Fallback to original generic email for backward compatibility
		return nh.emailService.SendFlightBookingEmail(
			booking.Customer.Email,
			booking.Customer.Name,
			booking.BookingID,
			booking.BookingDate,
			travelers,
			outboundSegments,
			returnSegments,
			booking.Pricing.BasePrice,
			booking.Pricing.Taxes,
			booking.Pricing.ServiceFee,
			booking.Pricing.Total,
		)
	}
}

// SendCarBookingEmail sends notification for car booking
func (nh *NotificationHelper) SendCarBookingEmail(
	ctx context.Context,
	booking *appmodels.CarBooking,
) error {
	// Use the detailed car booking email template
	return nh.emailService.SendCarBookingEmail(
		booking.Customer.Email,
		booking.Customer.Name,
		booking.BookingID,
		booking.BookingDate,
		booking.CarType,
		booking.Passengers,
		booking.Pickup.Location,
		booking.Pickup.Address,
		booking.Pickup.Date.Format("Monday, January 2, 2006"),
		booking.Pickup.Time,
		booking.Dropoff.Location,
		booking.Dropoff.Address,
		booking.Dropoff.Date.Format("Monday, January 2, 2006"),
		booking.Dropoff.Time,
		booking.Driver.Name,
		booking.Driver.PhoneNumber,
		booking.Driver.LicenseNumber,
		booking.Pricing.RentalPrice,
		booking.Pricing.Taxes,
		booking.Pricing.Total,
	)
}

// SendTransferBookingEmail sends notification for transfer booking
func (nh *NotificationHelper) SendTransferBookingEmail(
	ctx context.Context,
	booking *appmodels.TransferBooking,
) error {
	// Use the detailed transfer booking email template
	return nh.emailService.SendTransferBookingEmail(
		booking.Customer.Email,
		booking.Customer.Name,
		booking.Customer.PhoneNumber,
		booking.BookingID,
		booking.BookingDate,
		booking.Vehicle.Type,
		booking.Vehicle.Model,
		booking.Vehicle.Duration,
		booking.Vehicle.Image,
		booking.Vehicle.Passengers,
		booking.Vehicle.MeetGreet,
		booking.Vehicle.Price,
		booking.Pickup.Name,
		booking.Pickup.Address,
		booking.Pickup.Date.Format("Monday, January 2, 2006"),
		booking.Pickup.Time,
		booking.Dropoff.Name,
		booking.Dropoff.Address,
		booking.Dropoff.Date.Format("Monday, January 2, 2006"),
		booking.Dropoff.Time,
		booking.Pricing.BasePrice,
		booking.Pricing.Taxes,
		booking.Pricing.TaxOnFees,
		booking.Pricing.Total,
	)
}
