package handlers

import (
	"strconv"
	"time"

	"go_pco_arrivals/internal/config"
	"go_pco_arrivals/internal/services"
	"go_pco_arrivals/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BillboardHandler struct {
	config    *config.Config
	db        *gorm.DB
	logger    *utils.Logger
	billboard *services.BillboardService
	pco       *services.PCOService
}

type BillboardStateResponse struct {
	Success bool                     `json:"success"`
	State   *services.BillboardState `json:"state,omitempty"`
	Error   string                   `json:"error,omitempty"`
}

type CheckInsResponse struct {
	Success  bool                      `json:"success"`
	CheckIns []services.CheckInDisplay `json:"check_ins,omitempty"`
	Total    int                       `json:"total,omitempty"`
	Error    string                    `json:"error,omitempty"`
}

type StatsResponse struct {
	Success bool                   `json:"success"`
	Stats   map[string]interface{} `json:"stats,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

type SyncResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func NewBillboardHandler(config *config.Config, db *gorm.DB, logger *utils.Logger, billboard *services.BillboardService, pco *services.PCOService) *BillboardHandler {
	return &BillboardHandler{
		config:    config,
		db:        db,
		logger:    logger,
		billboard: billboard,
		pco:       pco,
	}
}

// GetBillboardState returns the current billboard state for a location
func (h *BillboardHandler) GetBillboardState(c *fiber.Ctx) error {
	locationID := c.Params("locationID")
	if locationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(BillboardStateResponse{
			Success: false,
			Error:   "Location ID is required",
		})
	}

	state, err := h.billboard.GetBillboardState(locationID)
	if err != nil {
		h.logger.Error("Failed to get billboard state", "error", err, "location_id", locationID)
		return c.Status(fiber.StatusInternalServerError).JSON(BillboardStateResponse{
			Success: false,
			Error:   "Failed to get billboard state",
		})
	}

	return c.JSON(BillboardStateResponse{
		Success: true,
		State:   state,
	})
}

// GetRecentCheckIns returns recent check-ins for a location
func (h *BillboardHandler) GetRecentCheckIns(c *fiber.Ctx) error {
	locationID := c.Params("locationID")
	if locationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(CheckInsResponse{
			Success: false,
			Error:   "Location ID is required",
		})
	}

	// Parse limit parameter
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10
	}

	checkIns, err := h.billboard.GetRecentCheckIns(locationID, limit)
	if err != nil {
		h.logger.Error("Failed to get recent check-ins", "error", err, "location_id", locationID)
		return c.Status(fiber.StatusInternalServerError).JSON(CheckInsResponse{
			Success: false,
			Error:   "Failed to get recent check-ins",
		})
	}

	// Get total count for today
	total, err := h.billboard.GetTodayCheckInCount(locationID)
	if err != nil {
		h.logger.Error("Failed to get today's check-in count", "error", err, "location_id", locationID)
		// Don't fail the request, just set total to 0
		total = 0
	}

	return c.JSON(CheckInsResponse{
		Success:  true,
		CheckIns: checkIns,
		Total:    total,
	})
}

// GetCheckInStats returns statistics for check-ins
func (h *BillboardHandler) GetCheckInStats(c *fiber.Ctx) error {
	locationID := c.Params("locationID")
	if locationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(StatsResponse{
			Success: false,
			Error:   "Location ID is required",
		})
	}

	// Parse days parameter
	daysStr := c.Query("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 || days > 365 {
		days = 7
	}

	stats, err := h.billboard.GetCheckInStats(locationID, days)
	if err != nil {
		h.logger.Error("Failed to get check-in stats", "error", err, "location_id", locationID)
		return c.Status(fiber.StatusInternalServerError).JSON(StatsResponse{
			Success: false,
			Error:   "Failed to get check-in statistics",
		})
	}

	return c.JSON(StatsResponse{
		Success: true,
		Stats:   stats,
	})
}

// SyncPCOCheckIns syncs check-ins from PCO for a location
func (h *BillboardHandler) SyncPCOCheckIns(c *fiber.Ctx) error {
	locationID := c.Params("locationID")
	if locationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(SyncResponse{
			Success: false,
			Error:   "Location ID is required",
		})
	}

	// Get current user's access token
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(SyncResponse{
			Success: false,
			Error:   "Authentication required",
		})
	}

	// TODO: Get user from session and use their access token
	// For now, we'll use the configured access token
	accessToken := h.config.PCO.AccessToken
	if accessToken == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(SyncResponse{
			Success: false,
			Error:   "No access token available",
		})
	}

	err := h.billboard.SyncPCOCheckIns(accessToken, locationID)
	if err != nil {
		h.logger.Error("Failed to sync PCO check-ins", "error", err, "location_id", locationID)
		return c.Status(fiber.StatusInternalServerError).JSON(SyncResponse{
			Success: false,
			Error:   "Failed to sync check-ins from PCO",
		})
	}

	return c.JSON(SyncResponse{
		Success: true,
		Message: "Successfully synced check-ins from PCO",
	})
}

// GetLocations returns all available locations
func (h *BillboardHandler) GetLocations(c *fiber.Ctx) error {
	locations, err := h.billboard.GetLocations()
	if err != nil {
		h.logger.Error("Failed to get locations", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to get locations",
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"locations": locations,
	})
}

// AddLocation adds a new location to the system
func (h *BillboardHandler) AddLocation(c *fiber.Ctx) error {
	var request struct {
		PCOLocationID string `json:"pco_location_id"`
		Name          string `json:"name"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if request.PCOLocationID == "" || request.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "PCO Location ID and Name are required",
		})
	}

	err := h.billboard.AddLocation(request.PCOLocationID, request.Name)
	if err != nil {
		h.logger.Error("Failed to add location", "error", err, "pco_location_id", request.PCOLocationID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to add location",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Location added successfully",
	})
}

