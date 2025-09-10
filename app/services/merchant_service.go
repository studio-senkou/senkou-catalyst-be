package services

import (
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"

	stderrors "errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MerchantService interface {
	CreateMerchant(merchant *dtos.CreateMerchantRequestDTO, userID uint32) (*models.Merchant, *errors.CustomError)
	GetMerchantByID(merchantID string) (*models.Merchant, *errors.CustomError)
	GetUserMerchants(userID uint32) ([]*models.Merchant, *errors.CustomError)
	GetMerchantOverview(merchantID string) (*dtos.MerchantOverview, *errors.CustomError)
	GetMerchantByUsername(username string) (*models.Merchant, *errors.CustomError)
	UpdateMerchantByID(merchantID string, updateData *dtos.UpdateMerchantRequestDTO) (*models.Merchant, *errors.CustomError)
	DeleteMerchantByID(merchantID string) *errors.CustomError
}

type MerchantServiceInstance struct {
	MerchantRepository repositories.MerchantRepository
	ProductRepository  repositories.ProductRepository
	CategoryRepository repositories.CategoryRepository
}

func NewMerchantService(merchantRepository repositories.MerchantRepository, productRepository repositories.ProductRepository, categoryRepository repositories.CategoryRepository) MerchantService {
	return &MerchantServiceInstance{
		MerchantRepository: merchantRepository,
		ProductRepository:  productRepository,
		CategoryRepository: categoryRepository,
	}
}

// Create a new merchant
// This function creates a new merchant for the user
// It returns the created merchant or an error if the creation fails
func (s *MerchantServiceInstance) CreateMerchant(merchant *dtos.CreateMerchantRequestDTO, userID uint32) (*models.Merchant, *errors.CustomError) {
	createdMerchant, err := s.MerchantRepository.Create(&models.Merchant{
		ID:      uuid.New().String(),
		Name:    merchant.Name,
		OwnerID: uint32(userID),
	})

	if err != nil {
		return nil, errors.Internal("Failed to create merchant", err.Error())
	}

	return createdMerchant, nil
}

// Get user merchants
// This function retrieves all merchants associated with a user
// It returns a slice of merchants or an error if the retrieval fails
func (s *MerchantServiceInstance) GetUserMerchants(userID uint32) ([]*models.Merchant, *errors.CustomError) {
	merchants, err := s.MerchantRepository.FindByUserID(userID)

	if err != nil {
		return nil, errors.Internal("Failed to retrieve merchants", err.Error())
	}

	return merchants, nil
}

// Get merchant by ID
// This function retrieves a merchant by its ID
// It returns the merchant or an error if the retrieval fails
func (s *MerchantServiceInstance) GetMerchantByID(merchantID string) (*models.Merchant, *errors.CustomError) {
	merchant, err := s.MerchantRepository.FindByID(merchantID)

	if err != nil {
		return nil, errors.Internal("Failed to retrieve merchant", err.Error())
	}

	return merchant, nil
}

// Get merchant's overview
// This function retrieves an overview of a specific merchant by its ID
// It returns the merchant overview or an error if the retrieval fails
func (s *MerchantServiceInstance) GetMerchantOverview(merchantID string) (*dtos.MerchantOverview, *errors.CustomError) {
	overview, err := s.MerchantRepository.FindOverview(merchantID)

	if err != nil {
		return nil, errors.Internal("Failed to retrieve merchant overview", err.Error())
	}

	return overview, nil
}

// Get merchant by username
// This function retrieves a merchant by its username
// It returns the merchant or an error if the retrieval fails
func (s *MerchantServiceInstance) GetMerchantByUsername(username string) (*models.Merchant, *errors.CustomError) {
	merchant, err := s.MerchantRepository.FindByUsername(username)

	if err != nil {

		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("Merchant not found")
		}

		return nil, errors.Internal("Failed to retrieve merchant by username", err.Error())
	}

	return merchant, nil
}

// Update merchant by ID
// This function updates an existing merchant with the provided data
// It returns the updated merchant or an error if the update fails
func (s *MerchantServiceInstance) UpdateMerchantByID(merchantID string, updateData *dtos.UpdateMerchantRequestDTO) (*models.Merchant, *errors.CustomError) {
	updatedMerchant, err := s.MerchantRepository.UpdateMerchant(merchantID, &models.Merchant{
		Name: updateData.Name,
	})

	if err != nil {
		return nil, errors.Internal("Failed to update merchant", err.Error())
	}

	return updatedMerchant, nil
}

// Delete merchant by ID
// This function deletes a merchant by its ID
// It returns an error if the deletion fails
func (s *MerchantServiceInstance) DeleteMerchantByID(merchantID string) *errors.CustomError {
	err := s.MerchantRepository.DeleteMerchant(merchantID)

	if err != nil {
		return errors.Internal("Failed to delete merchant", err.Error())
	}

	return nil
}
