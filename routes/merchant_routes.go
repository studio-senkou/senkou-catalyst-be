package routes

import (
	"senkou-catalyst-be/config"
	"senkou-catalyst-be/controllers"
	"senkou-catalyst-be/platform/middlewares"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/services"

	"github.com/gofiber/fiber/v2"
)

func MerchantRoutes(app *fiber.App) {
	merchantRepo := repositories.NewMerchantRepository(config.DB)
	merchantService := services.NewMerchantService(merchantRepo)
	merchantController := controllers.NewMerchantController(merchantService)

	app.Post("/merchants", middlewares.JWTProtected, merchantController.CreateMerchant)
	app.Get("/merchants", middlewares.JWTProtected, merchantController.GetUserMerchants)
	app.Get("/merchants/:id", middlewares.JWTProtected, merchantController.GetMerchantByID)
	app.Put("/merchants/:id", middlewares.JWTProtected, merchantController.UpdateMerchant)
	app.Delete("/merchants/:id", middlewares.JWTProtected, merchantController.DeleteMerchant)
}
