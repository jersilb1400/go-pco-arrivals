package services

import (
	"fmt"
	"time"

	"go_pco_arrivals/internal/config"
	"go_pco_arrivals/internal/models"
	"go_pco_arrivals/internal/utils"

	"gorm.io/gorm"
)

type BillboardService struct {
	config *config.Config
	db     *gorm.DB
	logger *utils.Logger
	pco    *PCOService
	ws     interface{} // Will be WebSocketService when implemented
}

type BillboardState struct {
	LocationID     string           `json:"location_id"`
	LocationName   string           `json:"location_name"`
	LastUpdated    time.Time        `json:"last_updated"`
	TotalCheckIns  int              `json:"total_check_ins"`
	RecentCheckIns []CheckInDisplay `json:"recent_check_ins"`
	IsOnline       bool             `json:"is_online"`
}

type CheckInDisplay struct {
	ID           string    `json:"id"`
	PersonName   string    `json:"person_name"`
	CheckInTime  time.Time `json:"check_in_time"`
	LocationName string    `json:"location_name"`
	Notes        string    `json:"notes"`
	TimeAgo      string    `json:"time_ago"`
}

type RealTimeUpdate struct {
	Type       string          `json:"type"`
	LocationID string          `json:"location_id"`
	CheckIn    *CheckInDisplay `json:"check_in,omitempty"`
	State      *BillboardState `json:"state,omitempty"`
	Timestamp  time.Time       `json:"timestamp"`
}

func NewBillboardService(config *config.Config, db *gorm.DB, logger *utils.Logger, pco *PCOService, ws interface{}) *BillboardService {
	return &BillboardService{
		config: config,
		db:     db,
		logger: logger,
		pco:    pco,
		ws:     ws,
	}
}

// GetBillboardState retrieves the current billboard state for a location
func (s *BillboardService) GetBillboardState(locationID string) (*BillboardState, error) {
	// Get location info
	var location models.Location
	if err := s.db.Where("pco_location_id = ?", locationID).First(&location).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create location if it doesn't exist
			location = models.Location{
				PCOLocationID: locationID,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			if err := s.db.Create(&location).Error; err != nil {
				return nil, fmt.Errorf("failed to create location: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get location: %w", err)
		}
	}

	// Get recent check-ins for this location
	recentCheckIns, err := s.GetRecentCheckIns(locationID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent check-ins: %w", err)
	}

	// Get total check-ins for today
	totalCheckIns, err := s.GetTodayCheckInCount(locationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total check-ins: %w", err)
	}

	state := &BillboardState{
		LocationID:     locationID,
		LocationName:   location.Name,
		LastUpdated:    time.Now(),
		TotalCheckIns:  totalCheckIns,
		RecentCheckIns: recentCheckIns,
		IsOnline:       true, // TODO: Implement actual online status check
	}

	// Save state to database
	if err := s.SaveBillboardState(state); err != nil {
		s.logger.Error("Failed to save billboard state", "error", err, "location_id", locationID)
	}

	return state, nil
}

// GetRecentCheckIns retrieves recent check-ins for a location
func (s *BillboardService) GetRecentCheckIns(locationID string, limit int) ([]CheckInDisplay, error) {
	var checkIns []models.CheckIn
	result := s.db.Where("location_id = ? AND check_in_time >= ?",
		locationID, time.Now().AddDate(0, 0, -1)). // Last 24 hours
		Order("check_in_time DESC").
		Limit(limit).
		Find(&checkIns)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get recent check-ins: %w", result.Error)
	}

	displayCheckIns := make([]CheckInDisplay, len(checkIns))
	for i, checkIn := range checkIns {
		displayCheckIns[i] = CheckInDisplay{
			ID:           checkIn.PCOCheckInID,
			PersonName:   checkIn.PersonName,
			CheckInTime:  checkIn.CheckInTime,
			LocationName: checkIn.LocationName,
			Notes:        checkIn.Notes,
			TimeAgo:      s.formatTimeAgo(checkIn.CheckInTime),
		}
	}

	return displayCheckIns, nil
}

