package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimitConfig holds rate limit configuration.
type RateLimitConfig struct {
	Max        int
	Expiration time.Duration
	KeyPrefix  string
}

// NewRateLimiter creates a scoped rate limiting middleware.
func NewRateLimiter(cfg RateLimitConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cfg.Max,
		Expiration: cfg.Expiration,

		// IMPORTANT: Scope rate limit buckets by purpose
		KeyGenerator: func(c *fiber.Ctx) string {
			return cfg.KeyPrefix + ":" + c.IP()
		},

		LimitReached: func(c *fiber.Ctx) error {
			// API / JSON requests
			if isAPIRequest(c) {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Too many requests. Please try again later.",
				})
			}

			// HTMX requests
			if c.Get("HX-Request") == "true" {
				c.Set("HX-Redirect", "/admin/error?error=Too+many+requests.+Please+try+again+later.")
				return c.SendStatus(fiber.StatusTooManyRequests)
			}

			// Browser requests
			return c.Redirect("/admin/error?error=Too+many+requests.+Please+try+again+later.")
		},
	})
}

// StrictRateLimiter — brute force protection (login only)
func StrictRateLimiter() fiber.Handler {
	return NewRateLimiter(RateLimitConfig{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyPrefix:  "login",
	})
}

// ModerateRateLimiter — CMS write protection
func ModerateRateLimiter() fiber.Handler {
	return NewRateLimiter(RateLimitConfig{
		Max:        30,
		Expiration: 1 * time.Minute,
		KeyPrefix:  "cms",
	})
}
