package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    ID         string         `json:"id" gorm:"primaryKey"`
    MerchantID string         `json:"merchant_id"`
    Name       string         `json:"name"`
    Email      string         `json:"email"`
    Password   string         `json:"password"`
    Role       string         `json:"role" gorm:"default:user"`
    CreatedAt  time.Time      `json:"created_at"`
    UpdatedAt  time.Time      `json:"updated_at"`
    DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"` 
}
