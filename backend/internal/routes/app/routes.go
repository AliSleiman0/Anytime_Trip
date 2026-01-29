package app

import (
	"Anytime_Travel/backend/internal/handlers/app"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up app routes
func SetupRoutes(router fiber.Router, handler *app.AppHandler) {
	// Public routes (no authentication required)
	router.Post("/signup", handler.Signup)
	router.Post("/login", handler.Login)

	// OTP routes
	router.Post("/send-otp", handler.SendOTP)
	router.Post("/verify-otp", handler.VerifyOTP)
	router.Post("/activate-account", handler.ActivateAccount)
	router.Post("/reset-password", handler.ResetPassword)

	// Protected routes (authentication required)
	router.Get("/dashboard", handler.GetDashboard)
	router.Get("/profile", handler.GetProfile)
}
