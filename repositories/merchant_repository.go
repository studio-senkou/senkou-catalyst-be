package repositories

import (
	"senkou-catalyst-be/app/models"

	"gorm.io/gorm"
)

type MerchantRepository interface {
	Create(merchant *models.Merchant) (*models.Merchant, error)
	FindMerchantsByUserID(userID uint32) ([]*models.Merchant, error)
	FindByID(merchantID string) (*models.Merchant, error)
	UpdateMerchant(merchantID string, updateData *models.Merchant) (*models.Merchant, error)
	DeleteMerchant(merchantID string) error
}

type MerchantRepositoryInstance struct {
	DB *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &MerchantRepositoryInstance{
		DB: db,
	}
}

// Create a new merchant
// This function creates a new merchant in the database
// It returns the created merchant or an error if the creation fails
func (r *MerchantRepositoryInstance) Create(merchant *models.Merchant) (*models.Merchant, error) {
	if err := r.DB.Create(merchant).Error; err != nil {
		return nil, err
	}

	return merchant, nil
}

// Find merchants by user ID
// This function retrieves all merchants associated with a specific user ID
// It returns a slice of merchants or an error if the retrieval fails
func (r *MerchantRepositoryInstance) FindMerchantsByUserID(userID uint32) ([]*models.Merchant, error) {
	merchants := make([]*models.Merchant, 0)

	if err := r.DB.
		Where("owner_id = ?", userID).
		Omit("Owner").
		Find(&merchants).Error; err != nil {
		return nil, err
	}

	if len(merchants) == 0 {
		return nil, nil
	}

	return merchants, nil
}

// Find a merchant by ID
// This function retrieves a merchant by its ID
// It returns the merchant or an error if the retrieval fails
func (r *MerchantRepositoryInstance) FindByID(merchantID string) (*models.Merchant, error) {
	var merchant models.Merchant

	if err := r.DB.
		Where("id = ?", merchantID).
		Omit("Owner").
		First(&merchant).Error; err != nil {
		return nil, err
	}

	return &merchant, nil
}

// Update a merchant
// This function updates an existing merchant with the provided data
// It returns the updated merchant or an error if the update fails
func (r *MerchantRepositoryInstance) UpdateMerchant(merchantID string, updateData *models.Merchant) (*models.Merchant, error) {
	var merchant models.Merchant

	if err := r.DB.
		Where("id = ?", merchantID).
		Omit("Owner").
		First(&merchant).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Model(&merchant).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &merchant, nil
}

// Delete a merchant
// This function deletes a merchant by its ID
// It returns an error if the deletion fails
func (r *MerchantRepositoryInstance) DeleteMerchant(merchantID string) error {
	var merchant models.Merchant

	if err := r.DB.
		Where("id = ?", merchantID).
		Omit("Owner").
		First(&merchant).Error; err != nil {
		return err
	}

	if err := r.DB.Delete(&merchant).Error; err != nil {
		return err
	}

	return nil
}
