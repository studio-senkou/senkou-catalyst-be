package services

import (
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"

	"github.com/google/uuid"
)

type MerchantService interface {
	CreateMerchant(merchant *dtos.CreateMerchantRequestDTO, userID uint) (*models.Merchant, error)
	GetUserMerchants(userID uint) (*[]models.Merchant, error)
}

type merchantService struct {
	merchantRepository repositories.MerchantRepository
}

func NewMerchantService(merchantRepository repositories.MerchantRepository) MerchantService {
	return &merchantService{
		merchantRepository: merchantRepository,
	}
}

func (s *merchantService) CreateMerchant(merchant *dtos.CreateMerchantRequestDTO, userID uint) (*models.Merchant, error) {
	createdMerchant, err := s.merchantRepository.Create(&models.Merchant{
		ID:      uuid.New().String(),
		Name:    merchant.Name,
		OwnerID: uint32(userID),
	})

	if err != nil {
		return nil, err
	}

	return createdMerchant, nil
}

func (s *merchantService) GetUserMerchants(userID uint) (*[]models.Merchant, error) {
	merchants, err := s.merchantRepository.FindMerchantsByUserID(userID)

	if err != nil {
		return nil, nil
	}

	return merchants, nil
}
