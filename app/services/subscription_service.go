package services

import (
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"
	"time"
)

type SubscriptionService interface {
	CreateNewSubscription(request *dtos.CreateSubscriptionDTO) (*models.Subscription, *errors.CustomError)
	SubscribeUserToSubscription(userID, subID uint32) *errors.CustomError
	CreateSubscriptionPlan(request *dtos.CreateSubscriptionPlanDTO, subID uint32) *errors.CustomError
	GetAllSubscriptions() ([]*models.Subscription, *errors.CustomError)
	GetSubscriptionByID(subID uint32) (*models.Subscription, *errors.CustomError)
	AssignFreeTierSubscription(userID uint32) *errors.CustomError
	UpdateSubscription(request *dtos.UpdateSubscriptionDTO, subID uint32) *errors.CustomError
	DeleteSubscription(subID uint32) *errors.CustomError
}

type SubscriptionServiceInstance struct {
	SubscriptionRepository     repositories.SubscriptionRepository
	SubscriptionPlanRepository repositories.SubscriptionPlanRepository
}

func NewSubscriptionService(
	subRepository repositories.SubscriptionRepository,
	subPlanRepository repositories.SubscriptionPlanRepository,
) SubscriptionService {
	return &SubscriptionServiceInstance{
		SubscriptionRepository:     subRepository,
		SubscriptionPlanRepository: subPlanRepository,
	}
}

// Create a new subscription
// This function will create a new subscription and return the created subscription
// It returns an error if the subscription could not be created
func (s *SubscriptionServiceInstance) CreateNewSubscription(request *dtos.CreateSubscriptionDTO) (*models.Subscription, *errors.CustomError) {
	subscription := &models.Subscription{
		Name:        request.Name,
		Description: request.Description,
		Price:       float32(request.Price),
		Duration:    request.Duration,
	}

	subscription, err := s.SubscriptionRepository.StoreNewSubscription(subscription)

	if err != nil {
		return nil, errors.Internal("Failed to create subscription", err.Error())
	}

	return subscription, nil
}

// Subscribe a user to a subscription
// This function will subscribe a user to a subscription by userID and subID
// It returns an error if the subscription does not exist or if the user could not be subscribed
func (s *SubscriptionServiceInstance) SubscribeUserToSubscription(userID, subID uint32) *errors.CustomError {
	subscription, err := s.SubscriptionRepository.FindByID(subID)
	if err != nil || subscription == nil {
		return errors.NotFound("Subscription not found")
	}

	if exist, err := s.SubscriptionRepository.VerifyUserHasActiveSubscription(userID, subID); err != nil || exist {
		return errors.BadRequest("User already has an active subscription", nil)
	}

	sub := &models.UserSubscription{
		UserID:        userID,
		SubID:         subID,
		IsActive:      false,
		PaymentStatus: "pending",
	}

	if err := s.SubscriptionRepository.SubscribeUser(sub); err != nil {
		return errors.Internal("Failed to subscribe user to subscription", err.Error())
	}

	return nil
}

// Assign a free tier subscription to a user
// This function will assign a free tier subscription to a user by userID
// It returns an error if the user could not be assigned the subscription
func (s *SubscriptionServiceInstance) AssignFreeTierSubscription(userID uint32) *errors.CustomError {
	freeTierSub, err := s.SubscriptionRepository.FindFreeTierSubscription()
	if err != nil || freeTierSub == nil {
		return errors.NotFound("Free tier subscription not found")
	}

	if err := s.SubscriptionRepository.SubscribeUser(&models.UserSubscription{
		UserID:        userID,
		SubID:         freeTierSub.ID,
		StartedAt:     time.Now(),
		ExpiredAt:     time.Now().AddDate(100, 0, 0),
		IsActive:      true,
		PaymentStatus: "settled",
	}); err != nil {
		return errors.Internal("Failed to assign free tier subscription", err.Error())
	}

	return nil
}

// Create a new subscription plan
// This function will create a new subscription plan and return error if it already exists
func (s *SubscriptionServiceInstance) CreateSubscriptionPlan(request *dtos.CreateSubscriptionPlanDTO, subID uint32) *errors.CustomError {
	if exists, err := s.SubscriptionPlanRepository.IsPlanExists(subID, request.Name); err != nil || exists {
		return errors.Conflict("Subscription plan already exists", nil)
	}

	plan := &models.SubscriptionPlan{
		SubID: subID,
		Name:  request.Name,
		Value: request.Value,
	}

	if err := s.SubscriptionPlanRepository.StoreNewPlan(plan); err != nil {
		return errors.Internal("Failed to create subscription plan", err.Error())
	}

	return nil
}

// Get all subscriptions
// This function retrieves all subscription from the database
// It returns a slice of subscription and an error if any
func (s *SubscriptionServiceInstance) GetAllSubscriptions() ([]*models.Subscription, *errors.CustomError) {
	subscriptions, err := s.SubscriptionRepository.FindAllSubscriptions()
	if err != nil {
		return nil, errors.Internal("Failed to retrieve subscriptions", err.Error())
	}

	return subscriptions, nil
}

func (s *SubscriptionServiceInstance) GetSubscriptionByID(subID uint32) (*models.Subscription, *errors.CustomError) {
	subscription, err := s.SubscriptionRepository.FindByID(subID)

	if err != nil || subscription == nil {
		return nil, errors.NotFound("Subscription not found")
	}

	return subscription, nil
}

// Update an existing subscription
// This function updates an existing subscription with the provided request data
// It returns an error if the subscription does not exist as if the update fails
func (s *SubscriptionServiceInstance) UpdateSubscription(request *dtos.UpdateSubscriptionDTO, subID uint32) *errors.CustomError {

	if subscription, err := s.SubscriptionRepository.FindByID(subID); err != nil || subscription == nil {
		return errors.NotFound("Subscription not found")
	}

	updatedSubscription := &models.Subscription{
		ID:          subID,
		Name:        *request.Name,
		Description: *request.Description,
		Price:       float32(*request.Price),
		Duration:    *request.Duration,
	}

	if _, err := s.SubscriptionRepository.UpdateSubscription(updatedSubscription); err != nil {
		return errors.Internal("Failed to update subscription", err.Error())
	}

	return nil
}

// Delete a subscription
// This function deletes a subscription by its ID
// It returns an error if the subscription does not exist or if the deletion fails
func (s *SubscriptionServiceInstance) DeleteSubscription(subID uint32) *errors.CustomError {
	if subscription, err := s.SubscriptionRepository.FindByID(subID); err != nil || subscription == nil {
		return errors.NotFound("Subscription not found")
	}

	subscription := &models.Subscription{ID: subID}

	if err := s.SubscriptionRepository.DeleteSubscription(subscription); err != nil {
		return errors.Internal("Failed to delete subscription", err.Error())
	}

	return nil
}
