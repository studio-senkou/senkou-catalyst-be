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

	if err != nil {
		return throw.InternalError(c, "Failed to create subscription", map[string]any{
			"error": err.Error(),
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
	planID, err := strconv.ParseUint(planIDStr, 10, 64)

	if err != nil {
		return throw.BadRequest(c, "Invalid plan ID", map[string]any{
			"error": "Invalid plan ID",
		})
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
		return throw.InternalError(c, "Failed to create subscription plan", map[string]any{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Subscription plan created successfully",
	})
}
