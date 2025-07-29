package models

import (
	"time"

	"gorm.io/gorm"
)

type BillboardState struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	EventID       uint           `json:"event_id"`
	EventName     string         `json:"event_name"`
	Date          time.Time      `json:"date"`
	LocationID    string         `json:"location_id"`
	LocationName  string         `json:"location_name"`
	SecurityCodes []string       `json:"security_codes" gorm:"serializer:json"`
	IsActive      bool           `json:"is_active" gorm:"default:false"`
	LastUpdated   time.Time      `json:"last_updated"`
	CreatedBy     string         `json:"created_by" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Event Event `json:"event" gorm:"foreignKey:EventID"`
}

func (BillboardState) TableName() string {
	return "billboard_states"
}
