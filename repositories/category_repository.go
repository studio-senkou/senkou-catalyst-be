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
	FindAllCategoriesByMerchantUsername(username string) ([]*models.Category, error)
	FindCategoryByNameAndMerchantUsername(name, username string) (*models.Category, error)
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
func (c *CategoryRepositoryInstance) StoreCategory(category *models.Category) (*models.Category, error) {
	if err := c.DB.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// Finding a category by its name and merchant ID
// This function requires the name of the category and the merchant ID to be passed in.
// It returns the category if found or an error if the operation fails.
func (c *CategoryRepositoryInstance) FindCategoryByName(name string, merchantID string) (*models.Category, error) {
	var category models.Category
	if err := c.DB.Where("name ILIKE ? AND merchant_id = ?", fmt.Sprintf("%%%s%%", name), merchantID).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Finding a category by its ID
// Requires the ID of the category to be passed in.
// It returns the category if found or an error if the operation fails.
func (c *CategoryRepositoryInstance) FindCategoryByID(id string) (*models.Category, error) {
	var category models.Category
	if err := c.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Finding all categories by merchant ID
// This function requires the merchant ID to be passed in.
// It returns a slice of categories associated with the merchant or an error if the operation fails.
func (c *CategoryRepositoryInstance) FindAllCategoriesByMerchantID(merchantID string) ([]*models.Category, error) {
	categories := make([]*models.Category, 0)
	if err := c.DB.Where("merchant_id = ?", merchantID).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// Finding all categories by merchant username
// This function requires the merchant username to be passed in.
// It returns a slice of categories associated with the merchant or an error if the operation fails.
func (c *CategoryRepositoryInstance) FindAllCategoriesByMerchantUsername(username string) ([]*models.Category, error) {
	categories := make([]*models.Category, 0)
	if err := c.DB.Joins("JOIN merchants ON merchants.id = categories.merchant_id").Where("merchants.username = ?", username).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// Finding a category by its name and merchant username
// This function requires the name of the category and the merchant username to be passed in.
// It returns the category if found or an error if the operation fails.
func (c *CategoryRepositoryInstance) FindCategoryByNameAndMerchantUsername(name, username string) (*models.Category, error) {
	var category models.Category
	if err := c.DB.Joins("JOIN merchants ON merchants.id = categories.merchant_id").Where("categories.name ILIKE ? AND merchants.username = ?", fmt.Sprintf("%%%s%%", name), username).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Update an existing category with the provided details
// This function requires an updated category model to be passed in.
// It returns the updated category or an error if the operation fails.
func (c *CategoryRepositoryInstance) UpdateCategory(category *models.Category) (*models.Category, error) {
	if err := c.DB.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// Delete a category by its ID
// This function requires the ID of the category to be passed in.
// It returns an error if the operation fails.
func (c *CategoryRepositoryInstance) DeleteCategory(id uint32) error {
	var category models.Category
	if err := c.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return err
	}

	if err := c.DB.Delete(&category).Error; err != nil {
		return err
	}

	return nil
}
