package seeder

import (
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/repositories"

	"gorm.io/gorm"
)

func SubscriptionSeeder(db *gorm.DB) {

	subsRepository := repositories.NewSubscriptionRepository(db)
	planRepository := repositories.NewSubscriptionPlanRepository(db)

	sub, err := subsRepository.StoreNewSubscription(&models.Subscription{
		Name:        "Standard",
		Description: "The minimal subscription plan",
		Price:       1000,
		Duration:    28,
	})
	if err != nil {
		panic(err)
	}

	planRepository.IsPlanExists(sub.ID, "Subscription-Product-Slot")
}
