package admin

import (
	"Anytime_Travel/backend/internal/handlers/admin"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up admin routes
func SetupRoutes(router fiber.Router, handler *admin.AdminHandler) {
	router.Get("/dashboard", handler.GetDashboard)
	router.Get("/users", handler.ManageUsers)
}
