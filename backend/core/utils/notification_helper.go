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
	fmt.Printf("[NOTIFICATION] ====== Starting SendCarBookingEmail ======\n")
	fmt.Printf("[NOTIFICATION] BookingID: %s, Email: %s\n", booking.BookingID, booking.Customer.Email)

	// Get user information by email
	user, err := nh.userRepo.FindByEmail(ctx, booking.Customer.Email)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: User with email %s not found, skipping email notification: %v\n", booking.Customer.Email, err)
		return nil // Don't fail if user doesn't exist
	}
	fmt.Printf("[NOTIFICATION] Found user: %s (%s)\n", user.Email, user.Name)

	// Priority 1: Check admin notification preferences
	adminPrefs, err := nh.adminNotificationPrefsRepo.FindAll(ctx)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: Could not fetch admin preferences: %v\n", err)
		// Continue anyway - don't block on admin prefs error
	}

	adminAllowsEmail := false
	if len(adminPrefs) > 0 {
		fmt.Printf("[NOTIFICATION] Checking %d admin preference(s)\n", len(adminPrefs))
		for _, pref := range adminPrefs {
			fmt.Printf("[NOTIFICATION] Admin pref - NewBooking enabled: %v\n", pref.Bookings.NewBooking)
			if pref.Bookings.NewBooking {
				adminAllowsEmail = true
				break
			}
		}
	} else {
		adminAllowsEmail = true
		fmt.Printf("[NOTIFICATION] No admin preferences found, defaulting to allow email\n")
	}
	fmt.Printf("[NOTIFICATION] Admin allows email: %v\n", adminAllowsEmail)

	if !adminAllowsEmail {
		fmt.Printf("[NOTIFICATION] Email blocked by admin preferences for booking %s\n", booking.BookingID)
		return nil
	}

	// Priority 2: Check user notification preferences
	fmt.Printf("[NOTIFICATION] User email preference enabled: %v\n", user.NotificationPreferences.Email)
	if !user.NotificationPreferences.Email {
		fmt.Printf("[NOTIFICATION] Email blocked by user preferences for user %s\n", user.Email)
		return nil
	}

	// Both admin and user allow email - send it
	fmt.Printf("[NOTIFICATION] Sending car booking email to %s for booking %s\n", user.Email, booking.BookingID)

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

// SendCarBookingReminder sends a car booking reminder with preference checking
func (nh *NotificationHelper) SendCarBookingReminder(
	ctx context.Context,
	userEmail string,
	userName string,
	bookingID string,
	carType string,
	pickupLocation string,
	pickupAddress string,
	pickupDate string,
	pickupTime string,
	totalPrice float64,
) error {
	fmt.Printf("[NOTIFICATION] ====== Starting SendCarBookingReminder ======\n")
	fmt.Printf("[NOTIFICATION] Email: %s, BookingID: %s\n", userEmail, bookingID)

	// Get user information by email
	user, err := nh.userRepo.FindByEmail(ctx, userEmail)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: User with email %s not found, skipping reminder: %v\n", userEmail, err)
		return nil
	}

	// Check admin preferences for reminders
	adminPrefs, err := nh.adminNotificationPrefsRepo.FindAll(ctx)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: Could not fetch admin preferences: %v\n", err)
	}

	adminAllowsEmail := false
	if len(adminPrefs) > 0 {
		for _, pref := range adminPrefs {
			fmt.Printf("[NOTIFICATION] Admin pref - BookingReminder enabled: %v\n", pref.Bookings.BookingReminder)
			if pref.Bookings.BookingReminder {
				adminAllowsEmail = true
				break
			}
		}
	} else {
		adminAllowsEmail = true
		fmt.Printf("[NOTIFICATION] No admin preferences found, defaulting to allow reminder\n")
	}

	if !adminAllowsEmail {
		fmt.Printf("[NOTIFICATION] Reminder blocked by admin preferences for booking %s\n", bookingID)
		return nil
	}

	// Check user notification preferences
	if !user.NotificationPreferences.Email {
		fmt.Printf("[NOTIFICATION] Reminder blocked by user preferences for user %s\n", userEmail)
		return nil
	}

	// Both admin and user allow email - send reminder
	fmt.Printf("[NOTIFICATION] Sending car booking reminder to %s for booking %s\n", userEmail, bookingID)
	return nh.emailService.SendCarBookingReminder(
		userEmail,
		userName,
		bookingID,
		carType,
		pickupLocation,
		pickupAddress,
		pickupDate,
		pickupTime,
		totalPrice,
	)
}

