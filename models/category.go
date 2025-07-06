package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID         uint32         `json:"id"          gorm:"type:int;primaryKey"`
	Name       string         `json:"name"        gorm:"type:varchar(100);not null"`
	MerchantID string         `json:"merchant_id" gorm:"type:uuid;not null"`
	Merchant   Merchant       `json:"merchant"    gorm:"foreignKey:MerchantID;references:ID"`
	CreatedAt  time.Time      `json:"created_at"  gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `json:"updated_at"  gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"-"           gorm:"type:timestamp;index"`
}
