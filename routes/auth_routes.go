package routes

import (
	"senkou-catalyst-be/controllers"
	"senkou-catalyst-be/platform/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authController *controllers.AuthController) {
	app.Post("/auth/login", authController.Login)
	app.Put("/auth/refresh", authController.RefreshToken)
	app.Delete("/auth/logout", middlewares.JWTProtected, authController.Logout)
}
