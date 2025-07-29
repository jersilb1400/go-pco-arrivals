package handlers

import (
	"fmt"
	"go_pco_arrivals/internal/models"
	"go_pco_arrivals/internal/services"
	"strconv"
	"time"

	"go_pco_arrivals/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type APIHandler struct {
	db                  *gorm.DB
	pcoService          *services.PCOService
	notificationService *services.NotificationService
	billboardService    *services.BillboardService
	websocketHub        *services.WebSocketHub
	logger              *utils.Logger
}

func NewAPIHandler(db *gorm.DB, pcoService *services.PCOService, notificationService *services.NotificationService, billboardService *services.BillboardService, websocketHub *services.WebSocketHub, logger *utils.Logger) *APIHandler {
	return &APIHandler{
		db:                  db,
		pcoService:          pcoService,
		notificationService: notificationService,
		billboardService:    billboardService,
		websocketHub:        websocketHub,
		logger:              logger,
	}
}

// Event endpoints
func (h *APIHandler) GetEvents(c *fiber.Ctx) error {
	// Get query parameters
	date := c.Query("date")

	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Parse the date parameter
	var targetDate time.Time
	if date != "" {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Use YYYY-MM-DD",
			})
		}
		targetDate = parsedDate
	} else {
		// Default to today if no date provided
		targetDate = time.Now()
	}

	// Try to get real events from PCO first
	events, err := h.pcoService.GetEvents(user.AccessToken, targetDate)
	if err != nil {
		h.logger.Warn("Failed to get events from PCO, using mock data", "error", err.Error())

		// Return enhanced mock data when PCO API fails
		mockEvents := []fiber.Map{
			{
				"id":          "mock-event-1",
				"name":        "Sunday Service",
				"date":        targetDate.Format("2006-01-02"),
				"location":    "Main Auditorium",
				"location_id": "main-auditorium",
				"time":        "09:00 AM",
				"start_time":  targetDate.Format("2006-01-02") + "T09:00:00Z",
				"end_time":    targetDate.Format("2006-01-02") + "T10:30:00Z",
				"is_active":   true,
				"created_at":  time.Now().Format(time.RFC3339),
			},
			{
				"id":          "mock-event-2",
				"name":        "Youth Group",
				"date":        targetDate.Format("2006-01-02"),
				"location":    "Youth Room",
				"location_id": "youth-room",
				"time":        "06:00 PM",
				"start_time":  targetDate.Format("2006-01-02") + "T18:00:00Z",
				"end_time":    targetDate.Format("2006-01-02") + "T19:30:00Z",
				"is_active":   true,
				"created_at":  time.Now().Format(time.RFC3339),
			},
			{
				"id":          "mock-event-3",
				"name":        "Wednesday Bible Study",
				"date":        targetDate.Format("2006-01-02"),
				"location":    "Fellowship Hall",
				"location_id": "fellowship-hall",
				"time":        "07:00 PM",
				"start_time":  targetDate.Format("2006-01-02") + "T19:00:00Z",
				"end_time":    targetDate.Format("2006-01-02") + "T20:30:00Z",
				"is_active":   true,
				"created_at":  time.Now().Format(time.RFC3339),
			},
		}

		return c.JSON(fiber.Map{
			"success": true,
			"events":  mockEvents,
			"note":    "Using mock data - PCO API access pending approval",
		})
	}

	// Convert PCO events to response format
	responseEvents := make([]fiber.Map, len(events))
	for i, event := range events {
		responseEvents[i] = fiber.Map{
			"id":          event.ID,
			"name":        event.Name,
			"date":        event.Date.Format("2006-01-02"),
			"location":    event.LocationName,
			"location_id": event.LocationID,
			"time":        event.StartTime.Format("03:04 PM"),
			"start_time":  event.StartTime.Format(time.RFC3339),
			"end_time":    event.EndTime.Format(time.RFC3339),
			"is_active":   event.IsActive,
			"created_at":  time.Now().Format(time.RFC3339),
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"events":  responseEvents,
	})
}

func (h *APIHandler) GetEvent(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get event endpoint"})
}

