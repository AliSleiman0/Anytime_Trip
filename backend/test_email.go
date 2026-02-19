//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"
	"time"

	"Anytime_Travel/backend/config"
	"Anytime_Travel/backend/core/utils"
	"Anytime_Travel/backend/internal/database"
	"Anytime_Travel/backend/internal/models/app"
	adminrepo "Anytime_Travel/backend/internal/repository/admin"
	apprepo "Anytime_Travel/backend/internal/repository/app"

	"github.com/google/uuid"
)

func main() {
	log.Println("=== Booking Email Notification Test ===")

	cfg := config.LoadConfig()
	db, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := apprepo.NewUserRepository(db.DB)
	adminNotifPrefsRepo := adminrepo.NewNotificationPreferencesRepository(db.DB)
	hotelBookingRepo := apprepo.NewHotelBookingRepository(db.DB)
	flightBookingRepo := apprepo.NewFlightBookingRepository(db.DB)

	// Initialize email service and notification helper
	emailService := utils.NewEmailService()
	notificationHelper := utils.NewNotificationHelper(
		emailService,
		adminNotifPrefsRepo,
		userRepo,
	)

	ctx := context.Background()

	// Step 1: Find or create a test user
	testEmail := "test@example.com"
	log.Printf("\n[1/5] Setting up test user: %s\n", testEmail)

	user, err := userRepo.FindByEmail(ctx, testEmail)
	if err != nil {
		// Create test user
		log.Println("   → Creating new test user...")
		user = &app.User{
			ID:          uuid.New().String(),
			Name:        "Test User",
			Email:       testEmail,
			PhoneNumber: "+1234567890",
			IsActive:    true,
			NotificationPreferences: app.NotificationPreferences{
				Email:   true, // Enable email notifications
				SMS:     false,
				Chatbot: false,
			},
			CreatedAt: time.Now(),
		}
		if err := userRepo.Create(ctx, user); err != nil {
			log.Fatal("   ✗ Failed to create test user:", err)
		}
		log.Println("   ✓ Test user created")
	} else {
		// Update user to enable email notifications
		log.Println("   → Updating existing test user...")
		user.NotificationPreferences.Email = true
		if err := userRepo.Update(ctx, user.ID, user); err != nil {
			log.Fatal("   ✗ Failed to update user:", err)
		}
		log.Println("   ✓ Email notifications enabled for user")
	}

	log.Printf("   User ID: %s\n", user.ID)
	log.Printf("   Email Notifications: %v\n", user.NotificationPreferences.Email)

	// Step 2: Check admin notification preferences
	log.Println("\n[2/5] Checking admin notification preferences...")
	adminPrefs, err := adminNotifPrefsRepo.FindAll(ctx)
	if err != nil {
		log.Printf("   ⚠ Warning: Could not fetch admin preferences: %v\n", err)
		log.Println("   → Emails will be allowed by default")
	} else if len(adminPrefs) == 0 {
		log.Println("   → No admin preferences found, emails allowed by default")
	} else {
		hasNewBookingEnabled := false
		for _, pref := range adminPrefs {
			if pref.Bookings.NewBooking {
				hasNewBookingEnabled = true
				log.Printf("   ✓ Admin %s has new booking notifications enabled\n", pref.AdminID)
				break
			}
		}
		if !hasNewBookingEnabled {
			log.Println("   ✗ WARNING: No admin has enabled new booking notifications!")
			log.Println("   → Emails will be blocked. Enable in admin panel: Settings > Notifications")
		}
	}

	// Step 3: Create test hotel booking
	log.Println("\n[3/5] Creating test hotel booking...")
	hotelBooking := &app.HotelBooking{
		ID:        uuid.New().String(),
		BookingID: "TEST-HB-" + time.Now().Format("20060102150405"),
		UserID:    user.ID,
		HotelID:   "hotel_001",
		Status:    app.HotelBookingStatusConfirmed,
		Customer: app.HotelCustomer{
			Name:  user.Name,
			Email: user.Email,
		},
		Details:       "Test Deluxe Room - 2 Nights",
		BookingDate:   time.Now(),
		Amount:        299.99,
		Currency:      "USD",
		PaymentStatus: app.HotelPaymentStatusPaid,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := hotelBookingRepo.Create(ctx, hotelBooking); err != nil {
		log.Fatal("   ✗ Failed to create hotel booking:", err)
	}
	log.Printf("   ✓ Hotel booking created: %s\n", hotelBooking.BookingID)

	// Step 4: Create test flight booking
	log.Println("\n[4/5] Creating test flight booking...")
	flightBooking := &app.FlightBooking{
		ID:        uuid.New().String(),
		BookingID: "TEST-FB-" + time.Now().Format("20060102150405"),
		UserID:    user.ID,
		FlightID:  "flight_001",
		Status:    app.FlightBookingStatusConfirmed,
		Customer: app.FlightCustomer{
			Name:  user.Name,
			Email: user.Email,
		},
		Details:       "JFK-LAX AA100",
		BookingDate:   time.Now(),
		Amount:        450.00,
		Currency:      "USD",
		PaymentStatus: app.FlightPaymentStatusPaid,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := flightBookingRepo.Create(ctx, flightBooking); err != nil {
		log.Fatal("   ✗ Failed to create flight booking:", err)
	}
	log.Printf("   ✓ Flight booking created: %s\n", flightBooking.BookingID)

	// Step 5: Send email notifications
	log.Println("\n[5/5] Sending email notifications...")

	log.Println("   → Sending hotel booking email...")
	if err := notificationHelper.SendHotelBookingEmail(ctx, hotelBooking); err != nil {
		log.Printf("   ✗ Error sending hotel booking email: %v\n", err)
	} else {
		log.Println("   ✓ Hotel booking email sent successfully!")
	}

	time.Sleep(1 * time.Second) // Brief delay between emails

	log.Println("   → Sending flight booking email...")
	if err := notificationHelper.SendFlightBookingEmail(ctx, flightBooking); err != nil {
		log.Printf("   ✗ Error sending flight booking email: %v\n", err)
	} else {
		log.Println("   ✓ Flight booking email sent successfully!")
	}

	// Summary
	log.Println("\n=== Test Summary ===")
	log.Printf("Test User: %s (%s)\n", user.Name, user.Email)
	log.Printf("Hotel Booking: %s - $%.2f\n", hotelBooking.BookingID, hotelBooking.Amount)
	log.Printf("Flight Booking: %s - $%.2f\n", flightBooking.BookingID, flightBooking.Amount)
	log.Println("\n✓ Test completed!")
	log.Printf("Check the email inbox for: %s\n", user.Email)
	log.Println("(If SENDER_PASSWORD is not set, emails are only logged to console)")
}
