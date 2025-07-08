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

// Store a new category
// This function requires a category model to be passed in, which contains the detail of the category that will be stored.
// It returns the stored category or an error if the operation fails.
func (categoryRepo *categoryRepository) StoreCategory(category *models.Category) (*models.Category, error) {
	if err := categoryRepo.db.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// Finding a category by its name and merchant ID
// This function requires the name of the category and the merchant ID to be passed in.
// It returns the category if found or an error if the operation fails.
func (categoryRepo *categoryRepository) FindCategoryByName(name string, merchantID string) (*models.Category, error) {
	var category models.Category
	if err := categoryRepo.db.Where("name = ? AND merchant_id = ?", name, merchantID).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Finding a category by its ID
// Requires the ID of the category to be passed in.
// It returns the category if found or an error if the operation fails.
func (categoryRepo *categoryRepository) FindCategoryByID(id string) (*models.Category, error) {
	var category models.Category
	if err := categoryRepo.db.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Finding all categories by merchant ID
// This function requires the merchant ID to be passed in.
// It returns a slice of categories associated with the merchant or an error if the operation fails.
func (categoryRepo *categoryRepository) FindAllCategoriesByMerchantID(merchantID string) ([]models.Category, error) {
	var categories []models.Category
	if err := categoryRepo.db.Where("merchant_id = ?", merchantID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Update an existing category with the provided details
// This function requires an updated category model to be passed in.
// It returns the updated category or an error if the operation fails.
func (categoryRepo *categoryRepository) UpdateCategory(category *models.Category) (*models.Category, error) {
	if err := categoryRepo.db.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// Delete a category by its ID
// This function requires the ID of the category to be passed in.
// It returns an error if the operation fails.
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