// SendFlightBookingReminder sends a flight booking reminder with preference checking
func (nh *NotificationHelper) SendFlightBookingReminder(
	ctx context.Context,
	userEmail string,
	userName string,
	bookingID string,
	airlineName string,
	flightNumber string,
	departureAirport string,
	arrivalAirport string,
	departureDate string,
	departureTime string,
	totalPrice float64,
) error {
	fmt.Printf("[NOTIFICATION] ====== Starting SendFlightBookingReminder ======\n")
	fmt.Printf("[NOTIFICATION] Email: %s, BookingID: %s\n", userEmail, bookingID)

	user, err := nh.userRepo.FindByEmail(ctx, userEmail)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: User with email %s not found, skipping reminder: %v\n", userEmail, err)
		return nil
	}

	adminPrefs, err := nh.adminNotificationPrefsRepo.FindAll(ctx)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: Could not fetch admin preferences: %v\n", err)
	}

	adminAllowsEmail := false
	if len(adminPrefs) > 0 {
		for _, pref := range adminPrefs {
			if pref.Bookings.BookingReminder {
				adminAllowsEmail = true
				break
			}
		}
	} else {
		adminAllowsEmail = true
	}

	if !adminAllowsEmail {
		fmt.Printf("[NOTIFICATION] Reminder blocked by admin preferences for booking %s\n", bookingID)
		return nil
	}

	if !user.NotificationPreferences.Email {
		fmt.Printf("[NOTIFICATION] Reminder blocked by user preferences for user %s\n", userEmail)
		return nil
	}

	fmt.Printf("[NOTIFICATION] Sending flight booking reminder to %s for booking %s\n", userEmail, bookingID)
	return nh.emailService.SendFlightBookingReminder(
		userEmail,
		userName,
		bookingID,
		airlineName,
		flightNumber,
		departureAirport,
		arrivalAirport,
		departureDate,
		departureTime,
		totalPrice,
	)
}

// SendHotelBookingReminder sends a hotel booking reminder with preference checking
func (nh *NotificationHelper) SendHotelBookingReminder(
	ctx context.Context,
	userEmail string,
	userName string,
	bookingID string,
	hotelName string,
	checkInDate string,
	nights int,
	totalPrice float64,
) error {
	fmt.Printf("[NOTIFICATION] ====== Starting SendHotelBookingReminder ======\n")
	fmt.Printf("[NOTIFICATION] Email: %s, BookingID: %s\n", userEmail, bookingID)

	user, err := nh.userRepo.FindByEmail(ctx, userEmail)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: User with email %s not found, skipping reminder: %v\n", userEmail, err)
		return nil
	}

	adminPrefs, err := nh.adminNotificationPrefsRepo.FindAll(ctx)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: Could not fetch admin preferences: %v\n", err)
	}

	adminAllowsEmail := false
	if len(adminPrefs) > 0 {
		for _, pref := range adminPrefs {
			if pref.Bookings.BookingReminder {
				adminAllowsEmail = true
				break
			}
		}
	} else {
		adminAllowsEmail = true
	}

	if !adminAllowsEmail {
		fmt.Printf("[NOTIFICATION] Reminder blocked by admin preferences for booking %s\n", bookingID)
		return nil
	}

	if !user.NotificationPreferences.Email {
		fmt.Printf("[NOTIFICATION] Reminder blocked by user preferences for user %s\n", userEmail)
		return nil
	}

	fmt.Printf("[NOTIFICATION] Sending hotel booking reminder to %s for booking %s\n", userEmail, bookingID)
	return nh.emailService.SendHotelBookingReminder(
		userEmail,
		userName,
		bookingID,
		hotelName,
		checkInDate,
		nights,
		totalPrice,
	)
}

