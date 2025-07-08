package repositories

import (
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	StoreCategory(category *models.Category) (*models.Category, error)
	FindCategoryByName(name string, merchantID string) (*models.Category, error)
	FindCategoryByID(id string) (*models.Category, error)
	FindAllCategoriesByMerchantID(merchantID string) ([]models.Category, error)
	UpdateCategory(category *models.Category) (*models.Category, error)
	DeleteCategory(id string) error
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

func (categoryRepo *categoryRepository) FindCategoryByID(id string) (*models.Category, error) {
	var category models.Category
	if err := categoryRepo.db.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (categoryRepo *categoryRepository) FindAllCategoriesByMerchantID(merchantID string) ([]models.Category, error) {
	var categories []models.Category
	if err := categoryRepo.db.Where("merchant_id = ?", merchantID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (categoryRepo *categoryRepository) UpdateCategory(category *models.Category) (*models.Category, error) {
	if err := categoryRepo.db.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (categoryRepo *categoryRepository) DeleteCategory(id string) error {
	var category models.Category
	if err := categoryRepo.db.Where("id = ?", id).First(&category).Error; err != nil {
		return err
	}

	if err := categoryRepo.db.Delete(&category).Error; err != nil {
		return err
	}

	return nil
}