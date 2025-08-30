package repositories

import (
	"senkou-catalyst-be/app/models"

	"gorm.io/gorm"
)

type ProductInteractionRepository interface {
	StoreProductInteractionLog(interaction *models.ProductMetric) error
}

type ProductInteractionRepositoryInstance struct {
	db *gorm.DB
}

func NewProductInteractionRepository(db *gorm.DB) ProductInteractionRepository {
	return &ProductInteractionRepositoryInstance{
		db: db,
	}
}

// StoreProductInteractionLog stores a new product interaction log in the database
func (r *ProductInteractionRepositoryInstance) StoreProductInteractionLog(interaction *models.ProductMetric) error {
	return r.db.Create(interaction).Error
}
