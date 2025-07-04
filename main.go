package main

import (
	"fmt"
	"senkou-catalyst-be/config"
	"senkou-catalyst-be/routes"
	"senkou-catalyst-be/utils"

	_ "senkou-catalyst-be/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title Fiber API with Service-Repository
// @version 1.0
// @description Fiber REST API using service-repo pattern
// @host localhost:3000
// @BasePath /

// @schemes http
func main() {
	app := fiber.New()
	config.ConnectDB()

	app.Get("/swagger/*", swagger.HandlerDefault)

	routes.UserRoutes(app)

	app.Listen(fmt.Sprintf(":%s", utils.GetEnv("APP_PORT", "8080")))
}
