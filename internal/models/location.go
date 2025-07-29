package models

import (
	"time"

	"gorm.io/gorm"
)

type Location struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	PCOLocationID string         `json:"pco_location_id" gorm:"uniqueIndex;not null"`
	Name          string         `json:"name" gorm:"not null"`
	Description   string         `json:"description"`
	Address       string         `json:"address"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	CheckIns []CheckIn `json:"check_ins" gorm:"foreignKey:LocationID"`
}

func (Location) TableName() string {
	return "locations"
}