func (h *APIHandler) CreateEvent(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Create event endpoint"})
}

func (h *APIHandler) UpdateEvent(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Update event endpoint"})
}

func (h *APIHandler) DeleteEvent(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Delete event endpoint"})
}

// Notification endpoints
func (h *APIHandler) GetNotifications(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var notifications []models.Notification
	query := h.db.Where("status = ? AND expires_at > ?", "active", time.Now())

	if err := query.Order("created_at DESC").Find(&notifications).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notifications",
		})
	}

	// Convert to response format
	var responseNotifications []fiber.Map
	for _, notification := range notifications {
		responseNotification := fiber.Map{
			"id":         notification.PCOCheckInID,
			"message":    notification.ChildName + " checked in",
			"type":       "info",
			"created_at": notification.CreatedAt.Format(time.RFC3339),
			"status":     notification.Status,
		}

		// Add additional fields if available
		if notification.SecurityCode != "" {
			responseNotification["security_code"] = notification.SecurityCode
		}
		if notification.LocationName != "" {
			responseNotification["location_name"] = notification.LocationName
		}
		if notification.EventName != "" {
			responseNotification["event_name"] = notification.EventName
		}
		if notification.ParentName != "" {
			responseNotification["parent_name"] = notification.ParentName
		}
		if !notification.ExpiresAt.IsZero() {
			responseNotification["expires_at"] = notification.ExpiresAt.Format(time.RFC3339)
		}

		responseNotifications = append(responseNotifications, responseNotification)
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"notifications": responseNotifications,
	})
}

func (h *APIHandler) CreateNotification(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Create notification endpoint"})
}

func (h *APIHandler) DeleteNotification(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Delete notification endpoint"})
}

// Security Code endpoints
func (h *APIHandler) GetSecurityCodes(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get all active security codes
	var securityCodes []models.SecurityCode
	if err := h.db.Where("is_active = ?", true).Find(&securityCodes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch security codes",
		})
	}

	// Convert to response format
	var codes []fiber.Map
	for _, code := range securityCodes {
		codes = append(codes, fiber.Map{
			"id":         code.ID,
			"code":       code.Code,
			"is_active":  code.IsActive,
			"created_by": code.CreatedBy,
			"created_at": code.CreatedAt.Format(time.RFC3339),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"codes":   codes,
	})
}

func (h *APIHandler) AddSecurityCode(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var request struct {
		Code string `json:"code"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if request.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Security code is required",
		})
	}

	// Check if code already exists
	var existingCode models.SecurityCode
	if err := h.db.Where("code = ?", request.Code).First(&existingCode).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Security code already exists",
		})
	}

	// Create new security code
	securityCode := models.SecurityCode{
		Code:      request.Code,
		IsActive:  true,
		CreatedBy: user.Name,
	}

	if err := h.db.Create(&securityCode).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create security code",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Security code added successfully",
	})
}

func (h *APIHandler) RemoveSecurityCode(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Security code is required",
		})
	}

	// Find and delete the security code
	var securityCode models.SecurityCode
	if err := h.db.Where("code = ?", code).First(&securityCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Security code not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find security code",
		})
	}

	if err := h.db.Delete(&securityCode).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove security code",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Security code removed successfully",
	})
}

// Billboard Control endpoints
func (h *APIHandler) GetBillboardControl(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get the current billboard state
	var billboardState models.BillboardState
	if err := h.db.Where("is_active = ?", true).First(&billboardState).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No active billboard state
			return c.JSON(fiber.Map{
				"success": true,
				"control": fiber.Map{
					"event_id":       "",
					"event_name":     "",
					"location_id":    "",
					"location_name":  "",
					"security_codes": []string{},
					"is_active":      false,
					"last_updated":   time.Now().Format(time.RFC3339),
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch billboard control",
		})
	}

	// Parse security codes from JSON string
	var securityCodes []string
	if len(billboardState.SecurityCodes) > 0 {
		securityCodes = billboardState.SecurityCodes
	}

	return c.JSON(fiber.Map{
		"success": true,
		"control": fiber.Map{
			"event_id":       billboardState.EventID,
			"event_name":     billboardState.EventName,
			"location_id":    billboardState.LocationID,
			"location_name":  billboardState.LocationName,
			"security_codes": securityCodes,
			"is_active":      billboardState.IsActive,
			"last_updated":   billboardState.LastUpdated.Format(time.RFC3339),
		},
	})
}

func (h *APIHandler) LaunchBillboard(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var request struct {
		EventID       string   `json:"event_id"`
		LocationID    string   `json:"location_id"`
		SecurityCodes []string `json:"security_codes"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if request.EventID == "" || request.LocationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event ID and Location ID are required",
		})
	}

	// Get event details (for now, use mock event since we don't have real PCO events)
	// In a real implementation, you would query the events table
	eventName := "Sunday Service"     // Mock event name
	locationName := "Main Auditorium" // Mock location name

	// Get all active security codes from the database
	var securityCodes []models.SecurityCode
	if err := h.db.Where("is_active = ?", true).Find(&securityCodes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch security codes",
		})
	}

	// Extract code strings
	var codeStrings []string
	for _, code := range securityCodes {
		codeStrings = append(codeStrings, code.Code)
	}

	// Deactivate all existing billboard states
	h.db.Model(&models.BillboardState{}).Where("is_active = ?", true).Update("is_active", false)

	// Create new billboard state
	newBillboardState := models.BillboardState{
		EventID:       1, // Mock event ID
		EventName:     eventName,
		Date:          time.Now(),
		LocationID:    request.LocationID,
		LocationName:  locationName,
		SecurityCodes: codeStrings,
		IsActive:      true,
		LastUpdated:   time.Now(),
		CreatedBy:     user.Name,
	}

	if err := h.db.Create(&newBillboardState).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to launch billboard",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Billboard launched successfully",
	})
}

