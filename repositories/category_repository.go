package repositories

import (
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	StoreCategory(category *models.Category) (*models.Category, error)
	FindCategoryByName(name string, merchantID string) (*models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (categoryRepo *categoryRepository) StoreCategory(category *models.Category) (*models.Category, error) {
	if err := categoryRepo.db.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (categoryRepo *categoryRepository) FindCategoryByName(name string, merchantID string) (*models.Category, error) {
	var category models.Category
	if err := categoryRepo.db.Where("name = ? AND merchant_id = ?", name, merchantID).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}
