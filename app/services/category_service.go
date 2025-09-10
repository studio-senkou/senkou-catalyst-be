package services

import (
	stderror "errors"

	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"

	"gorm.io/gorm"
)

type CategoryService interface {
	CreateNewCategory(category *dtos.CreateCategoryDTO, merchantID string) (*models.Category, *errors.CustomError)
	GetCategoryByName(name string, merchantID string) (*models.Category, *errors.CustomError)
	GetCategoryByID(id string) (*models.Category, *errors.CustomError)
	GetAllCategoriesByMerchantID(merchantID string) ([]*models.Category, *errors.CustomError)
	UpdateCategory(category *models.Category) (*models.Category, *errors.CustomError)
	DeleteCategory(id uint32) *errors.CustomError
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
func (s *CategoryServiceInstance) CreateNewCategory(category *dtos.CreateCategoryDTO, merchantID string) (*models.Category, *errors.CustomError) {
	categoryModel := &models.Category{
		Name:       category.Name,
		MerchantID: merchantID,
	}

	newCategory, err := s.CategoryRepository.StoreCategory(categoryModel)

	if err != nil {
		return nil, errors.Internal("Failed to create category", err.Error())
	}

	return newCategory, nil
}

// Getting a category by its name
// This will search for a category by its name and the merchant ID.
// It returns the category if found or an error if the operation fails.
func (s *CategoryServiceInstance) GetCategoryByName(name string, merchantID string) (*models.Category, *errors.CustomError) {
	category, err := s.CategoryRepository.FindCategoryByName(name, merchantID)
	if err != nil {
		if stderror.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("Category not found")
		}

		return nil, errors.Internal("Failed to get category by name", err.Error())
	}

	return category, nil
}

// Getting a category by its ID
// This function retrieves a category by its ID.
// It returns the category if found or an error if the operation fails.
func (s *CategoryServiceInstance) GetCategoryByID(id string) (*models.Category, *errors.CustomError) {
	category, err := s.CategoryRepository.FindCategoryByID(id)

	if err != nil {
		if stderror.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("Category not found")
		}

		return nil, errors.Internal("Failed to get category by ID", err.Error())
	}
	return category, nil
}

// Getting all categories by its merchant
// This function will retrieve all categories associated with a specific merchant ID.
// It returns a slice of categories or an error if the operation fails.
func (s *CategoryServiceInstance) GetAllCategoriesByMerchantID(merchantID string) ([]*models.Category, *errors.CustomError) {
	categories, err := s.CategoryRepository.FindAllCategoriesByMerchantID(merchantID)
	if err != nil {
		return nil, errors.Internal("Failed to get categories by merchant ID", err.Error())
	}
	return categories, nil
}

// Update an existing category
// This function updates a category with the provided details.
// It requires an updated category model to be passed in.
// It returns the updated category or an error if the operation fails.
func (s *CategoryServiceInstance) UpdateCategory(category *models.Category) (*models.Category, *errors.CustomError) {
	updatedCategory, err := s.CategoryRepository.UpdateCategory(category)
	if err != nil {
		return nil, errors.Internal("Failed to update category", err.Error())
	}
	return updatedCategory, nil
}

// Delete a category by its ID
// This function deletes a category by its ID.
// It requires the ID of the category to be passed in.
// It returns an error if the operation fails.
func (s *CategoryServiceInstance) DeleteCategory(id uint32) *errors.CustomError {
	err := s.CategoryRepository.DeleteCategory(id)
	if err != nil {
		return errors.Internal("Failed to delete category", err.Error())
	}
	return nil
}