func (h *APIHandler) ClearBillboard(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Deactivate all active billboard states
	if err := h.db.Model(&models.BillboardState{}).Where("is_active = ?", true).Update("is_active", false).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear billboard",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Billboard cleared successfully",
	})
}

// Check-in endpoints
func (h *APIHandler) GetCheckIns(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get query parameters
	locationID := c.Query("location_id")
	since := c.Query("since")

	var sinceTime time.Time
	if since != "" {
		var err error
		sinceTime, err = time.Parse(time.RFC3339, since)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid since parameter format. Use RFC3339 format (e.g., 2023-01-01T00:00:00Z)",
			})
		}
	} else {
		// Default to 24 hours ago
		sinceTime = time.Now().Add(-24 * time.Hour)
	}

	// Get check-ins from PCO
	checkIns, err := h.pcoService.GetCheckIns(user.AccessToken, locationID, sinceTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch check-ins from PCO",
		})
	}

	return c.JSON(fiber.Map{
		"check_ins": checkIns,
	})
}

func (h *APIHandler) GetCheckInsByLocation(c *fiber.Ctx) error {
	locationID := c.Params("locationId")
	if locationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Location ID is required",
		})
	}

	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get query parameters
	since := c.Query("since")

	var sinceTime time.Time
	if since != "" {
		var err error
		sinceTime, err = time.Parse(time.RFC3339, since)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid since parameter format. Use RFC3339 format (e.g., 2023-01-01T00:00:00Z)",
			})
		}
	} else {
		// Default to 24 hours ago
		sinceTime = time.Now().Add(-24 * time.Hour)
	}

	// Get check-ins from PCO for specific location
	checkIns, err := h.pcoService.GetCheckIns(user.AccessToken, locationID, sinceTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch check-ins from PCO",
		})
	}

	return c.JSON(fiber.Map{
		"check_ins":   checkIns,
		"location_id": locationID,
	})
}

func (h *APIHandler) GetCheckInsByEvent(c *fiber.Ctx) error {
	eventID := c.Params("eventId")
	if eventID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event ID is required",
		})
	}

	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// For now, we'll return a placeholder since PCO API doesn't directly support event-based check-ins
	// In a real implementation, you might need to:
	// 1. Get all locations for the event
	// 2. Get check-ins for each location
	// 3. Filter by event time

	return c.JSON(fiber.Map{
		"message":  "Event-based check-ins not yet implemented",
		"event_id": eventID,
	})
}

