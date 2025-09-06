package dtos

type RegisterUserDTO struct {
	Name                 string  `json:"name" validate:"required,min=3,max=100"`
	MerchantUsername     *string `json:"merchant_username,omitempty" validate:"omitempty,min=3,max=100"`
	Email                string  `json:"email" validate:"required,email"`
	Phone                string  `json:"phone" validate:"required,min=10,max=20"`
	Password             string  `json:"password" validate:"required,min=8,max=100"`
	PasswordConfirmation string  `json:"password_confirmation" validate:"required,eqfield=Password"`
}

func (dto *RegisterUserDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"Name.required":                 "Name is required",
		"Name.min":                      "Name must be at least 3 characters",
		"Name.max":                      "Name cannot exceed 100 characters",
		"MerchantUsername.min":          "Merchant username must be at least 3 characters",
		"MerchantUsername.max":          "Merchant username cannot exceed 100 characters",
		"Email.required":                "Email is required",
		"Email.email":                   "Email must be a valid email address",
		"Phone.required":                "Phone number is required",
		"Phone.min":                     "Phone number must be at least 10 characters",
		"Phone.max":                     "Phone number cannot exceed 20 characters",
		"Password.required":             "Password is required",
		"Password.min":                  "Password must be at least 8 characters",
		"Password.max":                  "Password cannot exceed 100 characters",
		"PasswordConfirmation.required": "Password confirmation is required",
		"PasswordConfirmation.eqfield":  "Password confirmation must match the password",
		"PasswordConfirmation.min":      "Password confirmation must be at least 8 characters",
	}
}
