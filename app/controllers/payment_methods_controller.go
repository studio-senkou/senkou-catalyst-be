package controllers

import (
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/utils/response"

	"github.com/gofiber/fiber/v2"
)

type PaymentMethodsController struct {
	PaymentMethodsService services.PaymentMethodsService
}

func NewPaymentMethodsController(paymentMethodsService services.PaymentMethodsService) *PaymentMethodsController {
	return &PaymentMethodsController{
		PaymentMethodsService: paymentMethodsService,
	}
}

// GetAllAvailablePaymentMethods retrieves all available payment methods
// @Summary Get all available payment methods
// @Description Get all available payment methods
// @Tags Payment Methods
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /payment-methods [get]
func (pc *PaymentMethodsController) GetAllAvailablePaymentMethods(c *fiber.Ctx) error {
	methods, err := pc.PaymentMethodsService.GetAllAvailablePaymentMethods()
	if err != nil {
		return response.InternalError(c, "Failed to retrieve payment methods", err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Available payment methods retrieved successfully",
		"data": fiber.Map{
			"payment_methods": methods,
		},
	})
}

// GetPaymentMethodsByType retrieves payment methods by type
// @Summary Get payment methods by type
// @Description Get payment methods by type
// @Tags Payment Methods
// @Param type path string true "Payment Method Type"
// @Success 200 {object} fiber.Map{message=string,data=fiber.Map{payment_methods=[]interface{}}}
// @Failure 500 {object} fiber.Map{success=bool,error=string}
// @Router /payment-methods/{type} [get]
func (pc *PaymentMethodsController) GetPaymentMethodsByType(c *fiber.Ctx) error {
	paymentType := c.Params("type")

	methods, err := pc.PaymentMethodsService.GetPaymentMethodsByType(paymentType)
	if err != nil {
		return response.InternalError(c, "Failed to retrieve payment methods by type", err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Payment methods by type retrieved successfully",
		"data": fiber.Map{
			"payment_methods": methods,
		},
	})
}

// GetPaymentMethodTypes retrieves all available payment method types
// @Summary Get all available payment method types
// @Description Get all available payment method types
// @Tags Payment Methods
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /payment-methods/types [get]
func (pc *PaymentMethodsController) GetPaymentMethodTypes(c *fiber.Ctx) error {
	types, err := pc.PaymentMethodsService.GetPaymentMethodTypes()
	if err != nil {
		return response.InternalError(c, "Failed to retrieve payment method types", err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Payment method types retrieved successfully",
		"data": fiber.Map{
			"payment_method_types": types,
		},
	})
}
