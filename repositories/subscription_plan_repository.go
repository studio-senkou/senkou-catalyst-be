package repositories

import (
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

type SubscriptionPlanRepository struct {
	DB *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) *SubscriptionPlanRepository {
	return &SubscriptionPlanRepository{
		DB: db,
	}
}

// Store a new subscription plan
// This function saves a new subscription plan to the database
func (r *SubscriptionPlanRepository) StoreNewPlan(plan *models.SubscriptionPlan) error {
	if err := r.DB.Create(plan).Error; err != nil {
		return err
	}
	return nil
}

func (r *SubscriptionPlanRepository) IsPlanExists(subID uint32, planName string) (bool, error) {
	plan := new(models.SubscriptionPlan)

	err := r.DB.Where("name = ? AND subscription_id = ?", planName, subID).First(&plan).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
