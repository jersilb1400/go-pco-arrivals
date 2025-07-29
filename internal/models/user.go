package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	PCOUserID    string         `json:"pco_user_id" gorm:"uniqueIndex;not null"`
	Name         string         `json:"name" gorm:"not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	Avatar       string         `json:"avatar"`
	IsAdmin      bool           `json:"is_admin" gorm:"default:false"`
	AccessToken  string         `json:"-" gorm:"not null"`
	RefreshToken string         `json:"-" gorm:"not null"`
	TokenExpiry  time.Time      `json:"token_expiry"`
	LastLogin    time.Time      `json:"last_login"`
	LastActivity time.Time      `json:"last_activity"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Sessions      []Session      `json:"-" gorm:"foreignKey:UserID"`
	Events        []Event        `json:"-" gorm:"foreignKey:CreatedBy"`
	Notifications []Notification `json:"-" gorm:"foreignKey:CreatedBy"`
}

func (User) TableName() string {
	return "users"
}
