package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           string         `json:"id"            gorm:"type:uuid;primaryKey"`
	MerchantID   string         `json:"merchant_id"   gorm:"type:text;not null"`
	Merchant     Merchant       `json:"merchant"      gorm:"foreignKey:MerchantID;references:ID"`
	CategoryID   uint32         `json:"category_id"   gorm:"type:int;not null"`
	Category     Category       `json:"category"      gorm:"foreignKey:CategoryID;references:ID"`
	Title        string         `json:"title"         gorm:"type:varchar(150);not null"`
	Price        float64        `json:"price"         gorm:"type:decimal(10,2);not null"`
	Description  string         `json:"description"   gorm:"type:text"`
	AffiliateURL string         `json:"affiliate_url" gorm:"type:text;not null"`
	CreatedAt    time.Time      `json:"created_at"    gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `json:"updated_at"    gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"type:timestamp;index"`
}
