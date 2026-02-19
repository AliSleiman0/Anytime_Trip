package services

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"Anytime_Travel/backend/core/utils"
	"Anytime_Travel/backend/internal/models/app"
)

// ReminderService handles sending reminder emails for upcoming bookings
type ReminderService struct {
	db           *mongo.Database
	emailService *utils.EmailService
	cron         *cron.Cron
}

// NewReminderService creates a new reminder service
func NewReminderService(db *mongo.Database) *ReminderService {
	return &ReminderService{
		db:           db,
		emailService: utils.NewEmailService(),
		cron:         cron.New(),
	}
}

// Start begins the cron job for checking and sending reminders
func (rs *ReminderService) Start() error {
	// Run every hour to check for upcoming bookings
	_, err := rs.cron.AddFunc("@hourly", func() {
		fmt.Println("[REMINDER] Running scheduled reminder check...")
		rs.ProcessReminders()
	})

	if err != nil {
		return fmt.Errorf("failed to schedule reminder job: %w", err)
	}

	rs.cron.Start()
	fmt.Println("[REMINDER] Reminder service started successfully")
	return nil
}

// Stop gracefully stops the reminder service
func (rs *ReminderService) Stop() {
	if rs.cron != nil {
		rs.cron.Stop()
		fmt.Println("[REMINDER] Reminder service stopped")
	}
}

// ProcessReminders checks all booking types and sends reminders
func (rs *ReminderService) ProcessReminders() {
	ctx := context.Background()

	// Calculate when events should occur to trigger reminders
	// Reminder sent 48 hours before the event
	now := time.Now()

	// Window: events between 47-49 hours from now (2-hour window for hourly checks)
	windowStart := now.Add(47 * time.Hour)
	windowEnd := now.Add(49 * time.Hour)

	fmt.Printf("[REMINDER] Checking for events between %s and %s (in ~48 hours)\n",
		windowStart.Format("2006-01-02 15:04"),
		windowEnd.Format("2006-01-02 15:04"))

	// Process each booking type
	rs.processCarBookingReminders(ctx, windowStart, windowEnd)
	rs.processFlightBookingReminders(ctx, windowStart, windowEnd)
	rs.processHotelBookingReminders(ctx, windowStart, windowEnd)
	rs.processTransferBookingReminders(ctx, windowStart, windowEnd)
}

// processCarBookingReminders sends reminders for car bookings
func (rs *ReminderService) processCarBookingReminders(ctx context.Context, windowStart, windowEnd time.Time) {
	collection := rs.db.Collection("car_bookings")

	// Find car bookings with pickup date in the reminder window and not cancelled
	filter := bson.M{
		"pickup.date": bson.M{
			"$gte": windowStart,
			"$lte": windowEnd,
		},
		"status": bson.M{
			"$ne": string(app.CarBookingStatusCancelled),
		},
	}

	fmt.Printf("[REMINDER-DEBUG] Car filter: pickup.date between %s and %s\n",
		windowStart.Format("2006-01-02 15:04:05 MST"),
		windowEnd.Format("2006-01-02 15:04:05 MST"))

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Printf("[REMINDER] Error finding car bookings: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var booking app.CarBooking
		if err := cursor.Decode(&booking); err != nil {
			fmt.Printf("[REMINDER] Error decoding car booking: %v\n", err)
			continue
		}

		fmt.Printf("[REMINDER-DEBUG] Found booking %s with pickup at %s\n",
			booking.BookingID,
			booking.Pickup.Date.Format("2006-01-02 15:04:05 MST"))

		// Send reminder email
		err = rs.emailService.SendCarBookingReminder(
			booking.Customer.Email,
			booking.Customer.Name,
			booking.BookingID,
			booking.CarType,
			booking.Pickup.Location,
			booking.Pickup.Address,
			booking.Pickup.Date.Format("Monday, January 2, 2006"),
			booking.Pickup.Time,
			booking.Pricing.Total,
		)

		if err != nil {
			fmt.Printf("[REMINDER] Error sending car reminder for %s: %v\n", booking.BookingID, err)
		} else {
			fmt.Printf("[REMINDER] Sent car reminder for booking %s to %s\n", booking.BookingID, booking.Customer.Email)
			count++
		}
	}

	fmt.Printf("[REMINDER] Sent %d car booking reminders\n", count)
}

