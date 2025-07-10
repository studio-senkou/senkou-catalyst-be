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
}

type predefinedCategoryService struct {
	pcRepository repositories.PredefinedCategoryRepository
}

func NewPredefinedCategoryService(pcRepository repositories.PredefinedCategoryRepository) PredefinedCategoryService {
	return &predefinedCategoryService{
		pcRepository: pcRepository,
	}
}

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