// GetTodayCheckInCount gets the total number of check-ins for today
func (s *BillboardService) GetTodayCheckInCount(locationID string) (int, error) {
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var count int64
	result := s.db.Model(&models.CheckIn{}).
		Where("location_id = ? AND check_in_time >= ? AND check_in_time < ?",
			locationID, startOfDay, endOfDay).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to count today's check-ins: %w", result.Error)
	}

	return int(count), nil
}

// ProcessNewCheckIn processes a new check-in and broadcasts updates
func (s *BillboardService) ProcessNewCheckIn(checkIn *models.CheckIn) error {
	// TODO: Broadcast to location-specific room when WebSocket service is ready
	// displayCheckIn := CheckInDisplay{
	// 	ID:           checkIn.PCOCheckInID,
	// 	PersonName:   checkIn.PersonName,
	// 	CheckInTime:  checkIn.CheckInTime,
	// 	LocationName: checkIn.LocationName,
	// 	Notes:        checkIn.Notes,
	// 	TimeAgo:      s.formatTimeAgo(checkIn.CheckInTime),
	// }
	// update := RealTimeUpdate{
	// 	Type:       "new_check_in",
	// 	LocationID: checkIn.LocationID,
	// 	CheckIn:    &displayCheckIn,
	// 	Timestamp:  time.Now(),
	// }
	// if err := s.ws.BroadcastLocationUpdate(checkIn.LocationID, update); err != nil {
	// 	s.logger.Error("Failed to broadcast check-in update", "error", err, "location_id", checkIn.LocationID)
	// }

	// Update billboard state
	if _, err := s.GetBillboardState(checkIn.LocationID); err != nil {
		s.logger.Error("Failed to update billboard state", "error", err, "location_id", checkIn.LocationID)
	}

	return nil
}

// SyncPCOCheckIns syncs check-ins from PCO and processes them
func (s *BillboardService) SyncPCOCheckIns(accessToken string, locationID string) error {
	// Get check-ins from the last hour
	since := time.Now().Add(-1 * time.Hour)

	pcoCheckIns, err := s.pco.GetCheckIns(accessToken, locationID, since)
	if err != nil {
		return fmt.Errorf("failed to get PCO check-ins: %w", err)
	}

	// Process each check-in
	for _, pcoCheckIn := range pcoCheckIns {
		// Check if check-in already exists
		var existingCheckIn models.CheckIn
		result := s.db.Where("pco_check_in_id = ?", pcoCheckIn.ID).First(&existingCheckIn)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// Create new check-in
				newCheckIn := models.CheckIn{
					PCOCheckInID: pcoCheckIn.ID,
					PersonID:     pcoCheckIn.PersonID,
					PersonName:   pcoCheckIn.PersonName,
					LocationID:   pcoCheckIn.LocationID,
					LocationName: pcoCheckIn.LocationName,
					CheckInTime:  pcoCheckIn.CheckedInAt,
					Notes:        pcoCheckIn.Notes,
					Status:       "active",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}

				if err := s.db.Create(&newCheckIn).Error; err != nil {
					s.logger.Error("Failed to create check-in", "error", err, "check_in_id", pcoCheckIn.ID)
					continue
				}

				// Process the new check-in
				if err := s.ProcessNewCheckIn(&newCheckIn); err != nil {
					s.logger.Error("Failed to process new check-in", "error", err, "check_in_id", pcoCheckIn.ID)
				}
			} else {
				s.logger.Error("Failed to query check-in", "error", result.Error, "check_in_id", pcoCheckIn.ID)
			}
		}
	}

	return nil
}

