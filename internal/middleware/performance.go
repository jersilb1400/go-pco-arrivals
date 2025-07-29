package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func PerformanceMonitoring() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		// Log performance metrics
		duration := time.Since(start)
		if duration > 100*time.Millisecond {
			// Log slow requests
			c.Locals("slow_request", true)
		}

		return err
	}
}

func Compression() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
}
