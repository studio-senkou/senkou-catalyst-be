package dtos

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (dto *LoginRequestDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"Email.required":    "Email is required",
		"Email.email":       "Email must be a valid email address",
		"Password.required": "Password is required",
		"Password.min":      "Password must be at least 8 characters",
	}

}

type GeneratedToken struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"token_expiry"`
}

type LoginResponseDTO struct {
	AccessToken        string `json:"access_token"`
	AccessTokenExpiry  string `json:"access_token_expiry"`
	RefreshToken       string `json:"refresh_token"`
	RefreshTokenExpiry string `json:"refresh_token_expiry"`
}

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (dto *RefreshTokenRequestDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"RefreshToken.required": "Refresh token is required",
	}
}
