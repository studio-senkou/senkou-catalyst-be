package controllers

import (
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/integrations/midtrans"
	"senkou-catalyst-be/utils/converter"
	"senkou-catalyst-be/utils/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PaymentController struct {
	PaymentService services.PaymentService
}

func NewPaymentController(paymentService services.PaymentService) *PaymentController {
	return &PaymentController{
		PaymentService: paymentService,
	}
}

func (p *PaymentController) PaymentNotifications(c *fiber.Ctx) error {
	midtransNotification := new(dtos.BaseMidtransNotification)

	if err := c.BodyParser(midtransNotification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body: " + err.Error(),
		})
	}

	var settlementTime, transactionTime *time.Time

	if midtransNotification.SettlementTime != "" {
		convertedTime, err := converter.ParseMidtransTime(midtransNotification.SettlementTime)
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid settlement time: " + err.Error(),
			})
		}

		settlementTime = &convertedTime
	}

	if midtransNotification.TransactionTime != "" {
		convertedTime, err := converter.ParseMidtransTime(midtransNotification.TransactionTime)
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid transaction time: " + err.Error(),
			})
		}

		transactionTime = &convertedTime
	}

	updatedPaymentTransaction := &models.PaymentTransaction{
		Status:          string(midtrans.ParseStatus(midtransNotification.TransactionStatus)),
		FraudStatus:     midtransNotification.FraudStatus,
		SignatureKey:    &midtransNotification.SignatureKey,
		TransactionTime: transactionTime,
		SettledAt:       settlementTime,
	}

	if err := p.PaymentService.UpdatePayment(midtransNotification.TransactionID, updatedPaymentTransaction); err != nil {
		return response.InternalError(c, "Failed to update payment transaction", err.Error())
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully update transaction",
	})
}
