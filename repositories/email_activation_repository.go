package repositories

import (
	"senkou-catalyst-be/app/models"

	"gorm.io/gorm"
)

type EmailActivationRepository interface {
	Create(activation *models.EmailActivationToken) (*models.EmailActivationToken, error)
}

type EmailActivationRepositoryInstance struct {
	DB *gorm.DB
}

func NewEmailActivationRepository(db *gorm.DB) EmailActivationRepository {
	return &EmailActivationRepositoryInstance{
		DB: db,
	}
}

func (r *EmailActivationRepositoryInstance) Create(activation *models.EmailActivationToken) (*models.EmailActivationToken, error) {
	if err := r.DB.Create(activation).Error; err != nil {
		return nil, err
	}

	return activation, nil
}
