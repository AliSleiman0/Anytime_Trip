package admin

import (
	"Anytime_Travel/backend/internal/handlers/admin"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up admin routes
func SetupRoutes(router fiber.Router, handler *admin.AdminHandler) {
	// Admin landing page
	router.Get("/", handler.ServePage)
	router.Get("/dashboard", handler.GetDashboard)
	router.Get("/dashboard-frag", handler.GetDashboardFragment)
	router.Get("/bookings", handler.GetBookings)
	router.Get("/bookings-frag", handler.GetBookingsFragment)
	router.Get("/sidebar", handler.GetSidebar)
	router.Get("/header", handler.GetHeader)
	router.Get("/users", handler.ManageUsers)
}
