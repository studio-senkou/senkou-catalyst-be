package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailActivationToken struct {
	ID        uint32         `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint32         `json:"user_id" gorm:"not null;index"`
	User      User           `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Token     string         `json:"token" gorm:"type:varchar(255);not null;uniqueIndex"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time     `json:"used_at"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
