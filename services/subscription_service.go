package services

import (
	"errors"
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"
)

type SubscriptionService struct {
	SubscriptionRepository     *repositories.SubscriptionRepository
	SubscriptionPlanRepository *repositories.SubscriptionPlanRepository
}

func NewSubscriptionService(
	subRepository *repositories.SubscriptionRepository,
	subPlanRepository *repositories.SubscriptionPlanRepository,
) *SubscriptionService {
	return &SubscriptionService{
		SubscriptionRepository:     subRepository,
		SubscriptionPlanRepository: subPlanRepository,
	}
}

func (s *SubscriptionService) CreateSubscription() error {
	return nil
}

func (s *SubscriptionService) CreateNewSubscription(request *dtos.CreateSubscriptionDTO) (*models.Subscription, error) {
	subscription := &models.Subscription{
		Name:        request.Name,
		Description: request.Description,
		Price:       float32(request.Price),
		Duration:    request.Duration,
	}

	subscription, err := s.SubscriptionRepository.StoreNewSubscription(subscription)

	if err != nil {
		return nil, errors.New("failed to create subscription")
	}

	return subscription, nil
}

// Create a new subscription plan
// This function will create a new subscription plan and return error if it already exists
func (s *SubscriptionService) CreateSubscriptionPlan(planRequest *dtos.CreateSubscriptionPlanDTO, subID uint32) error {
	if exists, err := s.SubscriptionPlanRepository.IsPlanExists(subID, planRequest.Name); err != nil || exists {
		return errors.New("subscription plan already exists")
	}

	plan := &models.SubscriptionPlan{
		SubID: subID,
		Name:  planRequest.Name,
		Value: planRequest.Value,
	}

	if err := s.SubscriptionPlanRepository.StoreNewPlan(plan); err != nil {
		return errors.New("failed to create subscription plan")
	}

	return nil
}
