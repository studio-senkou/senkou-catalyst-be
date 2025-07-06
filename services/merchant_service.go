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
	GetMerchantByID(merchantID string) (*models.Merchant, error)
	UpdateMerchantByID(merchantID string, updateData *dtos.UpdateMerchantRequestDTO) (*models.Merchant, error)
	DeleteMerchantByID(merchantID string) error
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

func (s *merchantService) GetMerchantByID(merchantID string) (*models.Merchant, error) {
	merchant, err := s.merchantRepository.FindByID(merchantID)

	if err != nil {
		return nil, err
	}

	return merchant, nil
}

func (s *merchantService) UpdateMerchantByID(merchantID string, updateData *dtos.UpdateMerchantRequestDTO) (*models.Merchant, error) {
	updatedMerchant, err := s.merchantRepository.UpdateMerchant(merchantID, &models.Merchant{
		Name: updateData.Name,
	})

	if err != nil {
		return nil, err
	}

	return updatedMerchant, nil
}

func (s *merchantService) DeleteMerchantByID(merchantID string) error {
	err := s.merchantRepository.DeleteMerchant(merchantID)

	if err != nil {
		return err
	}

	return nil
}
