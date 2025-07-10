package services

import (
	"errors"
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"
)

type PredefinedCategoryService interface {
	StoreCategory(pdCategory *dtos.CreatePDCategoryDTO) error
	GetPredefinedCategoryByName(name string) (*models.PredefinedCategory, error)
	GetAllPredefinedCategories() (*[]models.PredefinedCategory, error)
	UpdatePredefinedCategory(pdCategory *dtos.UpdatePDCategoryDTO, pdCategoryID uint32) error
	DeletePredefinedCategory(pdCategoryID uint32) error
}

type predefinedCategoryService struct {
	pcRepository repositories.PredefinedCategoryRepository
}

func NewPredefinedCategoryService(pcRepository repositories.PredefinedCategoryRepository) PredefinedCategoryService {
	return &predefinedCategoryService{
		pcRepository: pcRepository,
	}
}

// Create a new category and store into predefined categories repository
// This function will return an error if the category already exists
func (s *predefinedCategoryService) StoreCategory(pdCategory *dtos.CreatePDCategoryDTO) error {
	if existingCategory, err := s.pcRepository.FindByName(pdCategory.Name); err == nil && existingCategory != nil {
		return errors.New("category already exists")
	}

	predefinedCategory := &models.PredefinedCategory{
		Name:        pdCategory.Name,
		Description: pdCategory.Description,
		ImageURL:    pdCategory.ImageURL,
	}

	if err := s.pcRepository.StoreCategory(predefinedCategory); err != nil {
		return err
	}

	return nil
}

// Get a predefined category by its name
// This function will return an error if the category is not found
// It returns the category if found
func (s *predefinedCategoryService) GetPredefinedCategoryByName(name string) (*models.PredefinedCategory, error) {
	category, err := s.pcRepository.FindByName(name)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("predefined category not found")
	}

	return category, nil
}

// Get all predefined categories
// This function retrieves all predefined categories from the repository
// It returns a slice of predefined categories or an error if the operation fails
func (s *predefinedCategoryService) GetAllPredefinedCategories() (*[]models.PredefinedCategory, error) {
	categories, err := s.pcRepository.FindAll()
	if err != nil {
		return nil, err
	}

	if categories == nil || len(*categories) == 0 {
		return nil, errors.New("no predefined categories found")
	}

	return categories, nil
}

// Update predefined category by its ID
// This function requires the ID of the category and the updated data to be passed in
// It returns an error if the operation fails
func (s *predefinedCategoryService) UpdatePredefinedCategory(pdCategory *dtos.UpdatePDCategoryDTO, pdCategoryID uint32) error {
	if pdCategoryID == 0 {
		return errors.New("invalid category ID")
	}

	updatedCategory := &models.PredefinedCategory{
		ID: 		 pdCategoryID,
		Name:        pdCategory.Name,
		Description: pdCategory.Description,
		ImageURL:    pdCategory.ImageURL,
	}

	if err := s.pcRepository.UpdateByID(pdCategoryID, updatedCategory); err != nil {
		return err
	}

	return nil
}

// Delete predefined category by its ID
// This function requires the ID of the category to be passed in
// It returns an error if the operation fails
func (s *predefinedCategoryService) DeletePredefinedCategory(pdCategoryID uint32) error {
	if pdCategoryID == 0 {
		return errors.New("invalid category ID")
	}

	if err := s.pcRepository.RemoveByID(pdCategoryID); err != nil {
		return err
	}
	
	return nil
}