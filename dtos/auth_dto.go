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
