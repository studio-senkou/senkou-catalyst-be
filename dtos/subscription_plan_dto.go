package dtos

type CreateSubscriptionPlanDTO struct {
	Name  string `json:"name" validate:"required,max=100"`
	Value string `json:"value" validate:"required"`
}

func (dto *CreateSubscriptionPlanDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"name.required":  "Name is required",
		"name.max":       "Name must be at most 100 characters long",
		"value.required": "Value is required",
	}
}