// Location endpoints
func (h *APIHandler) GetLocations(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get locations from PCO
	locations, err := h.pcoService.GetLocations(user.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch locations from PCO",
		})
	}

	return c.JSON(fiber.Map{
		"locations": locations,
	})
}

func (h *APIHandler) GetLocation(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get location endpoint"})
}

// GetLocationStatus returns detailed status for a specific location
func (h *APIHandler) GetLocationStatus(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get path parameters
	locationId := c.Params("locationId")
	if locationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Location ID is required",
		})
	}

	// Get active notifications for this location
	var notifications []models.Notification
	if err := h.db.Where("location_name = ? AND status = ? AND expires_at > ?", locationId, "active", time.Now()).
		Order("created_at DESC").Find(&notifications).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch location notifications",
		})
	}

	// Get recent check-ins for this location (last 24 hours)
	todayStart := time.Now().AddDate(0, 0, -1)
	var recentCheckIns []models.CheckIn
	if err := h.db.Where("location_id = ? AND check_in_time >= ?", locationId, todayStart).
		Order("check_in_time DESC").Limit(10).Find(&recentCheckIns).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch recent check-ins",
		})
	}

	// Calculate location metrics
	var totalCheckIns int64
	var todayCheckIns int64
	var activeChildren int64

	// Total check-ins
	h.db.Model(&models.CheckIn{}).Where("location_id = ?", locationId).Count(&totalCheckIns)

	// Today's check-ins
	today := time.Now().Truncate(24 * time.Hour)
	h.db.Model(&models.CheckIn{}).Where("location_id = ? AND check_in_time >= ?", locationId, today).Count(&todayCheckIns)

	// Active children (notifications)
	h.db.Model(&models.Notification{}).Where("location_name = ? AND status = ? AND expires_at > ?", locationId, "active", time.Now()).Count(&activeChildren)

	// Get location details from PCO (if available)
	locationName := locationId // Default to ID if no name available
	locations, err := h.pcoService.GetLocations(user.AccessToken)
	if err == nil {
		for _, loc := range locations {
			if loc.ID == locationId {
				locationName = loc.Name
				break
			}
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"location": fiber.Map{
			"id":               locationId,
			"name":             locationName,
			"active_children":  activeChildren,
			"total_check_ins":  totalCheckIns,
			"today_check_ins":  todayCheckIns,
			"notifications":    notifications,
			"recent_check_ins": recentCheckIns,
			"last_updated":     time.Now().Format(time.RFC3339),
		},
	})
}

