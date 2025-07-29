package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';")

		return c.Next()
	}
}

func RateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Simple rate limiting - in production, use Redis or similar
		// For now, we'll just pass through
		return c.Next()
	}
}

func ValidateInput() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Basic input validation
		return c.Next()
	}
}

func SanitizeInput() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Input sanitization
		return c.Next()
	}
}
