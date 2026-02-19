package middleware

import (
	"net/url"
	"strings"

	"Anytime_Travel/backend/core/utils"

	"github.com/gofiber/fiber/v2"
)

// decide if the request should get JSON (API/HTMX) or redirect (HTML pages)
func isAPIRequest(c *fiber.Ctx) bool {
	accept := strings.ToLower(c.Get("Accept"))
	// Treat HX requests as API-ish to allow HX-Redirect handling
	if c.Get("HX-Request") == "true" {
		return true
	}
	return strings.HasPrefix(c.Path(), "/api") || strings.Contains(accept, "application/json")
}

// unified unauthorized handling: redirect for HTML, JSON for API
func handleUnauthorized(c *fiber.Ctx, msg string) error {
	if isAPIRequest(c) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": msg})
	}

	// For HTML requests, redirect to login (HTMX uses HX-Redirect)
	if c.Get("HX-Request") == "true" {
		c.Set("HX-Redirect", "/admin/login")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Redirect("/admin/login")
}

// unified forbidden handling for role-based restrictions
func handleForbidden(c *fiber.Ctx, msg string) error {
	if isAPIRequest(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": msg})
	}

	// URL encode the message for query parameter
	encodedMsg := url.QueryEscape(msg)

	if c.Get("HX-Request") == "true" {
		c.Set("HX-Redirect", "/admin/error?error="+encodedMsg)
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Redirect("/admin/error?error=" + encodedMsg)
}

// RBACMiddleware enforces per-path role restrictions.
// Admins are blocked from the provided forbidden prefixes; superadmins can access everything.
func RBACMiddleware(forbiddenPrefixes []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleVal := c.Locals("user_role")
		role, _ := roleVal.(string)

		// Superadmin can access everything.
		if role == "superadmin" {
			return c.Next()
		}

		// Admin blocked from certain paths.
		if role == "admin" {
			path := strings.ToLower(c.Path())
			for _, p := range forbiddenPrefixes {
				if strings.HasPrefix(path, strings.ToLower(p)) {
					return handleForbidden(c, "Super admin access required")
				}
			}
		}

		return c.Next()
	}
}

// AuthMiddleware handles JWT authentication
// Expects JWT token in Authorization header as "Bearer <token>"
// or in cookie named "auth_token" or "admin_token"
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenStr string

		// Try to get token from Authorization header first
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			// Expected format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenStr = parts[1]
			}
		}

		// If not found in header, try cookies
		if tokenStr == "" {
			tokenStr = c.Cookies("auth_token")
		}
		if tokenStr == "" {
			tokenStr = c.Cookies("admin_token")
		}

		// If no token found, return unauthorized
		if tokenStr == "" {
			return handleUnauthorized(c, "Missing authentication token")
		}

		// Parse and validate JWT
		claims, err := utils.ParseJWT(tokenStr, jwtSecret)
		if err != nil {
			return handleUnauthorized(c, "Invalid or expired token")
		}

		// Store claims in context for use in handlers
		c.Locals("user_id", claims.Subject)
		c.Locals("user_role", claims.Role)
		c.Locals("user_email", claims.Subject) // For backward compatibility
		c.Locals("claims", claims)

		return c.Next()
	}
}

// AdminMiddleware checks for admin role
// Must be used after AuthMiddleware
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("user_role")
		if role == nil {
			return handleUnauthorized(c, "Authentication required")
		}

		roleStr, ok := role.(string)
		if !ok || (roleStr != "admin" && roleStr != "superadmin") {
			return handleForbidden(c, "Admin access required")
		}

		return c.Next()
	}
}

// SuperAdminMiddleware checks for super admin role
// Must be used after AuthMiddleware
func SuperAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("user_role")
		if role == nil {
			return handleUnauthorized(c, "Authentication required")
		}

		roleStr, ok := role.(string)
		if !ok || roleStr != "superadmin" {
			return handleForbidden(c, "Super admin access required")
		}

		return c.Next()
	}
}

// OptionalAuthMiddleware attempts to authenticate but doesn't fail if token is missing
// Useful for routes that work differently for authenticated vs unauthenticated users
func OptionalAuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenStr string

		// Try to get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenStr = parts[1]
			}
		}

		// If not found in header, try cookies
		if tokenStr == "" {
			tokenStr = c.Cookies("auth_token")
		}
		if tokenStr == "" {
			tokenStr = c.Cookies("admin_token")
		}

		// If token exists, try to parse it
		if tokenStr != "" {
			claims, err := utils.ParseJWT(tokenStr, jwtSecret)
			if err == nil {
				// Store claims in context
				c.Locals("user_id", claims.Subject)
				c.Locals("user_role", claims.Role)
				c.Locals("claims", claims)
			}
		}

		// Always continue, even if authentication failed
		return c.Next()
	}
}
