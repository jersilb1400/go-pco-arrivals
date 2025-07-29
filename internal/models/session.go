package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Token        string         `json:"token" gorm:"uniqueIndex;not null"`
	UserID       uint           `json:"user_id" gorm:"not null"`
	ExpiresAt    time.Time      `json:"expires_at" gorm:"not null"`
	IsRememberMe bool           `json:"is_remember_me" gorm:"default:false"`
	UserAgent    string         `json:"user_agent"`
	IPAddress    string         `json:"ip_address"`
	LastActivity time.Time      `json:"last_activity"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (Session) TableName() string {
	return "sessions"
}
