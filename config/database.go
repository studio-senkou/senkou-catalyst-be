package config

import (
	"fmt"
	"log"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dsn *string

func init() {
	host := utils.GetEnv("DB_HOST", "localhost")
	port := utils.GetEnv("DB_PORT", "5432")
	username := utils.GetEnv("DB_USERNAME", "postgres")
	password := utils.GetEnv("DB_PASSWORD", "")
	database := utils.GetEnv("DB_NAME", "senkou_catalyst")

	dsnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database,
	)

	dsn = &dsnStr
}

func ConnectDB() {
	db, err := gorm.Open(postgres.Open(*dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Connection failed:", err)
	}

	migrationError := db.AutoMigrate(
		&models.User{},
		&models.Merchant{},
		&models.Category{},
		&models.Product{},
		&models.Subscription{},
		&models.SubscriptionPlan{},
		&models.UserSubscription{},
		&models.PredefinedCategory{},
	)

	if migrationError != nil {
		log.Fatal("Migration failed:", migrationError)
	}

	DB = db
}
