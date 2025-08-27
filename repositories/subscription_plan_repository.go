package repositories

import (
	"senkou-catalyst-be/app/models"

	"gorm.io/gorm"
)

type SubscriptionPlanRepository interface {
	StoreNewPlan(plan *models.SubscriptionPlan) error
	IsPlanExists(subID uint32, planName string) (bool, error)
}

type SubscriptionPlanRepositoryInstance struct {
	DB *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) SubscriptionPlanRepository {
	return &SubscriptionPlanRepositoryInstance{
		DB: db,
	}
}

// Store a new subscription plan
// This function saves a new subscription plan to the database
func (r *SubscriptionPlanRepositoryInstance) StoreNewPlan(plan *models.SubscriptionPlan) error {
	if err := r.DB.Create(plan).Error; err != nil {
		return err
	}
	return nil
}

// Check if the subscription plan already exists
// This function checks for the existence of a subscription plan by its name and subscription ID
func (r *SubscriptionPlanRepositoryInstance) IsPlanExists(subID uint32, planName string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.SubscriptionPlan{}).
		Where("name = ? AND sub_id = ?", planName, subID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
