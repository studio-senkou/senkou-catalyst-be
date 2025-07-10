package services

import (
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"
)

type CategoryService interface {
	CreateNewCategory(category *dtos.CreateCategoryDTO, merchantID string) (*models.Category, error)
	GetCategoryByName(name string, merchantID string) (*models.Category, error)
	GetCategoryByID(id string) (*models.Category, error)
	GetAllCategoriesByMerchantID(merchantID string) ([]models.Category, error)
	UpdateCategory(category *models.Category) (*models.Category, error)
	DeleteCategory(id uint32) error
}

type categoryService struct {
	categoryRepository repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepo,
	}
}

// Create a new category
// This function begins by creating a new category model from the provided data.
// It then attempts to store this category in the database using the category repository.
func (s *categoryService) CreateNewCategory(category *dtos.CreateCategoryDTO, merchantID string) (*models.Category, error) {
	categoryModel := &models.Category{
		Name:       category.Name,
		MerchantID: merchantID,
	}

	newCategory, err := s.categoryRepository.StoreCategory(categoryModel)

	if err != nil {
		return nil, err
	}

	return newCategory, nil
}

// Getting a category by its name
// This will search for a category by its name and the merchant ID.
// It returns the category if found or an error if the operation fails.
func (s *categoryService) GetCategoryByName(name string, merchantID string) (*models.Category, error) {
	category, err := s.categoryRepository.FindCategoryByName(name, merchantID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// Getting a category by its ID
// This function retrieves a category by its ID.
// It returns the category if found or an error if the operation fails.
func (s *categoryService) GetCategoryByID(id string) (*models.Category, error) {
	category, err := s.categoryRepository.FindCategoryByID(id)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// Getting all categories by its merchant
// This function will retrieve all categories associated with a specific merchant ID.
// It returns a slice of categories or an error if the operation fails.
func (s *categoryService) GetAllCategoriesByMerchantID(merchantID string) ([]models.Category, error) {
	categories, err := s.categoryRepository.FindAllCategoriesByMerchantID(merchantID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Update an existing category
// This function updates a category with the provided details.
// It requires an updated category model to be passed in.
// It returns the updated category or an error if the operation fails.
func (s *categoryService) UpdateCategory(category *models.Category) (*models.Category, error) {
	updatedCategory, err := s.categoryRepository.UpdateCategory(category)
	if err != nil {
		return nil, err
	}
	return updatedCategory, nil
}

// Delete a category by its ID
// This function deletes a category by its ID.
// It requires the ID of the category to be passed in.
// It returns an error if the operation fails.
func (s *categoryService) DeleteCategory(id uint32) error {
	err := s.categoryRepository.DeleteCategory(id)
	if err != nil {
		return err
	}
	return nil
}
