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
	router.Get("/payments", handler.GetPayments)
	router.Get("/payments-frag", handler.GetPaymentsFragment)
	router.Get("/reports", handler.GetReports)
	router.Get("/reports-frag", handler.GetReportsFragment)
	router.Get("/cms/flights", handler.CMSFlights)
	router.Get("/cms/flights-frag", handler.GetCMSFlightsFragment)
	router.Get("/cms/cars", handler.CMSCars)
	router.Get("/cms/cars-frag", handler.GetCMSCarsFragment)
	router.Get("/cms/hotels", handler.CMSHotels)
	router.Get("/cms/hotels-frag", handler.GetCMSHotelsFragment)
	router.Get("/cms/homepage/travel", handler.CMSTravelExperience)
	router.Get("/cms/homepage/travel-frag", handler.GetCMSTravelExperienceFragment)
	router.Get("/cms/homepage/banner", handler.CMSBanner)
	router.Get("/cms/homepage/banner-frag", handler.GetCMSBannerFragment)
	router.Get("/cms/homepage/popular", handler.CMSPopularLocations)
	router.Get("/cms/homepage/popular-frag", handler.GetCMSPopularLocationsFragment)
	router.Get("/sidebar", handler.GetSidebar)
	router.Get("/header", handler.GetHeader)
	router.Get("/users", handler.ManageUsers)
	router.Get("/user-frag", handler.GetUsersFragment)
	router.Get("/view-user", handler.ViewUser)
	router.Get("/view-user-frag", handler.GetViewUserFragment)
	router.Get("/service-providers", handler.ManageServiceProviders)
	router.Get("/service-provider-frag", handler.GetServiceProvidersFragment)
	router.Get("/view-service", handler.ViewService)
	router.Get("/view-service-frag", handler.GetViewServiceFragment)
}
