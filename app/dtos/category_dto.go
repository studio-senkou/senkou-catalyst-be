package dtos

type CreateCategoryDTO struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

func (dto *CreateCategoryDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"Name.required": "Name is required",
		"Name.min":      "Name must be at least 3 characters long",
		"Name.max":      "Name must not exceed 100 characters",
	}
}

type UpdateCategoryDTO struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

func (dto *UpdateCategoryDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"Name.required": "Name is required",
		"Name.min":      "Name must be at least 3 characters long",
		"Name.max":      "Name must not exceed 100 characters",
	}
}
