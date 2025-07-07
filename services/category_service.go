package services

import (
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"
)

type CategoryService interface {
	CreateNewCategory(category *dtos.CreateCategoryDTO, merchantID string) (*models.Category, error)
	GetCategoryByName(name string, merchantID string) (*models.Category, error)
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
