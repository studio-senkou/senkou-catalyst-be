package services

import (
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"

	"github.com/google/uuid"
)

type ProductInteractionService interface {
	StoreLog(productID uuid.UUID, request *dtos.SendProductInteractionDTO) *errors.AppError
}

type ProductInteractionServiceInstance struct {
	PIRepo repositories.ProductInteractionRepository
}

func NewProductInteractionService(piRepo repositories.ProductInteractionRepository) ProductInteractionService {
	return &ProductInteractionServiceInstance{
		PIRepo: piRepo,
	}
}

func (s *ProductInteractionServiceInstance) StoreLog(productID uuid.UUID, request *dtos.SendProductInteractionDTO) *errors.AppError {
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
		return errors.NewAppError(500, "failed to store product interaction log")
	}

	return nil
}