// GetLocationBillboard returns the billboard for a specific location
func (h *BillboardHandler) GetLocationBillboard(c *fiber.Ctx) error {
	locationID := c.Params("locationID")
	if locationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(BillboardStateResponse{
			Success: false,
			Error:   "Location ID is required",
		})
	}

	state, err := h.billboard.GetLocationBillboard(locationID)
	if err != nil {
		h.logger.Error("Failed to get location billboard", "error", err, "location_id", locationID)
		return c.Status(fiber.StatusInternalServerError).JSON(BillboardStateResponse{
			Success: false,
			Error:   "Failed to get location billboard",
		})
	}

	return c.JSON(BillboardStateResponse{
		Success: true,
		State:   state,
	})
}

// CleanupOldData removes old check-ins and expired sessions
func (h *BillboardHandler) CleanupOldData(c *fiber.Ctx) error {
	// Cleanup old check-ins
	if err := h.billboard.CleanupOldCheckIns(); err != nil {
		h.logger.Error("Failed to cleanup old check-ins", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to cleanup old check-ins",
		})
	}

	// TODO: Add session cleanup when auth service is properly integrated
	// if err := h.auth.CleanupExpiredSessions(); err != nil {
	// 	h.logger.Error("Failed to cleanup expired sessions", "error", err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"success": false,
	// 		"error":   "Failed to cleanup expired sessions",
	// 	})
	// }

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   "Cleanup completed successfully",
		"timestamp": time.Now(),
	})
}

// GetSystemStatus returns the overall system status
func (h *BillboardHandler) GetSystemStatus(c *fiber.Ctx) error {
	// Get locations count
	var locationsCount int64
	h.db.Model(&struct{}{}).Table("locations").Count(&locationsCount)

	// Get check-ins count for today
	var todayCheckIns int64
	todayStart := time.Now().Truncate(24 * time.Hour)
	todayEnd := todayStart.Add(24 * time.Hour)
	h.db.Model(&struct{}{}).Table("check_ins").
		Where("check_in_time >= ? AND check_in_time < ?", todayStart, todayEnd).
		Count(&todayCheckIns)

	// Get active sessions count
	var activeSessions int64
	h.db.Model(&struct{}{}).Table("sessions").
		Where("expires_at > ?", time.Now()).
		Count(&activeSessions)

	status := fiber.Map{
		"success": true,
		"system": fiber.Map{
			"locations": fiber.Map{
				"total": locationsCount,
			},
			"check_ins": fiber.Map{
				"today": todayCheckIns,
			},
			"sessions": fiber.Map{
				"active": activeSessions,
			},
			"timestamp": time.Now(),
		},
	}

	return c.JSON(status)
}
