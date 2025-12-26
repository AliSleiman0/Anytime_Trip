package admin

import (
	"Anytime_Travel/backend/internal/repository/admin"

	"github.com/gofiber/fiber/v2"
)

// AdminHandler handles admin-level requests
type AdminHandler struct {
	adminRepo *admin.AdminRepository
}

func NewAdminHandler(adminRepo *admin.AdminRepository) *AdminHandler {
	return &AdminHandler{
		adminRepo: adminRepo,
	}
}

func (h *AdminHandler) GetDashboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Admin Dashboard",
	})
}

func (h *AdminHandler) ManageUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Admin User Management",
	})
}
