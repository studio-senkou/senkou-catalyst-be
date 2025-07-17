package dtos

type CreateSubscriptionDTO struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Price       float64 `json:"price" validate:"required,number,min=0"`
	Description string  `json:"description" validate:"omitempty,max=500"`
	Duration    int16   `json:"duration" validate:"required,number,min=1"`
}

func (dto *CreateSubscriptionDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"name.required":     "Name is required",
		"name.max":          "Name must be at most 100 characters long",
		"price.required":    "Price is required",
		"price.number":      "Price must be a valid number",
		"price.min":         "Price must be greater than or equal to 0",
		"description.max":   "Description must be at most 500 characters long",
		"duration.required": "Duration is required",
		"duration.number":   "Duration must be a valid number",
		"duration.min":      "Duration must be at least 1 month",
	}
}

type UpdateSubscriptionDTO struct {
	Name        *string  `json:"name" validate:"required,max=100"`
	Price       *float64 `json:"price" validate:"omitempty,number,min=0"`
	Description *string  `json:"description" validate:"omitempty,max=500"`
	Duration    *int16   `json:"duration" validate:"omitempty,number,min=1"`
}
