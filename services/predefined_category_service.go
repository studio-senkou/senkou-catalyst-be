package services

import (
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/errors"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"

	"gorm.io/gorm"
)

type PredefinedCategoryService interface {
	StoreCategory(pdCategory *dtos.CreatePDCategoryDTO) (*models.PredefinedCategory, *errors.AppError)
	GetPredefinedCategoryByName(name string) (*models.PredefinedCategory, *errors.AppError)
	GetAllPredefinedCategories() ([]*models.PredefinedCategory, *errors.AppError)
	UpdatePredefinedCategory(pdCategory *dtos.UpdatePDCategoryDTO, pdCategoryID uint32) *errors.AppError
	DeletePredefinedCategory(pdCategoryID uint32) *errors.AppError
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
func (s *PredefinedCategoryServiceInstance) StoreCategory(pdCategory *dtos.CreatePDCategoryDTO) (*models.PredefinedCategory, *errors.AppError) {
	if existingCategory, err := s.PredefinedCategoryRepository.FindByName(pdCategory.Name); err == nil && existingCategory != nil {
		return nil, errors.NewAppError(400, "Predefined category already exists")
	}

	predefinedCategory := &models.PredefinedCategory{
		Name:        pdCategory.Name,
		Description: pdCategory.Description,
		ImageURL:    pdCategory.ImageURL,
	}

	if err := s.PredefinedCategoryRepository.StoreCategory(predefinedCategory); err != nil {
		return nil, errors.NewAppError(500, "Failed to create predefined category")
	}

	return predefinedCategory, nil
}

// Get a predefined category by its name
// This function will return an error if the category is not found
// It returns the category if found
func (s *PredefinedCategoryServiceInstance) GetPredefinedCategoryByName(name string) (*models.PredefinedCategory, *errors.AppError) {
	category, err := s.PredefinedCategoryRepository.FindByName(name)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to retrieve predefined category")
	}

	if category == nil {
		return nil, errors.NewAppError(404, "Predefined category not found")
	}

	return category, nil
}

// Get all predefined categories
// This function retrieves all predefined categories from the repository
// It returns a slice of predefined categories or an error if the operation fails
func (s *PredefinedCategoryServiceInstance) GetAllPredefinedCategories() ([]*models.PredefinedCategory, *errors.AppError) {
	categories, err := s.PredefinedCategoryRepository.FindAll()
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to retrieve predefined categories")
	}

	if len(categories) == 0 {
		return nil, errors.NewAppError(404, "No predefined categories found")
	}

	return categories, nil
}

// Update predefined category by its ID
// This function requires the ID of the category and the updated data to be passed in
// It returns an error if the operation fails
func (s *PredefinedCategoryServiceInstance) UpdatePredefinedCategory(pdCategory *dtos.UpdatePDCategoryDTO, pdCategoryID uint32) *errors.AppError {
	if pdCategoryID == 0 {
		return errors.NewAppError(400, "Invalid category ID")
	}

	updatedCategory := &models.PredefinedCategory{
		ID:          pdCategoryID,
		Name:        pdCategory.Name,
		Description: pdCategory.Description,
		ImageURL:    pdCategory.ImageURL,
	}

	if err := s.PredefinedCategoryRepository.UpdateByID(pdCategoryID, updatedCategory); err != nil {
		return errors.NewAppError(500, "Failed to update predefined category")
	}

	return nil
}

// Delete predefined category by its ID
// This function requires the ID of the category to be passed in
// It returns an error if the operation fails
func (s *PredefinedCategoryServiceInstance) DeletePredefinedCategory(PDCategoryID uint32) *errors.AppError {
	if PDCategoryID == 0 {
		return errors.NewAppError(400, "Invalid category ID")
	}

	if err := s.PredefinedCategoryRepository.RemoveByID(PDCategoryID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError(404, "Predefined category not found")
		}

		return errors.NewAppError(500, "Failed to delete predefined category")
	}

	return nil
}