// GetLocationAnalytics returns comprehensive analytics for a location
func (h *APIHandler) GetLocationAnalytics(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get path parameters
	locationId := c.Params("locationId")
	if locationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Location ID is required",
		})
	}

	// Get query parameters
	daysStr := c.Query("days", "30")
	days := 30 // default to 30 days
	if parsedDays, err := strconv.ParseInt(daysStr, 10, 32); err == nil {
		days = int(parsedDays)
	}

	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Get daily check-in counts for the period
	var dailyStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	rows, err := h.db.Raw(`
		SELECT DATE(check_in_time) as date, COUNT(*) as count 
		FROM check_ins 
		WHERE location_id = ? AND check_in_time >= ? AND check_in_time <= ?
		GROUP BY DATE(check_in_time)
		ORDER BY date
	`, locationId, startDate, endDate).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var stat struct {
				Date  string `json:"date"`
				Count int64  `json:"count"`
			}
			rows.Scan(&stat.Date, &stat.Count)
			dailyStats = append(dailyStats, stat)
		}
	}

	// Get peak hours analysis
	var peakHours []struct {
		Hour  int   `json:"hour"`
		Count int64 `json:"count"`
	}

	peakRows, err := h.db.Raw(`
		SELECT EXTRACT(HOUR FROM check_in_time) as hour, COUNT(*) as count 
		FROM check_ins 
		WHERE location_id = ? AND check_in_time >= ?
		GROUP BY EXTRACT(HOUR FROM check_in_time)
		ORDER BY count DESC
		LIMIT 5
	`, locationId, startDate).Rows()

	if err == nil {
		defer peakRows.Close()
		for peakRows.Next() {
			var stat struct {
				Hour  int   `json:"hour"`
				Count int64 `json:"count"`
			}
			peakRows.Scan(&stat.Hour, &stat.Count)
			peakHours = append(peakHours, stat)
		}
	}

	// Get average wait times (time between check-in and pickup notification)
	var avgWaitTime float64
	h.db.Raw(`
		SELECT AVG(JULIANDAY(n.created_at) - JULIANDAY(c.check_in_time)) * 24 * 60 as avg_wait_minutes
		FROM check_ins c
		JOIN notifications n ON c.pco_check_in_id = n.pco_check_in_id
		WHERE c.location_id = ? AND c.check_in_time >= ?
	`, locationId, startDate).Scan(&avgWaitTime)

	// Get location efficiency metrics
	var totalNotifications int64
	var completedPickups int64
	var expiredNotifications int64

	h.db.Model(&models.Notification{}).Where("location_name = ? AND created_at >= ?", locationId, startDate).Count(&totalNotifications)
	h.db.Model(&models.Notification{}).Where("location_name = ? AND status = ? AND created_at >= ?", locationId, "completed", startDate).Count(&completedPickups)
	h.db.Model(&models.Notification{}).Where("location_name = ? AND expires_at < ? AND created_at >= ?", locationId, time.Now(), startDate).Count(&expiredNotifications)

	efficiencyRate := 0.0
	if totalNotifications > 0 {
		efficiencyRate = float64(completedPickups) / float64(totalNotifications) * 100
	}

	return c.JSON(fiber.Map{
		"success": true,
		"analytics": fiber.Map{
			"location_id":           locationId,
			"period_days":           days,
			"daily_stats":           dailyStats,
			"peak_hours":            peakHours,
			"avg_wait_time_mins":    avgWaitTime,
			"efficiency_rate":       efficiencyRate,
			"total_notifications":   totalNotifications,
			"completed_pickups":     completedPickups,
			"expired_notifications": expiredNotifications,
			"generated_at":          time.Now().Format(time.RFC3339),
		},
	})
}

// GetLocationsOverview returns overview of all locations with their current status
func (h *APIHandler) GetLocationsOverview(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get all locations from PCO
	locations, err := h.pcoService.GetLocations(user.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch locations from PCO",
		})
	}

	// Get active notifications grouped by location
	var notifications []models.Notification
	if err := h.db.Where("status = ? AND expires_at > ?", "active", time.Now()).
		Order("created_at DESC").Find(&notifications).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch active notifications",
		})
	}

	// Group notifications by location
	locationMap := make(map[string][]models.Notification)
	for _, notification := range notifications {
		locationName := notification.LocationName
		if locationName == "" {
			locationName = "Unknown Location"
		}
		locationMap[locationName] = append(locationMap[locationName], notification)
	}

	// Build location overview
	var locationOverviews []fiber.Map
	for _, location := range locations {
		activeChildren := len(locationMap[location.Name])

		// Get today's check-ins for this location
		today := time.Now().Truncate(24 * time.Hour)
		var todayCheckIns int64
		h.db.Model(&models.CheckIn{}).Where("location_id = ? AND check_in_time >= ?", location.ID, today).Count(&todayCheckIns)

		// Get total check-ins for this location
		var totalCheckIns int64
		h.db.Model(&models.CheckIn{}).Where("location_id = ?", location.ID).Count(&totalCheckIns)

		locationOverviews = append(locationOverviews, fiber.Map{
			"id":              location.ID,
			"name":            location.Name,
			"active_children": activeChildren,
			"today_check_ins": todayCheckIns,
			"total_check_ins": totalCheckIns,
			"is_active":       true, // All PCO locations are considered active
			"last_updated":    time.Now().Format(time.RFC3339),
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"locations": locationOverviews,
		"summary": fiber.Map{
			"total_locations":  len(locationOverviews),
			"active_locations": len(locationOverviews), // All locations are considered active for now
			"total_children":   len(notifications),
			"generated_at":     time.Now().Format(time.RFC3339),
		},
	})
}

