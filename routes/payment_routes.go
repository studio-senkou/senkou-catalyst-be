package routes

import (
	"senkou-catalyst-be/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func InitPaymentRoutes(app *fiber.App, paymentController *controllers.PaymentController) {
	app.Post("/api/v1/payments/notifications", paymentController.PaymentNotifications)
}
