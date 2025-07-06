package routes

import (
	"senkou-catalyst-be/config"
	"senkou-catalyst-be/controllers"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/services"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	authRepo := repositories.NewAuthRepository(config.DB)
	authService := services.NewAuthService(&authRepo)
	userRepo := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepo)
	authController := controllers.NewAuthController(authService, userService)

	app.Post("/auth/login", authController.Login)
}
