package services

import (
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"

	"github.com/google/uuid"
)

type MerchantService interface {
	CreateMerchant(merchant *dtos.CreateMerchantRequestDTO, userID uint32) (*models.Merchant, *errors.AppError)
	GetUserMerchants(userID uint32) ([]*models.Merchant, *errors.AppError)
	GetMerchantByID(merchantID string) (*models.Merchant, *errors.AppError)
	UpdateMerchantByID(merchantID string, updateData *dtos.UpdateMerchantRequestDTO) (*models.Merchant, *errors.AppError)
	DeleteMerchantByID(merchantID string) *errors.AppError
}

type MerchantServiceInstance struct {
	MerchantRepository repositories.MerchantRepository
}

func NewMerchantService(merchantRepository repositories.MerchantRepository) MerchantService {
	return &MerchantServiceInstance{
		MerchantRepository: merchantRepository,
	}
}

// Create a new merchant
// This function creates a new merchant for the user
// It returns the created merchant or an error if the creation fails
func (s *MerchantServiceInstance) CreateMerchant(merchant *dtos.CreateMerchantRequestDTO, userID uint32) (*models.Merchant, *errors.AppError) {
	createdMerchant, err := s.MerchantRepository.Create(&models.Merchant{
		ID:      uuid.New().String(),
		Name:    merchant.Name,
		OwnerID: uint32(userID),
	})

	if err != nil {
		return nil, errors.NewAppError(500, "Failed to create merchant")
	}

	return createdMerchant, nil
}

// Get user merchants
// This function retrieves all merchants associated with a user
// It returns a slice of merchants or an error if the retrieval fails
func (s *MerchantServiceInstance) GetUserMerchants(userID uint32) ([]*models.Merchant, *errors.AppError) {
	merchants, err := s.MerchantRepository.FindMerchantsByUserID(userID)

	if err != nil {
		return nil, errors.NewAppError(500, "Failed to retrieve merchants")
	}

	return merchants, nil
}

// Get merchant by ID
// This function retrieves a merchant by its ID
// It returns the merchant or an error if the retrieval fails
func (s *MerchantServiceInstance) GetMerchantByID(merchantID string) (*models.Merchant, *errors.AppError) {
	merchant, err := s.MerchantRepository.FindByID(merchantID)

	if err != nil {
		return nil, errors.NewAppError(500, "Failed to retrieve merchant")
	}

	return merchant, nil
}

// Update merchant by ID
// This function updates an existing merchant with the provided data
// It returns the updated merchant or an error if the update fails
func (s *MerchantServiceInstance) UpdateMerchantByID(merchantID string, updateData *dtos.UpdateMerchantRequestDTO) (*models.Merchant, *errors.AppError) {
	updatedMerchant, err := s.MerchantRepository.UpdateMerchant(merchantID, &models.Merchant{
		Name: updateData.Name,
	})

	if err != nil {
		return nil, errors.NewAppError(500, "Failed to update merchant")
	}

	return updatedMerchant, nil
}

// Delete merchant by ID
// This function deletes a merchant by its ID
// It returns an error if the deletion fails
func (s *MerchantServiceInstance) DeleteMerchantByID(merchantID string) *errors.AppError {
	err := s.MerchantRepository.DeleteMerchant(merchantID)

	if err != nil {
		return errors.NewAppError(500, "Failed to delete merchant")
	}

	return nil
}