// SendTransferBookingReminder sends a transfer booking reminder with preference checking
func (nh *NotificationHelper) SendTransferBookingReminder(
	ctx context.Context,
	userEmail string,
	userName string,
	bookingID string,
	vehicleType string,
	pickupLocation string,
	pickupAddress string,
	pickupDate string,
	pickupTime string,
	totalPrice float64,
) error {
	fmt.Printf("[NOTIFICATION] ====== Starting SendTransferBookingReminder ======\n")
	fmt.Printf("[NOTIFICATION] Email: %s, BookingID: %s\n", userEmail, bookingID)

	user, err := nh.userRepo.FindByEmail(ctx, userEmail)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: User with email %s not found, skipping reminder: %v\n", userEmail, err)
		return nil
	}

	adminPrefs, err := nh.adminNotificationPrefsRepo.FindAll(ctx)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: Could not fetch admin preferences: %v\n", err)
	}

	adminAllowsEmail := false
	if len(adminPrefs) > 0 {
		for _, pref := range adminPrefs {
			if pref.Bookings.BookingReminder {
				adminAllowsEmail = true
				break
			}
		}
	} else {
		adminAllowsEmail = true
	}

	if !adminAllowsEmail {
		fmt.Printf("[NOTIFICATION] Reminder blocked by admin preferences for booking %s\n", bookingID)
		return nil
	}

	if !user.NotificationPreferences.Email {
		fmt.Printf("[NOTIFICATION] Reminder blocked by user preferences for user %s\n", userEmail)
		return nil
	}

	fmt.Printf("[NOTIFICATION] Sending transfer booking reminder to %s for booking %s\n", userEmail, bookingID)
	return nh.emailService.SendTransferBookingReminder(
		userEmail,
		userName,
		bookingID,
		vehicleType,
		pickupLocation,
		pickupAddress,
		pickupDate,
		pickupTime,
		totalPrice,
	)
}

// SendCarBookingCancellation sends a car booking cancellation email with preference checking
func (nh *NotificationHelper) SendCarBookingCancellation(
	ctx context.Context,
	userEmail string,
	userName string,
	bookingID string,
	carType string,
	pickupLocation string,
	pickupDate string,
	pickupTime string,
	totalPrice float64,
) error {
	fmt.Printf("[NOTIFICATION] ====== Starting SendCarBookingCancellation ======\n")
	fmt.Printf("[NOTIFICATION] Email: %s, BookingID: %s\n", userEmail, bookingID)

	user, err := nh.userRepo.FindByEmail(ctx, userEmail)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: User with email %s not found, skipping cancellation email: %v\n", userEmail, err)
		return nil
	}

	adminPrefs, err := nh.adminNotificationPrefsRepo.FindAll(ctx)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: Could not fetch admin preferences: %v\n", err)
	}

	adminAllowsEmail := false
	if len(adminPrefs) > 0 {
		for _, pref := range adminPrefs {
			if pref.Bookings.BookingCancelled {
				adminAllowsEmail = true
				break
			}
		}
	} else {
		adminAllowsEmail = true
	}

	if !adminAllowsEmail {
		fmt.Printf("[NOTIFICATION] Cancellation email blocked by admin preferences for booking %s\n", bookingID)
		return nil
	}

	if !user.NotificationPreferences.Email {
		fmt.Printf("[NOTIFICATION] Cancellation email blocked by user preferences for user %s\n", userEmail)
		return nil
	}

	fmt.Printf("[NOTIFICATION] Sending car booking cancellation email to %s for booking %s\n", userEmail, bookingID)
	return nh.emailService.SendCarBookingCancellation(
		userEmail,
		userName,
		bookingID,
		carType,
		pickupLocation,
		pickupDate,
		pickupTime,
		totalPrice,
	)
}

// SendFlightBookingCancellation sends a flight booking cancellation email with preference checking
func (nh *NotificationHelper) SendFlightBookingCancellation(
	ctx context.Context,
	userEmail string,
	userName string,
	bookingID string,
	airlineName string,
	flightNumber string,
	departureAirport string,
	arrivalAirport string,
	departureDate string,
	departureTime string,
	totalPrice float64,
) error {
	fmt.Printf("[NOTIFICATION] ====== Starting SendFlightBookingCancellation ======\n")
	fmt.Printf("[NOTIFICATION] Email: %s, BookingID: %s\n", userEmail, bookingID)

	user, err := nh.userRepo.FindByEmail(ctx, userEmail)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: User with email %s not found, skipping cancellation email: %v\n", userEmail, err)
		return nil
	}

	adminPrefs, err := nh.adminNotificationPrefsRepo.FindAll(ctx)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: Could not fetch admin preferences: %v\n", err)
	}

	adminAllowsEmail := false
	if len(adminPrefs) > 0 {
		for _, pref := range adminPrefs {
			if pref.Bookings.BookingCancelled {
				adminAllowsEmail = true
				break
			}
		}
	} else {
		adminAllowsEmail = true
	}

	if !adminAllowsEmail {
		fmt.Printf("[NOTIFICATION] Cancellation email blocked by admin preferences for booking %s\n", bookingID)
		return nil
	}

	if !user.NotificationPreferences.Email {
		fmt.Printf("[NOTIFICATION] Cancellation email blocked by user preferences for user %s\n", userEmail)
		return nil
	}

	fmt.Printf("[NOTIFICATION] Sending flight booking cancellation email to %s for booking %s\n", userEmail, bookingID)
	return nh.emailService.SendFlightBookingCancellation(
		userEmail,
		userName,
		bookingID,
		airlineName,
		flightNumber,
		departureAirport,
		arrivalAirport,
		departureDate,
		departureTime,
		totalPrice,
	)
}