// processFlightBookingReminders sends reminders for flight bookings
func (rs *ReminderService) processFlightBookingReminders(ctx context.Context, windowStart, windowEnd time.Time) {
	collection := rs.db.Collection("flight_bookings")

	// For flights, we need to check the departure date/time of outbound flights
	// Since departure date is stored as string, we need to fetch and parse
	filter := bson.M{
		"status": bson.M{
			"$ne": string(app.FlightBookingStatusCancelled),
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Printf("[REMINDER] Error finding flight bookings: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var booking app.FlightBooking
		if err := cursor.Decode(&booking); err != nil {
			fmt.Printf("[REMINDER] Error decoding flight booking: %v\n", err)
			continue
		}

		// Check if any outbound flight is in the reminder window
		if len(booking.OutboundFlights) > 0 {
			firstSegment := booking.OutboundFlights[0]

			// Parse departure date and time
			departureStr := fmt.Sprintf("%s %s", firstSegment.DepartureDate, firstSegment.DepartureTime)
			departureTime, err := time.Parse("January 2, 2006 3:04 PM", departureStr)
			if err != nil {
				// Try alternative format
				departureTime, err = time.Parse("2006-01-02 15:04", departureStr)
				if err != nil {
					fmt.Printf("[REMINDER] Error parsing flight departure time for %s: %v\n", booking.BookingID, err)
					continue
				}
			}

			// Check if departure is in the reminder window
			if departureTime.After(windowStart) && departureTime.Before(windowEnd) {
				err = rs.emailService.SendFlightBookingReminder(
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

				if err != nil {
					fmt.Printf("[REMINDER] Error sending flight reminder for %s: %v\n", booking.BookingID, err)
				} else {
					fmt.Printf("[REMINDER] Sent flight reminder for booking %s to %s\n", booking.BookingID, booking.Customer.Email)
					count++
				}
			}
		}
	}

	fmt.Printf("[REMINDER] Sent %d flight booking reminders\n", count)
}

// processHotelBookingReminders sends reminders for hotel bookings
func (rs *ReminderService) processHotelBookingReminders(ctx context.Context, windowStart, windowEnd time.Time) {
	collection := rs.db.Collection("hotel_bookings")

	// Find hotel bookings with check-in date in the reminder window and not cancelled
	filter := bson.M{
		"check_in_date": bson.M{
			"$gte": windowStart,
			"$lte": windowEnd,
		},
		"status": bson.M{
			"$ne": string(app.HotelBookingStatusCancelled),
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Printf("[REMINDER] Error finding hotel bookings: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var booking app.HotelBooking
		if err := cursor.Decode(&booking); err != nil {
			fmt.Printf("[REMINDER] Error decoding hotel booking: %v\n", err)
			continue
		}

		// Send reminder email
		err = rs.emailService.SendHotelBookingReminder(
			booking.Customer.Email,
			booking.Customer.Name,
			booking.BookingID,
			booking.HotelName,
			booking.CheckInDate.Format("Monday, January 2, 2006"),
			booking.Nights,
			booking.Pricing.Total,
		)

		if err != nil {
			fmt.Printf("[REMINDER] Error sending hotel reminder for %s: %v\n", booking.BookingID, err)
		} else {
			fmt.Printf("[REMINDER] Sent hotel reminder for booking %s to %s\n", booking.BookingID, booking.Customer.Email)
			count++
		}
	}

	fmt.Printf("[REMINDER] Sent %d hotel booking reminders\n", count)
}

// processTransferBookingReminders sends reminders for transfer bookings
func (rs *ReminderService) processTransferBookingReminders(ctx context.Context, windowStart, windowEnd time.Time) {
	collection := rs.db.Collection("transfer_bookings")

	// Find transfer bookings with pickup date in the reminder window and not cancelled
	filter := bson.M{
		"pickup.date": bson.M{
			"$gte": windowStart,
			"$lte": windowEnd,
		},
		"status": bson.M{
			"$ne": string(app.TransferBookingStatusCancelled),
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Printf("[REMINDER] Error finding transfer bookings: %v\n", err)
		return
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var booking app.TransferBooking
		if err := cursor.Decode(&booking); err != nil {
			fmt.Printf("[REMINDER] Error decoding transfer booking: %v\n", err)
			continue
		}

		// Send reminder email
		err = rs.emailService.SendTransferBookingReminder(
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

		if err != nil {
			fmt.Printf("[REMINDER] Error sending transfer reminder for %s: %v\n", booking.BookingID, err)
		} else {
			fmt.Printf("[REMINDER] Sent transfer reminder for booking %s to %s\n", booking.BookingID, booking.Customer.Email)
			count++
		}
	}

	fmt.Printf("[REMINDER] Sent %d transfer booking reminders\n", count)
}
