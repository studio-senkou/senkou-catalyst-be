package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionOrder struct {

	gorm.Model

	ID                   uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID               uint32         `json:"user_id" gorm:"not null;index"`
	SubscriptionID       uint32         `json:"subscription_id" gorm:"not null;index"`
	PaymentTransactionID *uuid.UUID     `json:"payment_transaction_id,omitempty" gorm:"type:uuid;index"`
	Amount               float64        `json:"amount" gorm:"type:decimal(15,2);not null"`
	Status               string         `json:"status" gorm:"type:varchar(20);default:'pending'"`
	CreatedAt            time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt `json:"-" gorm:"index"`

	User               *User               `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Subscription       *Subscription       `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
	PaymentTransaction *PaymentTransaction `json:"payment_transaction,omitempty" gorm:"foreignKey:PaymentTransactionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
