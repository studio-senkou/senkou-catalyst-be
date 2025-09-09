package repositories

import (
	"senkou-catalyst-be/app/models"
	"time"

	"gorm.io/gorm"
)

type EmailActivationRepository interface {
	Create(activation *models.EmailActivationToken) (*models.EmailActivationToken, error)
	FindByToken(token string) (*models.EmailActivationToken, error)
	Update(activation *models.EmailActivationToken) (*models.EmailActivationToken, error)
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
func (r *EmailActivationRepositoryInstance) FindByToken(token string) (*models.EmailActivationToken, error) {
	activation := new(models.EmailActivationToken)

	if err := r.DB.Where("token = ? AND expires_at > ?", token, time.Now()).First(activation).Error; err != nil {
		return nil, err
	}

	return activation, nil
}

func (r *EmailActivationRepositoryInstance) Update(activation *models.EmailActivationToken) (*models.EmailActivationToken, error) {
	if err := r.DB.Save(activation).Error; err != nil {
		return nil, err
	}

	return activation, nil
}
