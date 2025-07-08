package routes

import (
	"senkou-catalyst-be/config"
	"senkou-catalyst-be/controllers"
	"senkou-catalyst-be/platform/middlewares"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/services"

	"github.com/gofiber/fiber/v2"
)

func CategoryRoutes(app *fiber.App) {
	categoryRepo := repositories.NewCategoryRepository(config.DB)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryController := controllers.NewCategoryController(categoryService)

	app.Post("/merchants/:merchantID/categories", middlewares.JWTProtected, categoryController.CreateCategory)
	app.Get("/merchants/:merchantID/categories", middlewares.JWTProtected, categoryController.GetCategories)
	app.Put("/merchants/:merchantID/categories/:categoryID", middlewares.JWTProtected, categoryController.UpdateCategory)
	app.Delete("/merchants/:merchantID/categories/:categoryID", middlewares.JWTProtected, categoryController.DeleteCategory)
}
