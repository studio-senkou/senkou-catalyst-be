package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint32         `json:"id"          gorm:"type:int;primaryKey"`
	MerchantID string         `json:"merchant_id" gorm:"type:char(32);not null;unique"`
	Name       string         `json:"name"        gorm:"type:varchar(100);not null"`
	Email      string         `json:"email"       gorm:"type:varchar(100);unique;not null"`
	Password   string         `json:"password"    gorm:"type:varchar(255);not null"`
	Role       string         `json:"role"        gorm:"type:varchar(20);not null;default:user"`
	CreatedAt  time.Time      `json:"created_at"  gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `json:"updated_at"  gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"-"           gorm:"type:timestamp;index"`
}
