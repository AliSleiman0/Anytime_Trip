package admin

import (
	"github.com/gofiber/fiber/v2"
)

// GetPredefinedAnswers returns all predefined answers as JSON
func (h *AdminHandler) GetPredefinedAnswers(c *fiber.Ctx) error {
	answers, err := h.predefinedAnswerRepo.FindAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch predefined answers",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    answers,
	})
}
