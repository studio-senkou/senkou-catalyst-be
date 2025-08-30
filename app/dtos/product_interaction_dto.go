package dtos

import (
	"senkou-catalyst-be/app/models"
)

type SendProductInteractionDTO struct {
	Origin      string                          `json:"origin" validate:"required"`
	Browser     string                          `json:"browser" validate:"required"`
	OS          string                          `json:"os" validate:"required"`
	Interaction models.ProductMetricInteraction `json:"interaction_type" validate:"required"`
}

func (dto *SendProductInteractionDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"origin.required":           "Origin is required",
		"browser.required":          "Browser is required",
		"os.required":               "OS is required",
		"interaction_type.required": "Interaction type is required",
	}
}

type ProductMetricStats struct {
	TotalViews  int64 `json:"total_views"`
	TotalClicks int64 `json:"total_clicks"`
}

type ProductReport struct {
	Name        string `json:"product_name"`
	TotalViews  int64  `json:"total_views"`
	TotalClicks int64  `json:"total_clicks"`
}

type OverallProductMetrics struct {
	OverallStats *ProductMetricStats `json:"overall_stats"`
	ProductsStat []ProductReport     `json:"products_stat"`
}
