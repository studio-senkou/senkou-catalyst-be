package models

import (
	"time"

	"github.com/google/uuid"
)

type ProductMetricInteraction string

const (
	ProductMetricInteractionView  ProductMetricInteraction = "view"
	ProductMetricInteractionClick ProductMetricInteraction = "click"
)

type UserAgent struct {
	Browser string `json:"browser"`
	OS      string `json:"os"`
}

type ProductMetric struct {
	ID          uint                     `json:"id" gorm:"primaryKey"`
	ProductID   uuid.UUID                `json:"product_id" gorm:"index"`
	Product     Product                  `json:"-" gorm:"foreignKey:ProductID"`
	Origin      string                   `json:"origin" gorm:"type:varchar(20)"`
	UserAgent   UserAgent                `json:"user_agent" gorm:"embedded;embeddedPrefix:ua_"`
	Interaction ProductMetricInteraction `json:"interaction_type" gorm:"varchar(20)"`
	CreatedAt   time.Time                `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time                `json:"updated_at" gorm:"type:timestamp"`
}
