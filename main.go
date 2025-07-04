package main

import (
	"senkou-catalyst-be/config"
	"senkou-catalyst-be/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "senkou-catalyst-be/docs" // penting: import untuk dokumentasi Swag
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

	// Endpoint Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Routes
	routes.UserRoutes(app)

	// Jalankan server
	app.Listen(":3000")
}
