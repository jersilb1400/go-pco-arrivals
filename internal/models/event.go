package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	PCOEventID   string         `json:"pco_event_id" gorm:"uniqueIndex;not null"`
	Name         string         `json:"name" gorm:"not null"`
	Description  string         `json:"description"`
	Date         time.Time      `json:"date" gorm:"not null"`
	StartTime    time.Time      `json:"start_time"`
	EndTime      time.Time      `json:"end_time"`
	LocationID   string         `json:"location_id"`
	LocationName string         `json:"location_name"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedBy    string         `json:"created_by" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Notifications []Notification `json:"notifications" gorm:"foreignKey:EventID"`
	CheckIns      []CheckIn      `json:"check_ins" gorm:"foreignKey:EventID"`
}

func (Event) TableName() string {
	return "events"
}