// SaveBillboardState saves the current billboard state to the database
func (s *BillboardService) SaveBillboardState(state *BillboardState) error {
	billboardState := models.BillboardState{
		LocationID:   state.LocationID,
		LocationName: state.LocationName,
		LastUpdated:  state.LastUpdated,
		IsActive:     state.IsOnline,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Upsert the state
	var existing models.BillboardState
	result := s.db.Where("location_id = ?", state.LocationID).First(&existing)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new state
			if err := s.db.Create(&billboardState).Error; err != nil {
				return fmt.Errorf("failed to create billboard state: %w", err)
			}
		} else {
			return fmt.Errorf("failed to query billboard state: %w", result.Error)
		}
	} else {
		// Update existing state
		existing.LastUpdated = state.LastUpdated
		existing.IsActive = state.IsOnline
		existing.UpdatedAt = time.Now()

		if err := s.db.Save(&existing).Error; err != nil {
			return fmt.Errorf("failed to update billboard state: %w", err)
		}
	}

	return nil
}

// GetLocations retrieves all available locations
func (s *BillboardService) GetLocations() ([]models.Location, error) {
	var locations []models.Location
	result := s.db.Find(&locations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get locations: %w", result.Error)
	}
	return locations, nil
}

// AddLocation adds a new location to the system
func (s *BillboardService) AddLocation(pcoLocationID, name string) error {
	location := models.Location{
		PCOLocationID: pcoLocationID,
		Name:          name,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.db.Create(&location).Error; err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}

	return nil
}

// GetLocationBillboard gets the billboard for a specific location
func (s *BillboardService) GetLocationBillboard(locationID string) (*BillboardState, error) {
	return s.GetBillboardState(locationID)
}

// formatTimeAgo formats a time as a human-readable "time ago" string
func (s *BillboardService) formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}

// CleanupOldCheckIns removes check-ins older than the configured retention period
func (s *BillboardService) CleanupOldCheckIns() error {
	retentionDays := 30 // Default to 30 days - TODO: add to config

	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)

	result := s.db.Where("check_in_time < ?", cutoffDate).Delete(&models.CheckIn{})
	if result.Error != nil {
		return fmt.Errorf("failed to cleanup old check-ins: %w", result.Error)
	}

	s.logger.Info("Cleaned up old check-ins", "deleted_count", result.RowsAffected, "cutoff_date", cutoffDate)
	return nil
}

// GetCheckInStats gets statistics for check-ins
func (s *BillboardService) GetCheckInStats(locationID string, days int) (map[string]interface{}, error) {
	if days <= 0 {
		days = 7 // Default to 7 days
	}

	startDate := time.Now().AddDate(0, 0, -days)

	var stats struct {
		TotalCheckIns  int64 `json:"total_check_ins"`
		TodayCheckIns  int64 `json:"today_check_ins"`
		WeeklyCheckIns int64 `json:"weekly_check_ins"`
	}

	// Total check-ins for the period
	if err := s.db.Model(&models.CheckIn{}).
		Where("location_id = ? AND check_in_time >= ?", locationID, startDate).
		Count(&stats.TotalCheckIns).Error; err != nil {
		return nil, fmt.Errorf("failed to get total check-ins: %w", err)
	}

	// Today's check-ins
	todayStart := time.Now().Truncate(24 * time.Hour)
	if err := s.db.Model(&models.CheckIn{}).
		Where("location_id = ? AND check_in_time >= ?", locationID, todayStart).
		Count(&stats.TodayCheckIns).Error; err != nil {
		return nil, fmt.Errorf("failed to get today's check-ins: %w", err)
	}

	// This week's check-ins
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	if err := s.db.Model(&models.CheckIn{}).
		Where("location_id = ? AND check_in_time >= ?", locationID, weekStart).
		Count(&stats.WeeklyCheckIns).Error; err != nil {
		return nil, fmt.Errorf("failed to get weekly check-ins: %w", err)
	}

	return map[string]interface{}{
		"total_check_ins":  stats.TotalCheckIns,
		"today_check_ins":  stats.TodayCheckIns,
		"weekly_check_ins": stats.WeeklyCheckIns,
		"period_days":      days,
		"location_id":      locationID,
	}, nil
}
