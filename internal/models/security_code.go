package models

import (
	"time"

	"gorm.io/gorm"
)

// SecurityCode represents a security code that can be used to access the billboard
type SecurityCode struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Code      string         `json:"code" gorm:"uniqueIndex;not null"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedBy string         `json:"created_by" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for SecurityCode
func (SecurityCode) TableName() string {
	return "security_codes"
}