// SendHotelBookingCancellation sends a hotel booking cancellation email with preference checking
func (nh *NotificationHelper) SendHotelBookingCancellation(
	ctx context.Context,
	userEmail string,
	userName string,
	bookingID string,
	hotelName string,
	checkInDate string,
	nights int,
	totalPrice float64,
) error {
	fmt.Printf("[NOTIFICATION] ====== Starting SendHotelBookingCancellation ======\n")
	fmt.Printf("[NOTIFICATION] Email: %s, BookingID: %s\n", userEmail, bookingID)

	user, err := nh.userRepo.FindByEmail(ctx, userEmail)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: User with email %s not found, skipping cancellation email: %v\n", userEmail, err)
		return nil
	}

	adminPrefs, err := nh.adminNotificationPrefsRepo.FindAll(ctx)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: Could not fetch admin preferences: %v\n", err)
	}

	adminAllowsEmail := false
	if len(adminPrefs) > 0 {
		for _, pref := range adminPrefs {
			if pref.Bookings.BookingCancelled {
				adminAllowsEmail = true
				break
			}
		}
	} else {
		adminAllowsEmail = true
	}

	if !adminAllowsEmail {
		fmt.Printf("[NOTIFICATION] Cancellation email blocked by admin preferences for booking %s\n", bookingID)
		return nil
	}

	if !user.NotificationPreferences.Email {
		fmt.Printf("[NOTIFICATION] Cancellation email blocked by user preferences for user %s\n", userEmail)
		return nil
	}

	fmt.Printf("[NOTIFICATION] Sending hotel booking cancellation email to %s for booking %s\n", userEmail, bookingID)
	return nh.emailService.SendHotelBookingCancellation(
		userEmail,
		userName,
		bookingID,
		hotelName,
		checkInDate,
		nights,
		totalPrice,
	)
}

// SendTransferBookingCancellation sends a transfer booking cancellation email with preference checking
func (nh *NotificationHelper) SendTransferBookingCancellation(
	ctx context.Context,
	userEmail string,
	userName string,
	bookingID string,
	vehicleType string,
	pickupLocation string,
	pickupAddress string,
	pickupDate string,
	pickupTime string,
	totalPrice float64,
) error {
	fmt.Printf("[NOTIFICATION] ====== Starting SendTransferBookingCancellation ======\n")
	fmt.Printf("[NOTIFICATION] Email: %s, BookingID: %s\n", userEmail, bookingID)

	user, err := nh.userRepo.FindByEmail(ctx, userEmail)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: User with email %s not found, skipping cancellation email: %v\n", userEmail, err)
		return nil
	}

	adminPrefs, err := nh.adminNotificationPrefsRepo.FindAll(ctx)
	if err != nil {
		fmt.Printf("[NOTIFICATION] Warning: Could not fetch admin preferences: %v\n", err)
	}

	adminAllowsEmail := false
	if len(adminPrefs) > 0 {
		for _, pref := range adminPrefs {
			if pref.Bookings.BookingCancelled {
				adminAllowsEmail = true
				break
			}
		}
	} else {
		adminAllowsEmail = true
	}

	if !adminAllowsEmail {
		fmt.Printf("[NOTIFICATION] Cancellation email blocked by admin preferences for booking %s\n", bookingID)
		return nil
	}

	if !user.NotificationPreferences.Email {
		fmt.Printf("[NOTIFICATION] Cancellation email blocked by user preferences for user %s\n", userEmail)
		return nil
	}

	fmt.Printf("[NOTIFICATION] Sending transfer booking cancellation email to %s for booking %s\n", userEmail, bookingID)
	return nh.emailService.SendTransferBookingCancellation(
		userEmail,
		userName,
		bookingID,
		vehicleType,
		pickupLocation,
		pickupAddress,
		pickupDate,
		pickupTime,
		totalPrice,
	)
}
