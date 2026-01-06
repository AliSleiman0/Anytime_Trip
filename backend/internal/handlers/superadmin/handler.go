package superadmin

import (
	"Anytime_Travel/backend/internal/repository/superadmin"

	"github.com/gofiber/fiber/v2"
)

// SuperAdminHandler handles super admin-level requests
type SuperAdminHandler struct {
	configRepo *superadmin.SystemConfigRepository
}

func NewSuperAdminHandler(configRepo *superadmin.SystemConfigRepository) *SuperAdminHandler {
	return &SuperAdminHandler{
		configRepo: configRepo,
	}
}

func (h *SuperAdminHandler) GetDashboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Super Admin Dashboard",
	})
}

func (h *SuperAdminHandler) ManageSystem(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Super Admin System Management",
	})
}

// Analytics handlers
func (h *SuperAdminHandler) GetAnalytics(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/full-page/analytics.html")
}

func (h *SuperAdminHandler) GetAnalyticsFragment(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/fragments/analytics-frag.html")
}

// Accounting handlers
func (h *SuperAdminHandler) GetAccounting(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/full-page/accounting.html")
}

func (h *SuperAdminHandler) GetAccountingFragment(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/fragments/accounting-frag.html")
}

// Manage Agents handlers
func (h *SuperAdminHandler) GetManageAgents(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/full-page/manage-agents.html")
}

func (h *SuperAdminHandler) GetManageAgentsFragment(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/fragments/manage-agents-frag.html")
}

// Predefined Answers handlers
func (h *SuperAdminHandler) GetPredefinedAnswers(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/full-page/predefined-answers.html")
}

func (h *SuperAdminHandler) GetPredefinedAnswersFragment(c *fiber.Ctx) error {
	return c.SendFile("./templates/superadmin/fragments/predefined-answers-frag.html")
}
