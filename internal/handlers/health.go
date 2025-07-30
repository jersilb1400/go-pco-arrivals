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

func (h *HealthHandler) DetailedHealth(c *fiber.Ctx) error {
	// Check database connection
	var result int
	err := h.db.Raw("SELECT 1").Scan(&result).Error

	status := "healthy"
	if err != nil {
		status = "unhealthy"
	}

	return c.JSON(fiber.Map{
		"status":  status,
		"service": "go_pco_arrivals",
		"database": map[string]interface{}{
			"status": status,
			"error":  err,
		},
	})
}
