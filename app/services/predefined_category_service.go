package services

import (
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"

	"gorm.io/gorm"
)

type PredefinedCategoryService interface {
	StoreCategory(pdCategory *dtos.CreatePDCategoryDTO) (*models.PredefinedCategory, *errors.CustomError)
	GetPredefinedCategoryByName(name string) (*models.PredefinedCategory, *errors.CustomError)
	GetAllPredefinedCategories() ([]*models.PredefinedCategory, *errors.CustomError)
	UpdatePredefinedCategory(pdCategory *dtos.UpdatePDCategoryDTO, pdCategoryID uint32) *errors.CustomError
	DeletePredefinedCategory(pdCategoryID uint32) *errors.CustomError
}

type PredefinedCategoryServiceInstance struct {
	PredefinedCategoryRepository repositories.PredefinedCategoryRepository
}

func NewPredefinedCategoryService(PCRepository repositories.PredefinedCategoryRepository) PredefinedCategoryService {
	return &PredefinedCategoryServiceInstance{
		PredefinedCategoryRepository: PCRepository,
	}
}

// Create a new category and store into predefined categories repository
// This function will return an error if the category already exists
func (s *PredefinedCategoryServiceInstance) StoreCategory(pdCategory *dtos.CreatePDCategoryDTO) (*models.PredefinedCategory, *errors.CustomError) {
	if existingCategory, err := s.PredefinedCategoryRepository.FindByName(pdCategory.Name); err == nil && existingCategory != nil {
		return nil, errors.Conflict("Predefined category already exists", nil)
	}

	predefinedCategory := &models.PredefinedCategory{
		Name:        pdCategory.Name,
		Description: pdCategory.Description,
		ImageURL:    pdCategory.ImageURL,
	}

	if err := s.PredefinedCategoryRepository.StoreCategory(predefinedCategory); err != nil {
		return nil, errors.Internal("Failed to create predefined category", err.Error())
	}

	return predefinedCategory, nil
}

// Get a predefined category by its name
// This function will return an error if the category is not found
// It returns the category if found
func (s *PredefinedCategoryServiceInstance) GetPredefinedCategoryByName(name string) (*models.PredefinedCategory, *errors.CustomError) {
	category, err := s.PredefinedCategoryRepository.FindByName(name)
	if err != nil {
		return nil, errors.Internal("Failed to retrieve predefined category", err.Error())
	}

	if category == nil {
		return nil, errors.NotFound("Predefined category not found")
	}

	return category, nil
}

// Get all predefined categories
// This function retrieves all predefined categories from the repository
// It returns a slice of predefined categories or an error if the operation fails
func (s *PredefinedCategoryServiceInstance) GetAllPredefinedCategories() ([]*models.PredefinedCategory, *errors.CustomError) {
	categories, err := s.PredefinedCategoryRepository.FindAll()
	if err != nil {
		return nil, errors.Internal("Failed to retrieve predefined categories", err.Error())
	}

	if len(categories) == 0 {
		return nil, errors.NotFound("No predefined categories found")
	}

	return categories, nil
}

// Update predefined category by its ID
// This function requires the ID of the category and the updated data to be passed in
// It returns an error if the operation fails
func (s *PredefinedCategoryServiceInstance) UpdatePredefinedCategory(pdCategory *dtos.UpdatePDCategoryDTO, pdCategoryID uint32) *errors.CustomError {
	if pdCategoryID == 0 {
		return errors.BadRequest("Invalid category ID", nil)
	}

	updatedCategory := &models.PredefinedCategory{
		ID:          pdCategoryID,
		Name:        pdCategory.Name,
		Description: pdCategory.Description,
		ImageURL:    pdCategory.ImageURL,
	}

	if err := s.PredefinedCategoryRepository.UpdateByID(pdCategoryID, updatedCategory); err != nil {
		return errors.Internal("Failed to update predefined category", err.Error())
	}

	return nil
}

// Delete predefined category by its ID
// This function requires the ID of the category to be passed in
// It returns an error if the operation fails
func (s *PredefinedCategoryServiceInstance) DeletePredefinedCategory(PDCategoryID uint32) *errors.CustomError {
	if PDCategoryID == 0 {
		return errors.BadRequest("Invalid category ID", nil)
	}

	if err := s.PredefinedCategoryRepository.RemoveByID(PDCategoryID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NotFound("Predefined category not found")
		}

		return errors.Internal("Failed to delete predefined category", err.Error())
	}

	return nil
}
