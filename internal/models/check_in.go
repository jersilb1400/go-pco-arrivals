package models

import (
	"time"

	"gorm.io/gorm"
)

type CheckIn struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	PCOCheckInID string         `json:"pco_check_in_id" gorm:"uniqueIndex;not null"`
	PersonID     string         `json:"person_id" gorm:"not null"`
	PersonName   string         `json:"person_name" gorm:"not null"`
	LocationID   string         `json:"location_id" gorm:"not null"`
	LocationName string         `json:"location_name" gorm:"not null"`
	SecurityCode string         `json:"security_code" gorm:"not null"`
	CheckInTime  time.Time      `json:"check_in_time" gorm:"not null"`
	EventID      string         `json:"event_id" gorm:"not null"`
	EventName    string         `json:"event_name"`
	ParentName   string         `json:"parent_name"`
	ParentPhone  string         `json:"parent_phone"`
	Notes        string         `json:"notes"`
	Status       string         `json:"status" gorm:"default:'active'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (CheckIn) TableName() string {
	return "check_ins"
}
