package routes

import (
	"senkou-catalyst-be/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func InitStorageRoutes(app *fiber.App, storageController *controllers.StorageController) {
	app.Get("/files/*", storageController.GetFromStorage)
}
