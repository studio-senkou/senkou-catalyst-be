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
	DeleteCategory(id string) error
}

type categoryService struct {
	categoryRepository repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepo,
	}
}

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

func (s *categoryService) GetCategoryByName(name string, merchantID string) (*models.Category, error) {
	category, err := s.categoryRepository.FindCategoryByName(name, merchantID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) GetCategoryByID(id string) (*models.Category, error) {
	category, err := s.categoryRepository.FindCategoryByID(id)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) GetAllCategoriesByMerchantID(merchantID string) ([]models.Category, error) {
	categories, err := s.categoryRepository.FindAllCategoriesByMerchantID(merchantID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *categoryService) UpdateCategory(category *models.Category) (*models.Category, error) {
	updatedCategory, err := s.categoryRepository.UpdateCategory(category)
	if err != nil {
		return nil, err
	}
	return updatedCategory, nil
}

func (s *categoryService) DeleteCategory(id string) error {
	err := s.categoryRepository.DeleteCategory(id)
	if err != nil {
		return err
	}
	return nil
}