package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentTransaction struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PaymentType     string         `json:"payment_type" gorm:"type:varchar(50);not null"`
	PaymentChannel  string         `json:"payment_channel" gorm:"type:varchar(50);not null"`
	FraudStatus     string         `json:"fraud_status" gorm:"type:varchar(20);default:'pending'"`
	Amount          float64        `json:"amount" gorm:"type:decimal(15,2);not null"`
	Currency        string         `json:"currency" gorm:"type:varchar(10);not null;default:'IDR'"`
	Status          string         `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	TransactionID   *string        `json:"transaction_id,omitempty" gorm:"type:varchar(100)"`
	TransactionTime *time.Time     `json:"transaction_time,omitempty" gorm:"type:timestamp"`
	SignatureKey    *string        `json:"signature_key,omitempty" gorm:"type:varchar(255)"`
	ExpiredAt       *time.Time     `json:"expired_at,omitempty" gorm:"type:timestamp"`
	SettledAt       *time.Time     `json:"settled_at,omitempty" gorm:"type:timestamp"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	SubscriptionOrder *SubscriptionOrder `json:"subscription_order,omitempty" gorm:"foreignKey:PaymentTransactionID"`
}
