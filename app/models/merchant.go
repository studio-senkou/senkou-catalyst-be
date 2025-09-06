package models

import (
	"time"

	"gorm.io/gorm"
)

type Merchant struct {
	ID         string         `json:"id"         gorm:"type:char(16);primaryKey"`
	Name       string         `json:"name"       gorm:"type:varchar(100);not null"`
	Username   string         `json:"username"   gorm:"type:varchar(100);not null;unique"`
	OwnerID    uint32         `json:"owner_id"   gorm:"type:int;not null;unique"`
	Owner      User           `json:"-"          gorm:"foreignKey:OwnerID;references:ID"`
	Categories []Category     `json:"-"          gorm:"foreignKey:MerchantID;references:ID"`
	CreatedAt  time.Time      `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `json:"-"          gorm:"type:timestamp;index"`
}
