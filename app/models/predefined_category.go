package models

import (
	"time"

	"gorm.io/gorm"
)

type PredefinedCategory struct {
	ID          uint32         `json:"id"          gorm:"primaryKey"`
	Name        string         `json:"name"        gorm:"type:varchar(100);not null"`
	Description string         `json:"description" gorm:"type:text"`
	ImageURL    string         `json:"image_url"   gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"  gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at"  gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-"           gorm:"index"`
}
