package repositories

import (
	"senkou-catalyst-be/app/models"

	"gorm.io/gorm"
)

type PredefinedCategoryRepository interface {
	StoreCategory(pdCategory *models.PredefinedCategory) error
	FindByName(name string) (*models.PredefinedCategory, error)
	FindAll() ([]*models.PredefinedCategory, error)
	UpdateByID(pcID uint32, updatedCategory *models.PredefinedCategory) error
	RemoveByID(pcID uint32) error
}

type PredefinedCategoryRepositoryInstance struct {
	DB *gorm.DB
}

func NewPredefinedCategoryRepository(db *gorm.DB) PredefinedCategoryRepository {
	return &PredefinedCategoryRepositoryInstance{
		DB: db,
	}
}

// Store a new predefined category
// This function requires a predefined category model to be passed in, which contains the detail of the category.
func (r *PredefinedCategoryRepositoryInstance) StoreCategory(pdCategory *models.PredefinedCategory) error {
	return r.DB.Create(pdCategory).Error
}

// Find a predefined category by its name
// This function requires the name of the category to be passed in.
// It returns the category if found or an error if the operation fails.
func (r *PredefinedCategoryRepositoryInstance) FindByName(name string) (*models.PredefinedCategory, error) {
	category := new(models.PredefinedCategory)

	if err := r.DB.Select("id", "name", "description", "image_url").
		Where("name = ?", name).
		First(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// Find all predefined categories
// This function retrieves all predefined categories from the database
// It returns a slice of predefined categories or an error if the operation fails.
func (r *PredefinedCategoryRepositoryInstance) FindAll() ([]*models.PredefinedCategory, error) {
	categories := make([]*models.PredefinedCategory, 0)

	if err := r.DB.Select("id", "name", "description", "image_url").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// Update a predefined category by its ID
// This function requires the ID of the category and the updated data to be passed in.
// It returns an error if the operation fails.
func (r *PredefinedCategoryRepositoryInstance) UpdateByID(pcID uint32, updatedCategory *models.PredefinedCategory) error {
	return r.DB.Model(&models.PredefinedCategory{}).Where("id = ?", pcID).Updates(updatedCategory).Error
}

// Remove a predefined category by its ID
// This function requires the ID of the category to be passed in
// It returns an error if the operation fails.
func (r *PredefinedCategoryRepositoryInstance) RemoveByID(pcID uint32) error {
	return r.DB.Delete(&models.PredefinedCategory{}, pcID).Error
}
