package repositories

import (
	"fmt"
	"senkou-catalyst-be/app/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	StoreCategory(category *models.Category) (*models.Category, error)
	FindCategoryByName(name string, merchantID string) (*models.Category, error)
	FindCategoryByID(id string) (*models.Category, error)
	FindAllCategoriesByMerchantID(merchantID string) ([]*models.Category, error)
	UpdateCategory(category *models.Category) (*models.Category, error)
	DeleteCategory(id uint32) error
}

type CategoryRepositoryInstance struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &CategoryRepositoryInstance{
		DB: db,
	}
}

// Store a new category
// This function requires a category model to be passed in, which contains the detail of the category that will be stored.
// It returns the stored category or an error if the operation fails.
func (categoryRepo *CategoryRepositoryInstance) StoreCategory(category *models.Category) (*models.Category, error) {
	if err := categoryRepo.DB.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// Finding a category by its name and merchant ID
// This function requires the name of the category and the merchant ID to be passed in.
// It returns the category if found or an error if the operation fails.
func (categoryRepo *CategoryRepositoryInstance) FindCategoryByName(name string, merchantID string) (*models.Category, error) {
	var category models.Category
	if err := categoryRepo.DB.Where("name ILIKE ? AND merchant_id = ?", fmt.Sprintf("%%%s%%", name), merchantID).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Finding a category by its ID
// Requires the ID of the category to be passed in.
// It returns the category if found or an error if the operation fails.
func (categoryRepo *CategoryRepositoryInstance) FindCategoryByID(id string) (*models.Category, error) {
	var category models.Category
	if err := categoryRepo.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Finding all categories by merchant ID
// This function requires the merchant ID to be passed in.
// It returns a slice of categories associated with the merchant or an error if the operation fails.
func (categoryRepo *CategoryRepositoryInstance) FindAllCategoriesByMerchantID(merchantID string) ([]*models.Category, error) {
	categories := make([]*models.Category, 0)
	if err := categoryRepo.DB.Where("merchant_id = ?", merchantID).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// Update an existing category with the provided details
// This function requires an updated category model to be passed in.
// It returns the updated category or an error if the operation fails.
func (categoryRepo *CategoryRepositoryInstance) UpdateCategory(category *models.Category) (*models.Category, error) {
	if err := categoryRepo.DB.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// Delete a category by its ID
// This function requires the ID of the category to be passed in.
// It returns an error if the operation fails.
func (categoryRepo *CategoryRepositoryInstance) DeleteCategory(id uint32) error {
	var category models.Category
	if err := categoryRepo.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return err
	}

	if err := categoryRepo.DB.Delete(&category).Error; err != nil {
		return err
	}

	return nil
}
