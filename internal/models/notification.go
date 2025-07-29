package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	PCOCheckInID string         `json:"pco_check_in_id" gorm:"uniqueIndex;not null"`
	ChildName    string         `json:"child_name" gorm:"not null"`
	SecurityCode string         `json:"security_code" gorm:"not null"`
	LocationID   string         `json:"location_id"`
	LocationName string         `json:"location_name"`
	EventID      uint           `json:"event_id" gorm:"not null"`
	EventName    string         `json:"event_name"`
	ParentName   string         `json:"parent_name"`
	ParentPhone  string         `json:"parent_phone"`
	Notes        string         `json:"notes"`
	Status       string         `json:"status" gorm:"default:'active'"`
	ExpiresAt    time.Time      `json:"expires_at"`
	CreatedBy    string         `json:"created_by" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Event Event `json:"event" gorm:"foreignKey:EventID"`
}

func (Notification) TableName() string {
	return "notifications"
}
