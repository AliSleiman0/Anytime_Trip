package app

import (
	"context"
	"time"

	"Anytime_Travel/backend/internal/models/app"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmitRatingRequest struct {
	Rating   int    `json:"rating" validate:"required,min=1,max=5"`
	Feedback string `json:"feedback"`
}

func (h *AppHandler) SubmitChatbotRating(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userEmail := c.Locals("user_email")
	if userEmail == nil {
		userEmail = ""
	}

	userName := c.Locals("user_name")
	if userName == nil {
		userName = ""
	}

	// Parse request body
	var req SubmitRatingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate rating range
	if req.Rating < 1 || req.Rating > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Rating must be between 1 and 5",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create chatbot rating
	rating := &app.ChatbotRating{
		ID:        primitive.NewObjectID(),
		UserID:    userID.(string),
		UserEmail: userEmail.(string),
		UserName:  userName.(string),
		Rating:    req.Rating,
		Feedback:  req.Feedback,
		CreatedAt: time.Now(),
	}

	// Save to database
	insertedID, err := h.chatbotRatingRepo.Create(ctx, rating)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save rating",
		})
	}
	rating.ID = insertedID

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Rating submitted successfully",
		"data":    rating,
	})
}

func (h *AppHandler) GetUserRatings(c *fiber.Ctx) error {
	// Get user ID from JWT context
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get user's ratings
	ratings, err := h.chatbotRatingRepo.GetByUserID(ctx, userID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch ratings",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Ratings fetched successfully",
		"data":    ratings,
	})
}
