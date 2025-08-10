package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID         uint32         `json:"id"          gorm:"primaryKey;autoIncrement"`
	Name       string         `json:"name"        gorm:"type:varchar(100);not null"`
	MerchantID string         `json:"merchant_id" gorm:"type:char(16);not null"`
	Merchant   Merchant       `json:"-"           gorm:"foreignKey:MerchantID;references:ID"`
	CreatedAt  time.Time      `json:"created_at"  gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `json:"updated_at"  gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"-"           gorm:"type:timestamp;index"`
}
