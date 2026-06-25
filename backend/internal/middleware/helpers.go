package middleware

import (
	"travel/backend/core/utils"

	"github.com/gofiber/fiber/v2"
)

// GetUserID extracts the user ID from the context
// Returns empty string if not authenticated
func GetUserID(c *fiber.Ctx) string {
	userID := c.Locals("user_id")
	if userID == nil {
		return ""
	}

	str, ok := userID.(string)
	if !ok {
		return ""
	}

	return str
}

// GetUserRole extracts the user role from the context
// Returns empty string if not authenticated
func GetUserRole(c *fiber.Ctx) string {
	role := c.Locals("user_role")
	if role == nil {
		return ""
	}

	str, ok := role.(string)
	if !ok {
		return ""
	}

	return str
}

// GetClaims extracts the JWT claims from the context
// Returns nil if not authenticated
func GetClaims(c *fiber.Ctx) *utils.AdminClaims {
	claims := c.Locals("claims")
	if claims == nil {
		return nil
	}

	adminClaims, ok := claims.(*utils.AdminClaims)
	if !ok {
		return nil
	}

	return adminClaims
}

// IsAuthenticated checks if the current request is authenticated
func IsAuthenticated(c *fiber.Ctx) bool {
	return c.Locals("user_id") != nil
}

// IsAdmin checks if the current user has admin role
func IsAdmin(c *fiber.Ctx) bool {
	role := GetUserRole(c)
	return role == "admin" || role == "superadmin"
}

// IsSuperAdmin checks if the current user has superadmin role
func IsSuperAdmin(c *fiber.Ctx) bool {
	role := GetUserRole(c)
	return role == "superadmin"
}
