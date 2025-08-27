package models

import (
	"time"

	"gorm.io/gorm"
)

type UserSubscription struct {

	gorm.Model

	ID            uint32       `json:"id"             gorm:"type:int;primaryKey"`
	UserID        uint32       `json:"user_id"        gorm:"type:bigint;not null"`
	User          User         `json:"user"           gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubID         uint32       `json:"sub_id"         gorm:"type:bigint;not null"`
	Sub           Subscription `json:"sub"            gorm:"foreignKey:SubID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StartedAt     time.Time    `json:"started_at"     gorm:"type:timestamp;not null"`
	ExpiredAt     time.Time    `json:"expired_at"     gorm:"type:timestamp;not null"`
	IsActive      bool         `json:"is_active"      gorm:"type:boolean;default:false"`
	PaymentStatus string       `json:"payment_status" gorm:"type:varchar(50);not null;default:'pending'"`
	CreatedAt     time.Time    `json:"created_at"     gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time    `json:"updated_at"     gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
