package models

import (
	"time"
)

type SubscriptionPlan struct {
	ID           uint32       `json:"id"           gorm:"type:int;primaryKey"`
	SubID        uint32       `json:"sub_id"       gorm:"type:uuid;not null"`
	Subscription Subscription `json:"-" gorm:"foreignKey:SubID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name         string       `json:"name"         gorm:"type:varchar(100);not null"`
	Value        string       `json:"value"        gorm:"type:text;not null"`
	CreatedAt    time.Time    `json:"created_at"   gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time    `json:"updated_at"   gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
