//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"travel/backend/config"
	"travel/backend/internal/database"
	"travel/backend/internal/models/app"
)

// randomLastLogin returns a random timestamp within the past maxDays days.
// Roughly 20% of users will have an empty last login (never logged in).
func randomLastLogin(maxDays int) time.Time {
	if rand.Float64() < 0.2 {
		return time.Time{}
	}
	daysAgo := rand.Intn(maxDays)
	hoursAgo := rand.Intn(24)
	minutesAgo := rand.Intn(60)
	secondsAgo := rand.Intn(60)
	return time.Now().Add(-time.Duration(daysAgo) * 24 * time.Hour).
		Add(-time.Duration(hoursAgo) * time.Hour).
		Add(-time.Duration(minutesAgo) * time.Minute).
		Add(-time.Duration(secondsAgo) * time.Second)
}

// randomLastBooking returns a random timestamp within the past maxDays days.
// Roughly 30% of users will have an empty last booking (never booked).
func randomLastBooking(maxDays int) time.Time {
	if rand.Float64() < 0.3 {
		return time.Time{}
	}
	daysAgo := rand.Intn(maxDays)
	hoursAgo := rand.Intn(24)
	minutesAgo := rand.Intn(60)
	secondsAgo := rand.Intn(60)
	return time.Now().Add(-time.Duration(daysAgo) * 24 * time.Hour).
		Add(-time.Duration(hoursAgo) * time.Hour).
		Add(-time.Duration(minutesAgo) * time.Minute).
		Add(-time.Duration(secondsAgo) * time.Second)
}

// randomTotalBookings returns a random number of bookings (0 to maxBookings).
func randomTotalBookings(maxBookings int) int {
	if rand.Float64() < 0.2 {
		return 0
	}
	return rand.Intn(maxBookings) + 1
}

func main() {
	cfg := config.LoadConfig()
	db, err := database.NewDB(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	collection := db.DB.Collection("users")

	// Seed RNG for randomized timestamps
	rand.Seed(time.Now().UnixNano())

	// Sample users to seed
	users := []app.User{
		{
			ID:            "user_001",
			Name:          "John Doe",
			Email:         "john.doe@example.com",
			PhoneNumber:   "+1234567890",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_002",
			Name:          "Jane Smith",
			Email:         "jane.smith@example.com",
			PhoneNumber:   "+0987654321",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-72 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_003",
			Name:          "Alice Johnson",
			Email:         "alice.johnson@example.com",
			PhoneNumber:   "+1122334455",
			IsActive:      false,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-7 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_004",
			Name:          "Bob Williams",
			Email:         "bob.williams@example.com",
			PhoneNumber:   "+14085550101",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-15 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_005",
			Name:          "Charlie Brown",
			Email:         "charlie.brown@example.com",
			PhoneNumber:   "+14085550102",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-30 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_006",
			Name:          "Diana Prince",
			Email:         "diana.prince@example.com",
			PhoneNumber:   "+14085550103",
			IsActive:      false,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-45 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_007",
			Name:          "Ethan Hunt",
			Email:         "ethan.hunt@example.com",
			PhoneNumber:   "+14085550104",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-60 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_008",
			Name:          "Fiona Gallagher",
			Email:         "fiona.gallagher@example.com",
			PhoneNumber:   "+14085550105",
			IsActive:      false,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-5 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_009",
			Name:          "George Miller",
			Email:         "george.miller@example.com",
			PhoneNumber:   "+14085550106",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-10 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_010",
			Name:          "Hannah Lee",
			Email:         "hannah.lee@example.com",
			PhoneNumber:   "+14085550107",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-20 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_011",
			Name:          "Ivan Petrov",
			Email:         "ivan.petrov@example.com",
			PhoneNumber:   "+14085550108",
			IsActive:      false,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-25 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_012",
			Name:          "Julia Roberts",
			Email:         "julia.roberts@example.com",
			PhoneNumber:   "+14085550109",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-35 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_013",
			Name:          "Kevin Durant",
			Email:         "kevin.durant@example.com",
			PhoneNumber:   "+14085550110",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-40 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_014",
			Name:          "Laura Palmer",
			Email:         "laura.palmer@example.com",
			PhoneNumber:   "+14085550111",
			IsActive:      false,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-50 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_015",
			Name:          "Michael Scott",
			Email:         "michael.scott@example.com",
			PhoneNumber:   "+14085550112",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-65 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_016",
			Name:          "Nina Simone",
			Email:         "nina.simone@example.com",
			PhoneNumber:   "+14085550113",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-75 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_017",
			Name:          "Oscar Wilde",
			Email:         "oscar.wilde@example.com",
			PhoneNumber:   "+14085550114",
			IsActive:      false,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-85 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
		{
			ID:            "user_018",
			Name:          "Paula Abdul",
			Email:         "paula.abdul@example.com",
			PhoneNumber:   "+14085550115",
			IsActive:      true,
			TotalBookings: randomTotalBookings(15),
			LastBooking:   randomLastBooking(60),
			CreatedAt:     time.Now().Add(-90 * 24 * time.Hour),
			LastLogin:     randomLastLogin(90),
		},
	}

	for _, user := range users {
		_, err := collection.InsertOne(context.Background(), user)
		if err != nil {
			log.Printf("Error inserting user %s: %v", user.Name, err)
		} else {
			log.Printf("Successfully inserted user: %s", user.Name)
		}
	}

	log.Println("User seeding completed.")
}
