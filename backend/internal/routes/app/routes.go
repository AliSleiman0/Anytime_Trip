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
	router.Post("/apple-signin", handler.AppleSignIn)

	// OTP routes
	router.Post("/send-otp", handler.SendOTP)
	router.Post("/verify-otp", handler.VerifyOTP)
	router.Post("/activate-account", handler.ActivateAccount)

	// Password reset routes (email-based OTP)
	router.Post("/forgot-password", handler.ForgotPassword)
	router.Post("/verify-reset-otp", handler.VerifyResetOTP)
	router.Post("/reset-password-email", handler.ResetPasswordWithEmail)

	// Legacy phone-based password reset (keeping for backward compatibility)
	router.Post("/reset-password", handler.ResetPassword)

	// Protected routes (authentication required)
	router.Get("/dashboard", middleware.AuthMiddleware(jwtSecret), handler.GetDashboard)
	router.Get("/profile", middleware.AuthMiddleware(jwtSecret), handler.GetProfile)
	router.Put("/profile/update", middleware.AuthMiddleware(jwtSecret), handler.UpdateProfile)
	router.Post("/profile/upload-image", middleware.AuthMiddleware(jwtSecret), handler.UploadProfileImage)

	// Notification preferences routes
	router.Get("/notification-preferences", middleware.AuthMiddleware(jwtSecret), handler.GetNotificationPreferences)
	router.Post("/notification-preferences", middleware.AuthMiddleware(jwtSecret), handler.SaveNotificationPreferences)

	// Payment methods routes
	router.Post("/payment-methods", middleware.AuthMiddleware(jwtSecret), handler.AddPaymentMethod)
	router.Get("/payment-methods", middleware.AuthMiddleware(jwtSecret), handler.GetPaymentMethods)
	router.Put("/payment-methods/:id", middleware.AuthMiddleware(jwtSecret), handler.UpdatePaymentMethod)
	router.Delete("/payment-methods/:id", middleware.AuthMiddleware(jwtSecret), handler.DeletePaymentMethod)
	router.Post("/payment-methods/:id/set-default", middleware.AuthMiddleware(jwtSecret), handler.SetDefaultPaymentMethod)

	// Security preferences routes
	router.Get("/security-preferences", middleware.AuthMiddleware(jwtSecret), handler.GetSecurityPreferences)
	router.Post("/security-preferences", middleware.AuthMiddleware(jwtSecret), handler.SaveSecurityPreferences)

	// Banners route (public - no auth required)
	router.Get("/banners", handler.GetBanners)
	router.Get("/homepage-banners", handler.GetHomepageBanners)
	router.Get("/search-banners", handler.GetSearchBanners)

	// Popular locations route (public - no auth required)
	router.Get("/popular-locations", handler.GetPopularLocations)

	// Car search route (public - no auth required)
	router.Get("/search-cars", handler.SearchCars)

	// Booking routes (authentication required)
	router.Post("/bookings/hotel", middleware.AuthMiddleware(jwtSecret), handler.CreateHotelBooking)
	router.Post("/bookings/flight", middleware.AuthMiddleware(jwtSecret), handler.CreateFlightBooking)
	router.Post("/bookings/car", middleware.AuthMiddleware(jwtSecret), handler.CreateCarBooking)
	router.Post("/bookings/transfer", middleware.AuthMiddleware(jwtSecret), handler.CreateTransferBooking)

	// Get user's bookings
	router.Get("/my-bookings", middleware.AuthMiddleware(jwtSecret), handler.GetMyBookings)
	router.Get("/my-bookings/car", middleware.AuthMiddleware(jwtSecret), handler.GetMyCarBookings)
	router.Get("/my-bookings/flight", middleware.AuthMiddleware(jwtSecret), handler.GetMyFlightBookings)
	router.Get("/my-bookings/hotel", middleware.AuthMiddleware(jwtSecret), handler.GetMyHotelBookings)
	router.Get("/my-bookings/transfer", middleware.AuthMiddleware(jwtSecret), handler.GetMyTransferBookings)

	// Get specific booking and cancel
	router.Get("/bookings/:type/:id", middleware.AuthMiddleware(jwtSecret), handler.GetBookingByID)
	router.Post("/bookings/:type/:id/cancel", middleware.AuthMiddleware(jwtSecret), handler.CancelBooking)

	// Support ticket routes
	router.Post("/support/tickets", middleware.AuthMiddleware(jwtSecret), handler.CreateSupportTicket)
	router.Get("/support/tickets", middleware.AuthMiddleware(jwtSecret), handler.GetUserSupportTickets)
	router.Get("/support/tickets/:ticketId", middleware.AuthMiddleware(jwtSecret), handler.GetSupportTicketDetails)
	router.Post("/support/tickets/:ticketId/reply", middleware.AuthMiddleware(jwtSecret), handler.ReplyToSupportTicket)
	router.Get("/support/chat/history", middleware.AuthMiddleware(jwtSecret), handler.GetChatHistory)

	// Chatbot rating routes
	router.Post("/chatbot/rating", middleware.AuthMiddleware(jwtSecret), handler.SubmitChatbotRating)
	router.Get("/chatbot/ratings", middleware.AuthMiddleware(jwtSecret), handler.GetUserRatings)
}
