package app

import (
	"Anytime_Travel/backend/internal/repository/app"

	"github.com/gofiber/fiber/v2"
)

// AppHandler handles app-level requests
type AppHandler struct {
	userRepo *app.UserRepository
}

func NewAppHandler(userRepo *app.UserRepository) *AppHandler {
	return &AppHandler{
		userRepo: userRepo,
	}
}

func (h *AppHandler) GetDashboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "App Dashboard",
	})
}

func (h *AppHandler) GetProfile(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "App User Profile",
	})
}
