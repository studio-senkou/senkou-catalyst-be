package routes

import (
	"senkou-catalyst-be/app/controllers"
	"senkou-catalyst-be/platform/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitUserRoutes(app *fiber.App, userController *controllers.UserController) {
	app.Post(
		"/users",
		userController.CreateUser,
	)
	app.Post(
		"/users/activate",
		userController.ActivateAccount,
	)
	app.Get(
		"/users",
		middlewares.JWTProtected,
		userController.GetUsers,
	)
	app.Get(
		"/users/me",
		middlewares.JWTProtected,
		userController.GetUserDetail,
	)
}
