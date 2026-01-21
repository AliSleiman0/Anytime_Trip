package admin

import (
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetReports renders the reports page
func (h *AdminHandler) GetReports(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "reports.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Reports and Analytics",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetReportsFragment serves the reports fragment for HTMX partial loads
func (h *AdminHandler) GetReportsFragment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch all bookings
	hotelBookings, err := h.hotelBookingRepo.FindAll(ctx, 0, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching hotel bookings: " + err.Error())
	}

	carBookings, err := h.carBookingRepo.FindAll(ctx, 0, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching car bookings: " + err.Error())
	}

	flightBookings, err := h.flightBookingRepo.FindAll(ctx, 0, 0)
	if err != nil {
		return c.Status(500).SendString("Error fetching flight bookings: " + err.Error())
	}

	payments, err := h.paymentRepo.GetAllPayments(ctx)
	if err != nil {
		return c.Status(500).SendString("Error fetching payments: " + err.Error())
	}

	// Define time periods for comparison
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	sixtyDaysAgo := now.AddDate(0, 0, -60)

	// Calculate current period metrics
	var currentRevenue float64
	var currentHotelCount int
	var currentCarCount int
	var currentFlightCount int

	for _, p := range payments {
		if p.BookingDate.After(thirtyDaysAgo) {
			currentRevenue += p.Amount
		}
	}

	for _, b := range hotelBookings {
		if b.CreatedAt.After(thirtyDaysAgo) {
			currentHotelCount++
		}
	}

	for _, b := range carBookings {
		if b.CreatedAt.After(thirtyDaysAgo) {
			currentCarCount++
		}
	}

	for _, b := range flightBookings {
		if b.CreatedAt.After(thirtyDaysAgo) {
			currentFlightCount++
		}
	}

	// Calculate previous period metrics
	var prevRevenue float64
	var prevHotelCount int
	var prevCarCount int
	var prevFlightCount int

	for _, p := range payments {
		if p.BookingDate.After(sixtyDaysAgo) && p.BookingDate.Before(thirtyDaysAgo) {
			prevRevenue += p.Amount
		}
	}

	for _, b := range hotelBookings {
		if b.CreatedAt.After(sixtyDaysAgo) && b.CreatedAt.Before(thirtyDaysAgo) {
			prevHotelCount++
		}
	}

	for _, b := range carBookings {
		if b.CreatedAt.After(sixtyDaysAgo) && b.CreatedAt.Before(thirtyDaysAgo) {
			prevCarCount++
		}
	}

	for _, b := range flightBookings {
		if b.CreatedAt.After(sixtyDaysAgo) && b.CreatedAt.Before(thirtyDaysAgo) {
			prevFlightCount++
		}
	}

	// Helper function to calculate percentage change
	calcPercentage := func(current, previous float64) (float64, bool) {
		if previous == 0 {
			if current > 0 {
				return 100.0, true
			}
			return 0, true
		}
		change := ((current - previous) / previous) * 100
		return change, change >= 0
	}

	revenuePercent, revenuePositive := calcPercentage(currentRevenue, prevRevenue)
	hotelPercent, hotelPositive := calcPercentage(float64(currentHotelCount), float64(prevHotelCount))
	carPercent, carPositive := calcPercentage(float64(currentCarCount), float64(prevCarCount))
	flightPercent, flightPositive := calcPercentage(float64(currentFlightCount), float64(prevFlightCount))

	// Find top hotel (most bookings)
	hotelBookingCounts := make(map[string]int)
	for _, b := range hotelBookings {
		hotelBookingCounts[b.HotelID]++
	}
	topHotelID := ""
	topHotelCount := 0
	for hotelID, count := range hotelBookingCounts {
		if count > topHotelCount {
			topHotelCount = count
			topHotelID = hotelID
		}
	}

	// Find top car (most bookings)
	carBookingCounts := make(map[string]int)
	for _, b := range carBookings {
		carBookingCounts[b.CarID]++
	}
	topCarID := ""
	topCarCount := 0
	for carID, count := range carBookingCounts {
		if count > topCarCount {
			topCarCount = count
			topCarID = carID
		}
	}

	// Find top flight destination (most bookings)
	flightDestCounts := make(map[string]int)
	for _, b := range flightBookings {
		flightDestCounts[b.FlightID]++
	}
	topFlightID := ""
	topFlightCount := 0
	for flightID, count := range flightDestCounts {
		if count > topFlightCount {
			topFlightCount = count
			topFlightID = flightID
		}
	}

	// Fetch top provider details
	topHotelName := "N/A"
	topHotelLocation := "N/A"
	if topHotelID != "" {
		hotel, err := h.hotelRepo.FindByID(ctx, topHotelID)
		if err == nil && hotel != nil {
			topHotelName = hotel.ProviderName
			topHotelLocation = hotel.Location
		}
	}

	topCarName := "N/A"
	topCarLocation := "N/A"
	if topCarID != "" {
		car, err := h.carRepo.FindByID(ctx, topCarID)
		if err == nil && car != nil {
			topCarName = car.CarName
			topCarLocation = car.Location
		}
	}

	topFlightDest := "N/A"
	topFlightLocation := "N/A"
	if topFlightID != "" {
		flight, err := h.flightRepo.FindByID(ctx, topFlightID)
		if err == nil && flight != nil {
			topFlightDest = flight.ArrivalLocation
			topFlightLocation = flight.DepartureLocation
		}
	}

	// Calculate percentage changes for top providers
	prevTopHotelCount := 0
	prevTopCarCount := 0
	prevTopFlightCount := 0

	for _, b := range hotelBookings {
		if b.HotelID == topHotelID && b.CreatedAt.After(sixtyDaysAgo) && b.CreatedAt.Before(thirtyDaysAgo) {
			prevTopHotelCount++
		}
	}

	for _, b := range carBookings {
		if b.CarID == topCarID && b.CreatedAt.After(sixtyDaysAgo) && b.CreatedAt.Before(thirtyDaysAgo) {
			prevTopCarCount++
		}
	}

	for _, b := range flightBookings {
		if b.FlightID == topFlightID && b.CreatedAt.After(sixtyDaysAgo) && b.CreatedAt.Before(thirtyDaysAgo) {
			prevTopFlightCount++
		}
	}

	// Count current period bookings for top providers
	currentTopHotelCount := 0
	currentTopCarCount := 0
	currentTopFlightCount := 0

	for _, b := range hotelBookings {
		if b.HotelID == topHotelID && b.CreatedAt.After(thirtyDaysAgo) {
			currentTopHotelCount++
		}
	}

	for _, b := range carBookings {
		if b.CarID == topCarID && b.CreatedAt.After(thirtyDaysAgo) {
			currentTopCarCount++
		}
	}

	for _, b := range flightBookings {
		if b.FlightID == topFlightID && b.CreatedAt.After(thirtyDaysAgo) {
			currentTopFlightCount++
		}
	}

	topHotelPercent, topHotelPositive := calcPercentage(float64(currentTopHotelCount), float64(prevTopHotelCount))
	topCarPercent, topCarPositive := calcPercentage(float64(currentTopCarCount), float64(prevTopCarCount))
	topFlightPercent, topFlightPositive := calcPercentage(float64(currentTopFlightCount), float64(prevTopFlightCount))

	data := fiber.Map{
		"OverallRevenue":         formatCurrency("$", currentRevenue),
		"OverallRevenuePercent":  fmt.Sprintf("%.2f", revenuePercent),
		"OverallRevenuePositive": revenuePositive,
		"HotelBookingsCount":     currentHotelCount,
		"HotelBookingsPercent":   fmt.Sprintf("%.2f", hotelPercent),
		"HotelBookingsPositive":  hotelPositive,
		"CarBookingsCount":       currentCarCount,
		"CarBookingsPercent":     fmt.Sprintf("%.2f", carPercent),
		"CarBookingsPositive":    carPositive,
		"FlightBookingsCount":    currentFlightCount,
		"FlightBookingsPercent":  fmt.Sprintf("%.2f", flightPercent),
		"FlightBookingsPositive": flightPositive,
		"TopHotelName":           topHotelName,
		"TopHotelLocation":       topHotelLocation,
		"TopHotelBookings":       topHotelCount,
		"TopHotelPercent":        fmt.Sprintf("%.2f", topHotelPercent),
		"TopHotelPositive":       topHotelPositive,
		"TopCarName":             topCarName,
		"TopCarLocation":         topCarLocation,
		"TopCarBookings":         topCarCount,
		"TopCarPercent":          fmt.Sprintf("%.2f", topCarPercent),
		"TopCarPositive":         topCarPositive,
		"TopFlightDest":          topFlightDest,
		"TopFlightLocation":      topFlightLocation,
		"TopFlightBookings":      topFlightCount,
		"TopFlightPercent":       fmt.Sprintf("%.2f", topFlightPercent),
		"TopFlightPositive":      topFlightPositive,
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "reports-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template: " + err.Error())
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}
