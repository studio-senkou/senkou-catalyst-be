package config

import (
	"fmt"
	"log"
	"senkou-catalyst-be/utils/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dsn *string

func init() {
	host := config.GetEnv("DB_HOST", "localhost")
	port := config.GetEnv("DB_PORT", "5432")
	username := config.GetEnv("DB_USERNAME", "postgres")
	password := config.GetEnv("DB_PASSWORD", "")
	database := config.GetEnv("DB_NAME", "senkou_catalyst")

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

	DB = db
}

// Wire provider function
func GetDB() *gorm.DB {
	if DB == nil {
		ConnectDB()
	}
	return DB
}
