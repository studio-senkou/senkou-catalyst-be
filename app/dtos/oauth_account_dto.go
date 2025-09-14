package dtos

import "time"

type CreateOAuthAccountDTO struct {
	Provider       string    `json:"provider" validate:"required,oneof=google facebook github"`
	ProviderUserID string    `json:"provider_user_id" validate:"required"`
	AccessToken    string    `json:"access_token" validate:"required"`
	RefreshToken   string    `json:"refresh_token"`
	ExpiresAt      time.Time `json:"expiry"`
}
