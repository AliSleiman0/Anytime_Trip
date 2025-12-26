package admin

import (
"github.com/gofiber/fiber/v2"
"Anytime_Trip/backend/internal/handlers/admin"
)

// SetupRoutes sets up admin routes
func SetupRoutes(router fiber.Router, handler *admin.AdminHandler) {
router.Get("/dashboard", handler.GetDashboard)
router.Get("/users", handler.ManageUsers)
}
