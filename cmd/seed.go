package main

import (
	"fmt"
	"senkou-catalyst-be/config"
	"senkou-catalyst-be/database"
)

// @title Catalyst API Seeder
// @version 1.0
// @description Seeder for Catalyst API
func main() {
	config.ConnectDB()
	database.Seed(config.DB)

	fmt.Println("Database seeding completed successfully!")
}
