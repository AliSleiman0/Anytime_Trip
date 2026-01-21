package admin

import (
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetDashboard renders the dashboard page shell.
func (h *AdminHandler) GetDashboard(c *fiber.Ctx) error {
	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "full-page", "dashboard.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	data := fiber.Map{
		"Title": "Dashboard",
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}

// GetDashboardFragment serves the dashboard fragment for HTMX partial loads.
func (h *AdminHandler) GetDashboardFragment(c *fiber.Ctx) error {
	totalUsers, err := h.userRepo.CountTotalUsers(c.Context())
	if err != nil {
		return c.Status(500).SendString("Error fetching total users")
	}

	newUsersToday, err := h.userRepo.CountNewUsersToday(c.Context())
	if err != nil {
		return c.Status(500).SendString("Error fetching new users today")
	}

	activeUsers, err := h.userRepo.CountActiveUsers(c.Context())
	if err != nil {
		return c.Status(500).SendString("Error fetching active users")
	}

	unactiveUsers := totalUsers - activeUsers

	// Compute percent changes
	now := time.Now()
	// Start and end times for calculations
	startToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endToday := startToday.AddDate(0, 0, 1)

	// Active users today (last 30 days from start of today)
	activeUsersToday, err := h.userRepo.CountNewUsersBetween(c.Context(), startToday.AddDate(0, 0, -30), endToday)
	if err != nil {
		return c.Status(500).SendString("Error fetching active users today")
	}

	// New users and active users: month vs last month
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)

	activeThisMonth, err := h.userRepo.CountNewUsersBetween(c.Context(), startOfMonth.AddDate(0, 0, -30), endToday)
	if err != nil {
		return c.Status(500).SendString("Error fetching active users this month")
	}
	activeLastMonth, err := h.userRepo.CountNewUsersBetween(c.Context(), startOfLastMonth.AddDate(0, 0, -30), startOfMonth)
	if err != nil {
		return c.Status(500).SendString("Error fetching active users last month")
	}

	newUsersThisMonth, err := h.userRepo.CountNewUsersBetween(c.Context(), startOfMonth, endToday)
	if err != nil {
		return c.Status(500).SendString("Error fetching new users this month")
	}
	newUsersLastMonth, err := h.userRepo.CountNewUsersBetween(c.Context(), startOfLastMonth, startOfMonth)
	if err != nil {
		return c.Status(500).SendString("Error fetching new users last month")
	}

	// Helper to format percentage change
	formatPct := func(cur int64, prev int64) string {
		if prev == 0 {
			if cur == 0 {
				return "+0.00%"
			}
			return "+100.00%"
		}
		diff := float64(cur-prev) / float64(prev) * 100.0
		return fmt.Sprintf("%+.2f%%", diff)
	}

	// New users month-over-month percentage
	newUsersMonthPct := formatPct(newUsersThisMonth, newUsersLastMonth)

	// New users today as percentage of total users
	var newUsersTodayPct string
	if totalUsers == 0 || newUsersToday == 0 {
		newUsersTodayPct = "0.00%"
	} else {
		newUsersTodayPct = fmt.Sprintf("%.2f%%", (float64(newUsersToday)/float64(totalUsers))*100)
	}

	// Active users month-over-month percentage
	activeUsersMonthPct := formatPct(activeThisMonth, activeLastMonth)

	// Active users today as percentage of total users
	var activeUsersTodayPct string
	if totalUsers == 0 || activeUsersToday == 0 {
		activeUsersTodayPct = "0.00%"
	} else {
		activeUsersTodayPct = fmt.Sprintf("%.2f%%", (float64(activeUsersToday)/float64(totalUsers))*100)
	}

	// Inactive users calculations
	unactiveUsersToday := totalUsers - activeUsersToday
	unactiveThisMonth := totalUsers - activeThisMonth
	unactiveLastMonth := totalUsers - activeLastMonth

	// Inactive users month-over-month percentage
	unactiveUsersMonthPct := formatPct(unactiveThisMonth, unactiveLastMonth)

	// Inactive users today as percentage of total users
	var unactiveUsersTodayPct string
	if totalUsers == 0 || unactiveUsersToday == 0 {
		unactiveUsersTodayPct = "0.00%"
	} else {
		unactiveUsersTodayPct = fmt.Sprintf("%.2f%%", (float64(unactiveUsersToday)/float64(totalUsers))*100)
	}

	visitsPct := "-11.01%" // keep placeholder for visits

	// Build last N months labels and booking counts for this year and last year
	monthsCount := 7
	labels := make([]string, 0, monthsCount)
	thisYearCounts := make([]int64, 0, monthsCount)
	lastYearCounts := make([]int64, 0, monthsCount)
	// per-type (flights) datasets
	flightThisYearCounts := make([]int64, 0, monthsCount)
	flightLastYearCounts := make([]int64, 0, monthsCount)

	// start from monthsCount-1 months ago to current month
	for i := monthsCount - 1; i >= 0; i-- {
		t := now.AddDate(0, -i, 0)
		monthStart := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, now.Location())
		monthEnd := monthStart.AddDate(0, 1, 0)

		// This year counts (monthStart..monthEnd)
		carCntThis, err := h.carBookingRepo.CountBetween(c.Context(), monthStart, monthEnd)
		if err != nil {
			return c.Status(500).SendString("Error fetching car bookings for chart")
		}
		flightCntThis, err := h.flightBookingRepo.CountBetween(c.Context(), monthStart, monthEnd)
		if err != nil {
			return c.Status(500).SendString("Error fetching flight bookings for chart")
		}
		hotelCntThis, err := h.hotelBookingRepo.CountBetween(c.Context(), monthStart, monthEnd)
		if err != nil {
			return c.Status(500).SendString("Error fetching hotel bookings for chart")
		}

		// Last year same month
		monthStartLast := monthStart.AddDate(-1, 0, 0)
		monthEndLast := monthEnd.AddDate(-1, 0, 0)
		carCntLast, err := h.carBookingRepo.CountBetween(c.Context(), monthStartLast, monthEndLast)
		if err != nil {
			return c.Status(500).SendString("Error fetching car bookings for chart (last year)")
		}
		flightCntLast, err := h.flightBookingRepo.CountBetween(c.Context(), monthStartLast, monthEndLast)
		if err != nil {
			return c.Status(500).SendString("Error fetching flight bookings for chart (last year)")
		}
		hotelCntLast, err := h.hotelBookingRepo.CountBetween(c.Context(), monthStartLast, monthEndLast)
		if err != nil {
			return c.Status(500).SendString("Error fetching hotel bookings for chart (last year)")
		}

		labels = append(labels, monthStart.Format("Jan 2006"))
		thisYearCounts = append(thisYearCounts, carCntThis+flightCntThis+hotelCntThis)
		lastYearCounts = append(lastYearCounts, carCntLast+flightCntLast+hotelCntLast)
		// flight-only
		flightThisYearCounts = append(flightThisYearCounts, flightCntThis)
		flightLastYearCounts = append(flightLastYearCounts, flightCntLast)
	}

	data := fiber.Map{
		"TotalUsers":                totalUsers,
		"NewUsersToday":             newUsersToday,
		"ActiveUsers":               activeUsers,
		"ActiveUsersToday":          activeUsersToday,
		"UnactiveUsers":             unactiveUsers,
		"UnactiveUsersToday":        unactiveUsersToday,
		"Visits":                    71000, // placeholder
		"NewUsersMonthPercent":      newUsersMonthPct,
		"NewUsersTodayPercent":      newUsersTodayPct,
		"ActiveUsersMonthPercent":   activeUsersMonthPct,
		"ActiveUsersTodayPercent":   activeUsersTodayPct,
		"UnactiveUsersMonthPercent": unactiveUsersMonthPct,
		"UnactiveUsersTodayPercent": unactiveUsersTodayPct,
		"VisitsPercent":             visitsPct,
		"ActiveBookingLabels":       labels,
		"ActiveBookingThisYearData": thisYearCounts,
		"ActiveBookingLastYearData": lastYearCounts,
		"FlightsThisYearData":       flightThisYearCounts,
		"FlightsLastYearData":       flightLastYearCounts,
	}

	// Marshal arrays to JSON and expose as template.JS to avoid template-range JS parsing issues
	if b, err := json.Marshal(labels); err == nil {
		data["ActiveBookingLabelsJSON"] = template.JS(string(b))
	} else {
		data["ActiveBookingLabelsJSON"] = template.JS("[]")
	}
	if b, err := json.Marshal(thisYearCounts); err == nil {
		data["ActiveBookingThisYearDataJSON"] = template.JS(string(b))
	} else {
		data["ActiveBookingThisYearDataJSON"] = template.JS("[]")
	}
	if b, err := json.Marshal(lastYearCounts); err == nil {
		data["ActiveBookingLastYearDataJSON"] = template.JS(string(b))
	} else {
		data["ActiveBookingLastYearDataJSON"] = template.JS("[]")
	}
	if b, err := json.Marshal(flightThisYearCounts); err == nil {
		data["FlightsThisYearDataJSON"] = template.JS(string(b))
	} else {
		data["FlightsThisYearDataJSON"] = template.JS("[]")
	}
	if b, err := json.Marshal(flightLastYearCounts); err == nil {
		data["FlightsLastYearDataJSON"] = template.JS(string(b))
	} else {
		data["FlightsLastYearDataJSON"] = template.JS("[]")
	}

	tmpl, err := template.ParseFiles(filepath.Clean(filepath.Join("templates", "admin", "fragments", "dashboard-frag.html")))
	if err != nil {
		return c.Status(500).SendString("Error loading template")
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c, data)
}
