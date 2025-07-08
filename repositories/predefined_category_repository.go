package repositories

import (
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

type PredefinedCategoryRepository interface {
	StoreCategory(pdCategory *models.PredefinedCategory) error
	FindByName(name string) (*models.PredefinedCategory, error)
	FindAll() (*[]models.PredefinedCategory, error)
	UpdateByID(pcID uint, updatedName string) error
	RemoveByID(pcID uint) error
}

type predefinedCategoryRepository struct {
	DB *gorm.DB
}

func NewPredefinedCategoryRepository(db *gorm.DB) PredefinedCategoryRepository {
	return &predefinedCategoryRepository{
		DB: db,
	}
}

func (r *predefinedCategoryRepository) StoreCategory(pdCategory *models.PredefinedCategory) error {
	return r.DB.Create(pdCategory).Error
}

func (r *predefinedCategoryRepository) FindByName(name string) (*models.PredefinedCategory, error) {
	var category models.PredefinedCategory

	err := r.DB.Where("name = ?", name).First(&category).Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *predefinedCategoryRepository) FindAll() (*[]models.PredefinedCategory, error) {
	var categories []models.PredefinedCategory

	err := r.DB.Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return &categories, nil
}

func (r *predefinedCategoryRepository) UpdateByID(pcID uint, updatedName string) error {
	var category models.PredefinedCategory

	err := r.DB.Where("id = ?", pcID).First(&category).Error
	
	if err != nil {
		return err
	}

	category.Name = updatedName
	
	return r.DB.Save(&category).Error
}

func (r *predefinedCategoryRepository) RemoveByID(pcID uint) error {
	var category models.PredefinedCategory

	err := r.DB.Where("id = ?", pcID).First(&category).Error
	if err != nil {
		return err
	}

	return r.DB.Delete(&category).Error
}