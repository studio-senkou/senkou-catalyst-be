package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	ID          uint32             `json:"id"           gorm:"type:int;primaryKey"`
	Name        string             `json:"name"         gorm:"type:varchar(100);not null"`
	Price       float32            `json:"price"        gorm:"type:decimal(10,2);not null"`
	Description string             `json:"description"  gorm:"type:text"`
	Duration    int16              `json:"duration"     gorm:"type:int;not null"`
	Plans       []SubscriptionPlan `json:"plans"        gorm:"foreignKey:SubID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt   time.Time          `json:"created_at"   gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time          `json:"updated_at"   gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt     `json:"-"            gorm:"type:timestamp;index"`
}
