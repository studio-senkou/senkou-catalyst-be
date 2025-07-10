package routes

import (
	"senkou-catalyst-be/config"
	"senkou-catalyst-be/controllers"
	"senkou-catalyst-be/platform/middlewares"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/services"

	"github.com/gofiber/fiber/v2"
)

func PredefinedCategoryRoutes(app *fiber.App) {
	predefinedCategoryRepo := repositories.NewPredefinedCategoryRepository(config.DB)
	predefinedCategoryService := services.NewPredefinedCategoryService(predefinedCategoryRepo)
	predefinedCategoryController := controllers.NewPredefinedCategoryController(predefinedCategoryService)

	pdRoute := app.Group("/predefined-categories", middlewares.JWTProtected, middlewares.RoleMiddleware("admin"))

	pdRoute.Post("/", predefinedCategoryController.StoreCategory)
}
