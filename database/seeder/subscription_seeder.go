package seeder

import (
	"errors"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/repositories"

	"gorm.io/gorm"
)

type SubscriptionData struct {
	Name        string
	Description string
	Price       float64
	Duration    int
	Plans       []PlanData
}

type PlanData struct {
	Name  string
	Value string
}

func SeedSubscriptions(db *gorm.DB) error {
	subsRepository := repositories.NewSubscriptionRepository(db)
	planRepository := repositories.NewSubscriptionPlanRepository(db)

	subscriptions := []SubscriptionData{
		{
			Name:        "Free tier",
			Description: "The minimal subscription plan",
			Price:       0,
			Duration:    28,
			Plans: []PlanData{
				{"Subscription-Product-Slot", "20"},
				{"Subscription-Category-Limit", "4"},
				{"Subscription-Analytics", "false"},
				{"Subscription-Interaction-Metrics", "false"},
				{"Subscription-Merchant-Template-Customize", "false"},
			},
		},
		{
			Name:        "Content Creator",
			Description: "Subscription plan for content creators",
			Price:       10000,
			Duration:    28,
			Plans: []PlanData{
				{"Subscription-Product-Slot", "100"},
				{"Subscription-Category-Limit", "10"},
				{"Subscription-Analytics", "true"},
				{"Subscription-Interaction-Metrics", "false"},
				{"Subscription-Merchant-Template-Customize", "true"},
			},
		},
		{
			Name:        "Business",
			Description: "Subscription plan for businesses",
			Price:       30000,
			Duration:    28,
			Plans: []PlanData{
				{"Subscription-Product-Slot", "99999"},
				{"Subscription-Category-Limit", "999"},
				{"Subscription-Analytics", "true"},
				{"Subscription-Interaction-Metrics", "true"},
				{"Subscription-Merchant-Template-Customize", "true"},
			},
		},
	}

	for _, subData := range subscriptions {
		sub, err := createOrGetSubscription(subsRepository, subData)
		if err != nil {
			panic(err)
		}

		if err := createSubscriptionPlans(planRepository, sub.ID, subData.Plans); err != nil {
			panic(err)
		}
	}

	return nil
}

func createOrGetSubscription(repo repositories.SubscriptionRepository, data SubscriptionData) (*models.Subscription, error) {
	sub := &models.Subscription{
		Name:        data.Name,
		Description: data.Description,
		Price:       float32(data.Price),
		Duration:    int16(data.Duration),
	}

	return repo.StoreNewSubscription(sub)
}

func createSubscriptionPlans(repo repositories.SubscriptionPlanRepository, subID uint32, plans []PlanData) error {
	for _, planData := range plans {
		exists, err := repo.IsPlanExists(subID, planData.Name)

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if !exists || errors.Is(err, gorm.ErrRecordNotFound) {
			plan := &models.SubscriptionPlan{
				SubID: subID,
				Name:  planData.Name,
				Value: planData.Value,
			}

			if err := repo.StoreNewPlan(plan); err != nil {
				return err
			}
		}
	}
	return nil
}
