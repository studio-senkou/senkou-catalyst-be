package routes

import (
	"senkou-catalyst-be/controllers"
	"senkou-catalyst-be/platform/middlewares"

	"github.com/gofiber/fiber/v2"
)

func MerchantRoutes(app *fiber.App, merchantController *controllers.MerchantController) {
	app.Post("/merchants", middlewares.JWTProtected, merchantController.CreateMerchant)
	app.Get("/merchants", middlewares.JWTProtected, middlewares.RoleMiddleware("admin"), merchantController.GetUserMerchants)
	app.Get("/merchants/:id", middlewares.JWTProtected, merchantController.GetMerchantByID)
	app.Put("/merchants/:id", middlewares.JWTProtected, merchantController.UpdateMerchant)
	app.Delete("/merchants/:id", middlewares.JWTProtected, merchantController.DeleteMerchant)
}
