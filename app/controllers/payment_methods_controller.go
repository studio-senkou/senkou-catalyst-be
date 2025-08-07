package controllers

import (
	"senkou-catalyst-be/app/services"

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

func (pc *PaymentMethodsController) GetAllAvailablePaymentMethods(c *fiber.Ctx) error {
	methods, err := pc.PaymentMethodsService.GetAllAvailablePaymentMethods()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Available payment methods retrieved successfully",
		"data": fiber.Map{
			"payment_methods": methods,
		},
	})
}

func (pc *PaymentMethodsController) GetPaymentMethodsByType(c *fiber.Ctx) error {
	paymentType := c.Params("type")

	methods, err := pc.PaymentMethodsService.GetPaymentMethodsByType(paymentType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment methods by type retrieved successfully",
		"data": fiber.Map{
			"payment_methods": methods,
		},
	})
}

func (pc *PaymentMethodsController) GetPaymentMethodTypes(c *fiber.Ctx) error {
	types, err := pc.PaymentMethodsService.GetPaymentMethodTypes()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment method types retrieved successfully",
		"data": fiber.Map{
			"payment_method_types": types,
		},
	})
}
