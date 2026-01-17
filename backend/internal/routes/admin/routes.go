package admin

import (
	"Anytime_Travel/backend/internal/handlers/admin"
	"Anytime_Travel/backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up admin routes
func SetupRoutes(router fiber.Router, handler *admin.AdminHandler, jwtSecret string) {
	// Public routes (no authentication required)
	router.Post("/login", handler.HandleLogin)
	router.Get("/login", handler.Login)
	router.Get("/login-frag", handler.GetLoginFragment)
	router.Get("/forget", handler.Forget)
	router.Get("/forget-frag", handler.GetForgetFragment)
	router.Post("/forgot-password", handler.HandleForgotPassword)
	router.Get("/verify", handler.Verify)
	router.Get("/verify-frag", handler.GetVerifyFragment)
	router.Post("/verify-code", handler.HandleVerify)
	router.Get("/reset", handler.Reset)
	router.Get("/reset-frag", handler.GetResetFragment)
	router.Post("/reset-password", handler.HandleResetPassword)

	// Error pages
	router.Get("/error", handler.Error)
	router.Get("/error-frag", handler.GetErrorFragment)

	// Logout route (can be public or protected - works either way)
	router.Post("/logout", handler.HandleLogout)
	router.Get("/logout", handler.HandleLogout)

	// Protected routes (require authentication)
	protected := router.Use(middleware.AuthMiddleware(jwtSecret))
	protected.Use(middleware.AdminMiddleware())
	// Admin can access most pages; superadmin can access all. Block admin on sensitive sections.
	protected.Use(middleware.RBACMiddleware([]string{
		"/admin/accounting",         // accounting pages
		"/admin/manage-agents",      // manage agents pages
		"/admin/predefined-answers", // predefined answers pages
	}))

	// Protected page views (no rate limiting on authenticated pages to avoid redirect loops)
	protected.Get("/", handler.ServePage)
	protected.Get("/dashboard", handler.GetDashboard)
	protected.Get("/dashboard-frag", handler.GetDashboardFragment)
	protected.Get("/bookings", handler.GetBookings)
	protected.Get("/bookings-frag", handler.GetBookingsFragment)
	protected.Get("/view-booking", handler.ViewBooking)
	protected.Get("/view-booking-frag", handler.GetViewBookingFragment)
	protected.Post("/update-booking-email", handler.UpdateBookingEmail)
	protected.Get("/payments", handler.GetPayments)
	protected.Get("/payments-frag", handler.GetPaymentsFragment)
	protected.Get("/reports", handler.GetReports)
	protected.Get("/reports-frag", handler.GetReportsFragment)
	protected.Get("/support", handler.GetSupport)
	protected.Get("/support-frag", handler.GetSupportFragment)
	protected.Get("/settings", handler.GetSettings)
	protected.Get("/settings-frag", handler.GetSettingsFragment)
	protected.Post("/change-password", handler.ChangePassword)
	protected.Post("/change-currency", handler.ChangeCurrency)
	protected.Post("/change-language", handler.ChangeLanguage)
	protected.Get("/notification-preferences", handler.GetNotificationPreferences)
	protected.Post("/notification-preferences", handler.UpdateNotificationPreferences)
	protected.Get("/cms/flights", handler.CMSFlights)
	protected.Get("/cms/flights-frag", handler.GetCMSFlightsFragment)
	protected.Get("/cms/cars", handler.CMSCars)
	protected.Get("/cms/cars-frag", handler.GetCMSCarsFragment)
	protected.Get("/cms/hotels", handler.CMSHotels)
	protected.Get("/cms/hotels-frag", handler.GetCMSHotelsFragment)
	protected.Get("/cms/homepage/travel", handler.CMSTravelExperience)
	protected.Get("/cms/homepage/travel-frag", handler.GetCMSTravelExperienceFragment)
	// CMS travel upload/save endpoints
	protected.Post("/cms/homepage/travel/upload", handler.UploadTravelImage)
	protected.Post("/cms/homepage/travel/save", handler.SaveHomePage)
	protected.Get("/cms/homepage/banner", handler.CMSBanner)
	protected.Get("/cms/homepage/banner-frag", handler.GetCMSBannerFragment)
	// CMS banner upload/save endpoints
	protected.Post("/cms/homepage/banner/upload", handler.UploadBannerImage)
	protected.Post("/cms/homepage/banner/save", handler.SaveBannerImage)
	// CMS popular upload/save endpoints
	protected.Post("/cms/homepage/popular/upload", handler.UploadPopularImage)
	protected.Post("/cms/homepage/popular/save", handler.SavePopular)
	protected.Get("/cms/homepage/popular", handler.CMSPopularLocations)
	protected.Get("/cms/homepage/popular-frag", handler.GetCMSPopularLocationsFragment)
	protected.Get("/sidebar", handler.GetSidebar)
	protected.Get("/header", handler.GetHeader)
	protected.Get("/users", handler.ManageUsers)
	protected.Get("/user-frag", handler.GetUsersFragment)
	protected.Get("/view-user", handler.ViewUser)
	protected.Get("/view-user-frag", handler.GetViewUserFragment)
	protected.Post("/update-user-email", handler.UpdateUserEmail)
	protected.Get("/service-providers", handler.ManageServiceProviders)
	protected.Get("/service-provider-frag", handler.GetServiceProvidersFragment)
	protected.Get("/view-service", handler.ViewService)
	protected.Get("/view-service-frag", handler.GetViewServiceFragment)
	protected.Post("/update-service-profit", handler.UpdateServiceProviderProfitPercent)

	// Support ticket routes
	protected.Get("/support/tickets", handler.GetAllTickets)
	protected.Get("/tickets/:id", handler.GetTicketDetail)
	protected.Post("/support/tickets/:id/reply", handler.ReplyToTicket)
	protected.Patch("/support/tickets/:id/status", handler.UpdateTicketStatus)
	protected.Patch("/support/tickets/:id/priority", handler.UpdateTicketPriority)

	// CMS mutations
	protected.Post("/cms/flights/update", handler.UpdateFlight)
	protected.Delete("/cms/flights/delete", handler.DeleteFlight)
	protected.Post("/cms/cars/update", handler.UpdateCar)
	protected.Delete("/cms/cars/delete", handler.DeleteCar)

}
