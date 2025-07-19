package services

import (
	"fmt"
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/errors"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"
)

type CategoryService interface {
	CreateNewCategory(category *dtos.CreateCategoryDTO, merchantID string) (*models.Category, *errors.AppError)
	GetCategoryByName(name string, merchantID string) (*models.Category, *errors.AppError)
	GetCategoryByID(id string) (*models.Category, *errors.AppError)
	GetAllCategoriesByMerchantID(merchantID string) ([]*models.Category, *errors.AppError)
	UpdateCategory(category *models.Category) (*models.Category, *errors.AppError)
	DeleteCategory(id uint32) *errors.AppError
}

type CategoryServiceInstance struct {
	CategoryRepository repositories.CategoryRepository
}

func NewCategoryService(categoryRepository repositories.CategoryRepository) CategoryService {
	return &CategoryServiceInstance{
		CategoryRepository: categoryRepository,
	}
}

// Create a new category
// This function begins by creating a new category model from the provided data.
// It then attempts to store this category in the database using the category repository.
func (s *CategoryServiceInstance) CreateNewCategory(category *dtos.CreateCategoryDTO, merchantID string) (*models.Category, *errors.AppError) {
	categoryModel := &models.Category{
		Name:       category.Name,
		MerchantID: merchantID,
	}

	newCategory, err := s.CategoryRepository.StoreCategory(categoryModel)

	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to create category: %v", err.Error()))
	}

	return newCategory, nil
}

// Getting a category by its name
// This will search for a category by its name and the merchant ID.
// It returns the category if found or an error if the operation fails.
func (s *CategoryServiceInstance) GetCategoryByName(name string, merchantID string) (*models.Category, *errors.AppError) {
	category, err := s.CategoryRepository.FindCategoryByName(name, merchantID)
	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to get category by name: %v", err.Error()))
	}
	return category, nil
}

// Getting a category by its ID
// This function retrieves a category by its ID.
// It returns the category if found or an error if the operation fails.
func (s *CategoryServiceInstance) GetCategoryByID(id string) (*models.Category, *errors.AppError) {
	category, err := s.CategoryRepository.FindCategoryByID(id)
	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to get category by ID: %v", err.Error()))
	}
	return category, nil
}

// Getting all categories by its merchant
// This function will retrieve all categories associated with a specific merchant ID.
// It returns a slice of categories or an error if the operation fails.
func (s *CategoryServiceInstance) GetAllCategoriesByMerchantID(merchantID string) ([]*models.Category, *errors.AppError) {
	categories, err := s.CategoryRepository.FindAllCategoriesByMerchantID(merchantID)
	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to get categories by merchant ID: %v", err.Error()))
	}
	return categories, nil
}

// Update an existing category
// This function updates a category with the provided details.
// It requires an updated category model to be passed in.
// It returns the updated category or an error if the operation fails.
func (s *CategoryServiceInstance) UpdateCategory(category *models.Category) (*models.Category, *errors.AppError) {
	updatedCategory, err := s.CategoryRepository.UpdateCategory(category)
	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to update category: %v", err.Error()))
	}
	return updatedCategory, nil
}

// Delete a category by its ID
// This function deletes a category by its ID.
// It requires the ID of the category to be passed in.
// It returns an error if the operation fails.
func (s *CategoryServiceInstance) DeleteCategory(id uint32) *errors.AppError {
	err := s.CategoryRepository.DeleteCategory(id)
	if err != nil {
		return errors.NewAppError(500, fmt.Sprintf("Failed to delete category: %v", err.Error()))
	}
	return nil
}
