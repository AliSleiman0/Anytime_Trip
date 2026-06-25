package superadmin

import (
	"travel/backend/internal/handlers/superadmin"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up super admin routes
func SetupRoutes(router fiber.Router, handler *superadmin.SuperAdminHandler) {
	router.Get("/dashboard", handler.GetDashboard)
	router.Get("/system", handler.ManageSystem)

	// Analytics routes
	router.Get("/analytics", handler.GetAnalytics)
	router.Get("/analytics-frag", handler.GetAnalyticsFragment)

	// Accounting routes
	router.Get("/accounting", handler.GetAccounting)
	router.Get("/accounting-frag", handler.GetAccountingFragment)

	// Manage Agents routes
	router.Get("/manage-agents", handler.GetManageAgents)
	router.Get("/manage-agents-frag", handler.GetManageAgentsFragment)

	// Predefined Answers routes
	router.Get("/predefined-answers", handler.GetPredefinedAnswers)
	router.Get("/predefined-answers-frag", handler.GetPredefinedAnswersFragment)
	router.Post("/predefined-answers", handler.CreatePredefinedAnswer)
	router.Put("/predefined-answers/:id", handler.UpdatePredefinedAnswer)
	router.Put("/predefined-answers/:id/shortcut", handler.UpdatePredefinedAnswerShortcut)
	router.Delete("/predefined-answers/:id", handler.DeletePredefinedAnswer)
}
