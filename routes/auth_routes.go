package routes

import (
	"senkou-catalyst-be/app/controllers"
	"senkou-catalyst-be/platform/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitAuthRoutes(app *fiber.App, authController *controllers.AuthController) {
	app.Post("/auth/login", authController.Login)
	app.Put("/auth/refresh", authController.RefreshToken)
	app.Delete("/auth/logout", middlewares.JWTProtected, authController.Logout)
}
