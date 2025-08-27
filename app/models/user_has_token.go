package models

import (
	"time"
)

type UserHasToken struct {
	ID        uint      `json:"id"         gorm:"primaryKey"`
	UserID    uint32    `json:"user_id"    gorm:"not null"`
	User      User      `json:"user"       gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Token     string    `json:"token"      gorm:"type:varchar(255);not null;unique"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
