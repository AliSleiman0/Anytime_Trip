package superadmin

import (
	"Anytime_Trip/backend/internal/repository/superadmin"

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