// GetCheckInStats returns check-in statistics for a location over a specified period
func (h *APIHandler) GetCheckInStats(c *fiber.Ctx) error {
	// Get the current user's access token
	userID := c.Locals("user_id").(uint)
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get path parameters
	locationId := c.Params("locationId")
	if locationId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Location ID is required",
		})
	}

	// Get query parameters
	daysStr := c.Query("days", "7")
	days := 7 // default to 7 days
	if parsedDays, err := strconv.ParseInt(daysStr, 10, 32); err == nil {
		days = int(parsedDays)
	}

	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	// Get total check-ins for the period
	var totalCheckIns int64
	if err := h.db.Model(&models.CheckIn{}).
		Where("location_id = ? AND check_in_time >= ? AND check_in_time <= ?", locationId, startDate, endDate).
		Count(&totalCheckIns).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch check-in statistics",
		})
	}

	// Get today's check-ins
	todayStart := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())
	todayEnd := todayStart.Add(24 * time.Hour)
	var todayCheckIns int64
	if err := h.db.Model(&models.CheckIn{}).
		Where("location_id = ? AND check_in_time >= ? AND check_in_time < ?", locationId, todayStart, todayEnd).
		Count(&todayCheckIns).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch today's check-ins",
		})
	}

	// Get weekly check-ins (last 7 days)
	weekStart := endDate.AddDate(0, 0, -7)
	var weeklyCheckIns int64
	if err := h.db.Model(&models.CheckIn{}).
		Where("location_id = ? AND check_in_time >= ? AND check_in_time <= ?", locationId, weekStart, endDate).
		Count(&weeklyCheckIns).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch weekly check-ins",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"stats": fiber.Map{
			"total_check_ins":  totalCheckIns,
			"today_check_ins":  todayCheckIns,
			"weekly_check_ins": weeklyCheckIns,
			"period_days":      days,
			"location_id":      locationId,
		},
	})
}

// TestWebSocketBroadcast is a test endpoint for triggering WebSocket broadcasts
func (h *APIHandler) TestWebSocketBroadcast(c *fiber.Ctx) error {
	eventType := c.Query("type", "notification_update")

	switch eventType {
	case "notification_update":
		// Simulate a new notification
		notification := map[string]interface{}{
			"id":            "test-123",
			"child_name":    "Test Child",
			"security_code": "1234",
			"location_name": "Test Location",
			"parent_name":   "Test Parent",
			"created_at":    time.Now().Format(time.RFC3339),
		}
		h.websocketHub.Broadcast("notification_update", notification)

	case "billboard_state_change":
		// Simulate billboard state change
		billboardState := map[string]interface{}{
			"is_active":     true,
			"event_name":    "Test Event",
			"location_name": "Test Location",
			"event_id":      "test-event-123",
		}
		h.websocketHub.Broadcast("billboard_state_change", billboardState)

	case "security_code_added":
		// Simulate security code addition
		securityCodeData := map[string]interface{}{
			"code":     "5678",
			"event_id": "test-event-123",
		}
		h.websocketHub.Broadcast("security_code_added", securityCodeData)

	case "security_code_removed":
		// Simulate security code removal
		securityCodeData := map[string]interface{}{
			"code":     "5678",
			"event_id": "test-event-123",
		}
		h.websocketHub.Broadcast("security_code_removed", securityCodeData)

	case "billboard_launched":
		// Simulate billboard launch
		billboardState := map[string]interface{}{
			"is_active":     true,
			"event_name":    "Test Event",
			"location_name": "Test Location",
			"event_id":      "test-event-123",
		}
		h.websocketHub.Broadcast("billboard_launched", billboardState)

	case "billboard_cleared":
		// Simulate billboard clear
		clearData := map[string]interface{}{
			"event_id": "test-event-123",
		}
		h.websocketHub.Broadcast("billboard_cleared", clearData)

	default:
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid event type. Valid types: notification_update, billboard_state_change, security_code_added, security_code_removed, billboard_launched, billboard_cleared",
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"message":   fmt.Sprintf("Broadcasted %s event", eventType),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
