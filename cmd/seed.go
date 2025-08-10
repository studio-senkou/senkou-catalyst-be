package main

import (
	"fmt"
	"senkou-catalyst-be/database/seeder"
	"senkou-catalyst-be/platform/config"
)

// @title Catalyst API Seeder
// @version 1.0
// @description Seeder for Catalyst API
func main() {
	config.ConnectDB()
	seeder.Seed(config.DB)

	fmt.Println("Database seeding completed successfully!")
}
