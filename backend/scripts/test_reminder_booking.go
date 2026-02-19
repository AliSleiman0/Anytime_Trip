package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"Anytime_Travel/backend/config"
	"Anytime_Travel/backend/internal/database"
	"Anytime_Travel/backend/internal/models/app"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	dbConn, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Calculate pickup time 5 minutes from now
	pickupTime := time.Now().Add(5 * time.Minute)

	fmt.Printf("Creating test car booking with pickup at: %s\n", pickupTime.Format("2006-01-02 15:04:05"))

	// Create a test car booking
	carBooking := &app.CarBooking{
		ID:        uuid.New().String(),
		BookingID: fmt.Sprintf("CBK-TEST-%d", time.Now().Unix()),
		UserID:    "test-user-123",
		CarID:     "test-car-1",
		Status:    app.CarBookingStatusConfirmed,
		Customer: app.CarCustomer{
			Name:  "Elias Sakr",
			Email: "elias.s.sakr10@gmail.com",
		},
		CarType:    "SUV",
		Passengers: 4,
		Pickup: app.CarPickupDropoff{
			Location: "BEY Airport",
			Address:  "Rafic Hariri International Airport, Beirut",
			Date:     pickupTime,
			Time:     pickupTime.Format("3:04 PM"),
		},
		Dropoff: app.CarPickupDropoff{
			Location: "Beirut Hotel",
			Address:  "Downtown Beirut",
			Date:     pickupTime.Add(24 * time.Hour),
			Time:     pickupTime.Add(24 * time.Hour).Format("3:04 PM"),
		},
		Driver: app.CarDriver{
			Name:          "John Doe",
			PhoneNumber:   "+961 1 234567",
			LicenseNumber: "LB123456",
		},
		Pricing: app.CarPricing{
			RentalPrice: 100.00,
			Taxes:       15.00,
			Total:       115.00,
		},
		BookingDate:   time.Now(),
		Amount:        115.00,
		Currency:      "USD",
		PaymentStatus: app.CarPaymentStatusPaid,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Insert into database
	collection := dbConn.DB.Collection("car_bookings")
	ctx := context.Background()

	result, err := collection.InsertOne(ctx, carBooking)
	if err != nil {
		log.Fatalf("Failed to insert test booking: %v", err)
	}

	fmt.Printf("✅ Test booking created successfully!\n")
	fmt.Printf("   Booking ID: %s\n", carBooking.BookingID)
	fmt.Printf("   Customer: %s (%s)\n", carBooking.Customer.Name, carBooking.Customer.Email)
	fmt.Printf("   Pickup: %s at %s\n", pickupTime.Format("2006-01-02"), carBooking.Pickup.Time)
	fmt.Printf("   MongoDB ID: %v\n", result.InsertedID)
	fmt.Printf("\n⏰ Reminder should be sent in approximately 5 minutes!\n")
	fmt.Printf("   Expected reminder time: %s\n", time.Now().Add(5*time.Minute).Format("15:04:05"))
}
