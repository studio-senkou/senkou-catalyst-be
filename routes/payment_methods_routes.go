package routes

import (
	"senkou-catalyst-be/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func InitPaymentMethodsRoutes(app *fiber.App, controller *controllers.PaymentMethodsController) {
	api := app.Group("/payment-methods")

	api.Get("/", controller.GetAllAvailablePaymentMethods)
	api.Get("/types", controller.GetPaymentMethodTypes)
	api.Get("/type/:type", controller.GetPaymentMethodsByType)
}
