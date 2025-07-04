package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	environmentError := godotenv.Load(".env")

	if environmentError != nil {
		log.Fatal("Error loading the .env file, please ensure it exists and is correctly formatted.")
	}
}

func GetEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}