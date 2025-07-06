package repositories

import (
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

type MerchantRepository interface {
	Create(merchant *models.Merchant) (*models.Merchant, error)
	FindMerchantsByUserID(userID uint) (*[]models.Merchant, error)
}

type merchantRepository struct {
	DB *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{
		DB: db,
	}
}

func (r *merchantRepository) Create(merchant *models.Merchant) (*models.Merchant, error) {
	if err := r.DB.Create(merchant).Error; err != nil {
		return nil, err
	}

	return merchant, nil
}

func (r *merchantRepository) FindMerchantsByUserID(userID uint) (*[]models.Merchant, error) {
	var merchants []models.Merchant

	if err := r.DB.
		Where("owner_id = ?", userID).
		Omit("Owner").
		Find(&merchants).Error; err != nil {
		return nil, err
	}

	if len(merchants) == 0 {
		return nil, nil
	}

	return &merchants, nil
}
