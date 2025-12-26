package app

import (
	"Anytime_Trip/backend/internal/handlers/app"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up app routes
func SetupRoutes(router fiber.Router, handler *app.AppHandler) {
	router.Get("/dashboard", handler.GetDashboard)
	router.Get("/profile", handler.GetProfile)
}
