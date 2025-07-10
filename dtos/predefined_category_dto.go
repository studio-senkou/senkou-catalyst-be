package dtos

type CreatePDCategoryDTO struct {
	Name        string `json:"name"        validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,max=255"`
	ImageURL    string `json:"image_url"   validate:"required,url"`
}

func (dto *CreatePDCategoryDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"name.required":      "Name is required",
		"name.min":           "Name must be at least 3 characters long",
		"name.max":           "Name must not exceed 100 characters",
		"description.max":    "Description must not exceed 255 characters",
		"image_url.required": "Image URL is required",
		"image_url.url":      "Image URL must be a valid URL",
	}
}
