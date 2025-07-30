package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "go_pco_arrivals",
	})
}

	// Handle nil db case (MongoDB mode)
	if h.db == nil {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "go_pco_arrivals",
			"database": map[string]interface{}{
				"status": "healthy",
				"type":   "mongodb",
				"note":   "MongoDB mode - no GORM connection available",
			},
		})
	}

	// Check database connection
		},
	})
}
