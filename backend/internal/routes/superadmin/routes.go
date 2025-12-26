package superadmin

import (
	"Anytime_Travel/backend/internal/handlers/superadmin"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up super admin routes
func SetupRoutes(router fiber.Router, handler *superadmin.SuperAdminHandler) {
	router.Get("/dashboard", handler.GetDashboard)
	router.Get("/system", handler.ManageSystem)
}
