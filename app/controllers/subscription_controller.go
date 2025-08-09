package controllers

import (
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/integrations/midtrans"
	"senkou-catalyst-be/utils/response"
	"senkou-catalyst-be/utils/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type SubscriptionController struct {
	UserService              services.UserService
	SubscriptionService      services.SubscriptionService
	SubscriptionOrderService services.SubscriptionOrderService
	PaymentService           services.PaymentService
}

func NewSubscriptionController(userService services.UserService, subService services.SubscriptionService, subOrderService services.SubscriptionOrderService, paymentService services.PaymentService) *SubscriptionController {
	return &SubscriptionController{
		UserService:              userService,
		SubscriptionService:      subService,
		SubscriptionOrderService: subOrderService,
		PaymentService:           paymentService,
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

	if err := validator.Validate(c, createSubscriptionDTO); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	subscription, err := h.SubscriptionService.CreateNewSubscription(createSubscriptionDTO)

	if err != nil && err.Code == 500 {
		return response.InternalError(c, "Failed to create subscription", map[string]any{
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
		return response.BadRequest(c, "Cannot continue to create subscription plan", "Failed to parse subscription plan ID")
	}

	planRequest := new(dtos.CreateSubscriptionPlanDTO)

	if err := validator.Validate(c, planRequest); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	if err := h.SubscriptionService.CreateSubscriptionPlan(planRequest, uint32(planID)); err != nil {
		if err.Code == 400 {
			return response.BadRequest(c, "Subscription plan already exists", map[string]any{
				"errors": err.Details,
			})
		}

		return response.InternalError(c, "Failed to create subscription plan", map[string]any{
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
		return response.BadRequest(c, "Cannot continue to subscribe user subscription", "Failed to parse user ID")
	}

	subIDStr := fmt.Sprintf("%v", c.Params("subID"))
	subID, err := strconv.ParseUint(subIDStr, 10, 32)

	if subID == 0 || err != nil {
		return response.BadRequest(c, "Cannot continue to subscribe user subscription", "Failed to parse subscription ID")
	}

	subscribeRequest := new(dtos.CreateSubscriptionOrderDTO)

	if err := validator.Validate(c, subscribeRequest); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	subscription, subError := h.SubscriptionService.GetSubscriptionByID(uint32(subID))

	if subError != nil {
		return response.NotFound(c, "Cannot continue to subscribe user into subscription due to not found subscription")
	}

	user, userError := h.UserService.GetUserDetail(uint32(userID))

	if userError != nil {
		return response.NotFound(c, "Cannot continue to subscribe user into subscription due to error getting user details")
	}

	subOrder, subOrderError := h.SubscriptionOrderService.GetOrderByUserAndSubscription(user.ID, subscription.ID)
	if subOrderError != nil {
		if subOrderError.Code != 404 {
			return response.InternalError(c, "Failed to get subscription order", subOrderError.Details)
		}
	}

	if subOrder != nil && subOrder.Status == string(midtrans.PaymentStatusSettled) {
		return response.BadRequest(c, "User already has an active subscription", "User already subscribed to this subscription")
	}

	chargeReceipt, paymentErr := h.PaymentService.CreatePayment(user, &midtrans.PaymentMethodConfig{
		PaymentType: subscribeRequest.PaymentType,
		Channel:     subscribeRequest.PaymentChannel,
	}, subscribeRequest.Amount)
	if paymentErr != nil {
		return response.InternalError(c, "Failed to create payment", paymentErr)
	}

	appError := h.SubscriptionOrderService.CreateNewSubscriptionOrder(uint32(userID), uint32(subID), subscribeRequest)
	if appError != nil {
		return response.InternalError(c, "Failed to create subscription order", appError.Details)
	}

	// if err := h.SubscriptionService.SubscribeUserToSubscription(uint32(userID), uint32(subID)); err != nil {
	// return response.InternalError(c, "Failed to subscribe user to subscription", "User could not be subscribed to subscription")
	// switch err.Code {
	// case 400:
	// 	return response.BadRequest(c, "User already has an active subscription", "Could not subscribe user to subscription")
	// case 404:
	// 	return response.NotFound(c, "Subscription not found")
	// }

	// return response.InternalError(c, "Failed to subscribe user to subscription", "User could not be subscribed to subscription")
	// }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User subscribed to subscription successfully",
		"data": fiber.Map{
			"payment_request": chargeReceipt,
		},
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
		return response.InternalError(c, "Failed to get subscriptions", "Unable to retrieve subscriptions")
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
		return response.BadRequest(c, "Invalid subscription ID", map[string]any{
			"error": "Invalid subscription ID",
		})
	}

	updateSubscriptionDTO := new(dtos.UpdateSubscriptionDTO)

	if err := validator.Validate(c, updateSubscriptionDTO); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	if err := h.SubscriptionService.UpdateSubscription(updateSubscriptionDTO, uint32(subID)); err != nil {
		return response.InternalError(c, "Failed to update subscription", "Could not update subscription")
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
		return response.BadRequest(c, "Invalid subscription ID", map[string]any{
			"error": "Invalid subscription ID",
		})
	}

	if err := h.SubscriptionService.DeleteSubscription(uint32(subID)); err != nil {
		return response.InternalError(c, "Failed to delete subscription", map[string]any{
			// "error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Subscription deleted successfully",
	})
}
