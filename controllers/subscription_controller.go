package controllers

import (
	"fmt"
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/services"
	"senkou-catalyst-be/utils"
	"senkou-catalyst-be/utils/throw"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SubscriptionController struct {
	SubscriptionService *services.SubscriptionService
}

func NewSubscriptionController(subService *services.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{
		SubscriptionService: subService,
	}
}

// Create subscription
// @Summary Create a new subscription
// @Description Create a new subscription packet
// @Tags Subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param CreateSubscriptionDTO body dtos.CreateSubscriptionDTO true "Create subscription request object"
// @Success 201 {object} fiber.Map{message=string,data=fiber.Map{subscription=models.Subscription}}
// @Failure 400 {object} fiber.Map{message=string,errors=[]string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /subscriptions [post]
func (h *SubscriptionController) CreateSubscription(c *fiber.Ctx) error {
	createSubscriptionDTO := new(dtos.CreateSubscriptionDTO)

	if err := utils.Validate(c, createSubscriptionDTO); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return throw.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return throw.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	subscription, err := h.SubscriptionService.CreateNewSubscription(createSubscriptionDTO)

	if err != nil && err.Code == 500 {
		return throw.InternalError(c, "Failed to create subscription", map[string]any{
			"error": err.Details,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Subscription created successfully",
		"data": fiber.Map{
			"subscription": subscription,
		},
	})
}

// Create subscription plan
// @Summary Create a subscription plan
// @Description Create a new subscription plan for a subscription
// @Tags Subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param planID path string true "Plan ID"
// @Param CreateSubscriptionPlanDTO body dtos.CreateSubscriptionPlanDTO true "Create subscription plan request object"
// @Success 201 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{message=string,errors=[]string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /subscriptions/{planID}/plans [post]
func (h *SubscriptionController) CreateSubscriptionPlan(c *fiber.Ctx) error {
	planIDStr := fmt.Sprintf("%v", c.Locals("planID"))
	planID, err := strconv.ParseUint(planIDStr, 10, 32)

	if err != nil {
		return throw.BadRequest(c, "Cannot continue to create subscription plan", "Failed to parse subscription plan ID")
	}

	planRequest := new(dtos.CreateSubscriptionPlanDTO)

	if err := utils.Validate(c, planRequest); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return throw.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return throw.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	if err := h.SubscriptionService.CreateSubscriptionPlan(planRequest, uint32(planID)); err != nil {
		if err.Code == 400 {
			return throw.BadRequest(c, "Subscription plan already exists", map[string]any{
				"errors": err.Details,
			})
		}

		return throw.InternalError(c, "Failed to create subscription plan", map[string]any{
			"error": err.Details,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Subscription plan created successfully",
	})
}

// Subscribe to a subscription
// @Summary Subscribe to a subscription
// Description Subscribe a user to a subscription
// @Tags Subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param subID path string true "Subscription ID"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{message=string,error=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /subscriptions/{subID}/subscribe [post]
func (h *SubscriptionController) SubscribeSubscription(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseUint(userIDStr, 10, 32)

	if userID == 0 || err != nil {
		return throw.BadRequest(c, "Cannot continue to subscribe user subscription", "Failed to parse user ID")
	}

	subIDStr := fmt.Sprintf("%v", c.Params("subID"))
	subID, err := strconv.ParseUint(subIDStr, 10, 32)

	if err != nil {
		return throw.BadRequest(c, "Cannot continue to subscribe user subscription", "Failed to parse subscription ID")
	}

	if err := h.SubscriptionService.SubscribeUserToSubscription(uint32(userID), uint32(subID)); err != nil {
		switch err.Code {
		case 400:
			return throw.BadRequest(c, "User already has an active subscription", "Could not subscribe user to subscription")
		case 404:
			return throw.NotFound(c, "Subscription not found")
		}

		return throw.InternalError(c, "Failed to subscribe user to subscription", "User could not be subscribed to subscription")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User subscribed to subscription successfully",
	})
}

// Get all subscriptions
// @Summary Get all subscriptions
// @Description Retrieve all available subscriptions
// @Tags Subscription
// @Accept json
// @Produce json
// @Success 200 {object} fiber.Map{data=fiber.Map{subscriptions=[]models.Subscription}}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /subscriptions [get]
func (h *SubscriptionController) GetSubscriptions(c *fiber.Ctx) error {
	subscriptions, err := h.SubscriptionService.GetAllSubscriptions()

	if err != nil {
		return throw.InternalError(c, "Failed to get subscriptions", "Unable to retrieve subscriptions")
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"subscriptions": subscriptions,
		},
		"message": "Subscriptions retrieved successfully",
	})
}

// Update subscription
// @Summary Update a subscription
// @Description Update an existing subscription
// @Tags Subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param subID path string true "Subscription ID"
// @Param UpdateSubscriptionDTO body dtos.UpdateSubscriptionDTO true "Update subscription request object"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{message=string,errors=[]string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /subscriptions/{subID} [put]
func (h *SubscriptionController) UpdateSubscription(c *fiber.Ctx) error {
	subIDStr := fmt.Sprintf("%v", c.Params("subID"))
	subID, err := strconv.ParseUint(subIDStr, 10, 32)

	if err != nil {
		return throw.BadRequest(c, "Invalid subscription ID", map[string]any{
			"error": "Invalid subscription ID",
		})
	}

	updateSubscriptionDTO := new(dtos.UpdateSubscriptionDTO)

	if err := utils.Validate(c, updateSubscriptionDTO); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return throw.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return throw.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	if err := h.SubscriptionService.UpdateSubscription(updateSubscriptionDTO, uint32(subID)); err != nil {
		return throw.InternalError(c, "Failed to update subscription", "Could not update subscription")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Subscription updated successfully",
	})
}

// Delete subscription
// @Summary Delete a subscription
// @Description Delete a subscription by its ID
// @Tags Subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param subID path string true "Subscription ID"
// @Success 204 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{message=string,error=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /subscriptions/{subID} [delete]
func (h *SubscriptionController) DeleteSubscription(c *fiber.Ctx) error {
	subIDStr := fmt.Sprintf("%v", c.Params("subID"))
	subID, err := strconv.ParseUint(subIDStr, 10, 32)

	if err != nil {
		return throw.BadRequest(c, "Invalid subscription ID", map[string]any{
			"error": "Invalid subscription ID",
		})
	}

	if err := h.SubscriptionService.DeleteSubscription(uint32(subID)); err != nil {
		return throw.InternalError(c, "Failed to delete subscription", map[string]any{
			// "error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Subscription deleted successfully",
	})
}
