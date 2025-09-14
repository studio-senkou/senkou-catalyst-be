package models

import "time"

type OauthAccount struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Provider     string    `json:"provider" gorm:"not null"`
	UserID       uint      `json:"user_id" gorm:"not null;foreignKey:UserID;references:ID"`
	User         User      `json:"user" gorm:"constraint:OnDelete:CASCADE"`
	AccessToken  string    `json:"access_token" gorm:"text;not null"`
	RefreshToken *string   `json:"refresh_token" gorm:"text"`
	TokenExpiry  time.Time `json:"token_expiry"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
