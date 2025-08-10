package dtos

type CreateProductDTO struct {
	Title        string  `json:"title" validate:"required,max=150"`
	Description  string  `json:"description" validate:"omitempty,max=500"`
	Price        float64 `json:"price" validate:"required,number,min=0"`
	AffiliateURL string  `json:"affiliate_url" validate:"required,url"`
	Photo        string  `json:"photo" validate:"omitempty"`
	CategoryID   *uint32 `json:"category_id" validate:"omitempty,number"`
}

func (dto *CreateProductDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"title.required":         "Title is required",
		"title.max":              "Title must be at most 150 characters long",
		"description.max":        "Description must be at most 500 characters long",
		"price.required":         "Price is required",
		"price.number":           "Price must be a valid number",
		"price.min":              "Price must be greater than or equal to 0",
		"affiliate_url.required": "Affiliate URL is required",
		"affiliate_url.url":      "Affiliate URL must be a valid URL",
		"category_id.number":     "Category ID must be a valid number",
	}
}

type UpdateProductDTO struct {
	ID           string   `json:"id"`
	Title        *string  `json:"title" validate:"omitempty,max=150"`
	Description  *string  `json:"description" validate:"omitempty,max=500"`
	Price        *float64 `json:"price" validate:"omitempty,number,min=0"`
	AffiliateURL *string  `json:"affiliate_url" validate:"omitempty,url"`
	CategoryID   *uint32  `json:"category_id" validate:"omitempty,number"`
}

func (dto *UpdateProductDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"title.max":          "Title must be at most 150 characters long",
		"description.max":    "Description must be at most 500 characters long",
		"price.number":       "Price must be a valid number",
		"price.min":          "Price must be greater than or equal to 0",
		"affiliate_url.url":  "Affiliate URL must be a valid URL",
		"category_id.number": "Category ID must be a valid number",
	}
}
