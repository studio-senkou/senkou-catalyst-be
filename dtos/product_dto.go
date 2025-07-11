package dtos

type CreateProductDTO struct {
	Title        string  `json:"title" validate:"required,max=150"`
	Description  string  `json:"description" validate:"omitempty,max=500"`
	Price        float64 `json:"price" validate:"required,number,decimal=2"`
	AffiliateURL string  `json:"affiliate_url" validate:"required,url"`
	CategoryID   uint32  `json:"category_id" validate:"required,number"`
}

func (dto *CreateProductDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"title.required":         "Title is required",
		"title.max":              "Title must be at most 150 characters long",
		"price.required":         "Price is required",
		"price.number":           "Price must be a valid number",
		"price.decimal":          "Price must have at most 2 decimal places",
		"affiliate_url.required": "Affiliate URL is required",
		"affiliate_url.url":      "Affiliate URL must be a valid URL",
		"category_id.required":   "Category ID is required",
		"category_id.number":     "Category ID must be a valid number",
	}
}
