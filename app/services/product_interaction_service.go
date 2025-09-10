package services

import (
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"

	"senkou-catalyst-be/utils/query"

	"github.com/google/uuid"
)

type ProductInteractionService interface {
	StoreLog(productID uuid.UUID, request *dtos.SendProductInteractionDTO) *errors.CustomError
	GetProductMetrics(merchantID string, params *query.QueryParams) (*dtos.OverallProductMetrics, *errors.CustomError)
}

type ProductInteractionServiceInstance struct {
	PIRepo repositories.ProductInteractionRepository
}

func NewProductInteractionService(piRepo repositories.ProductInteractionRepository) ProductInteractionService {
	return &ProductInteractionServiceInstance{
		PIRepo: piRepo,
	}
}

func (s *ProductInteractionServiceInstance) StoreLog(productID uuid.UUID, request *dtos.SendProductInteractionDTO) *errors.CustomError {
	userAgent := models.UserAgent{
		Browser: request.Browser,
		OS:      request.OS,
	}

	productInteraction := &models.ProductMetric{
		ProductID:   productID,
		Origin:      request.Origin,
		UserAgent:   userAgent,
		Interaction: request.Interaction,
	}

	if err := s.PIRepo.StoreProductInteractionLog(productInteraction); err != nil {
		return errors.Internal("Failed to store product interaction log", err.Error())
	}

	return nil
}

func (s *ProductInteractionServiceInstance) GetProductMetrics(merchantID string, params *query.QueryParams) (*dtos.OverallProductMetrics, *errors.CustomError) {
	merchantProductsStat, err := s.PIRepo.GetMerchantProductsMetric(merchantID, params)
	if err != nil {
		return nil, errors.Internal("Failed to get product metrics", err.Error())
	}

	overallStats, err := s.PIRepo.GetMerchantProductsMetricStats(merchantID, params)
	if err != nil {
		return nil, errors.Internal("Failed to get merchant products stats", err.Error())
	}

	result := &dtos.OverallProductMetrics{
		OverallStats: overallStats,
		ProductsStat: merchantProductsStat,
	}

	return result, nil
}
