package app

import (
	"Anytime_Travel/backend/internal/handlers/app"
	"Anytime_Travel/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up app routes
func SetupRoutes(router fiber.Router, handler *app.AppHandler, jwtSecret string) {
	// Public routes (no authentication required)
	router.Post("/signup", handler.Signup)
	router.Post("/login", handler.Login)
	router.Post("/google-signin", handler.GoogleSignIn)

	// OTP routes
	router.Post("/send-otp", handler.SendOTP)
	router.Post("/verify-otp", handler.VerifyOTP)
	router.Post("/activate-account", handler.ActivateAccount)
	router.Post("/reset-password", handler.ResetPassword)

	// Protected routes (authentication required)
	router.Get("/dashboard", middleware.AuthMiddleware(jwtSecret), handler.GetDashboard)
	router.Get("/profile", middleware.AuthMiddleware(jwtSecret), handler.GetProfile)
	router.Put("/profile/update", middleware.AuthMiddleware(jwtSecret), handler.UpdateProfile)
	router.Post("/profile/upload-image", middleware.AuthMiddleware(jwtSecret), handler.UploadProfileImage)
}
