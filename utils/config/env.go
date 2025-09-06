// utils/config/env.go
package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	environmentError := godotenv.Load(".env")

	if environmentError != nil {
		log.Fatal("Error loading the .env file, please ensure it exists and is correctly formatted.")
	}
}

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GetEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func GetEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}

func GetEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s is required but not set", key))
	}
	return value
}

func loadEnvForTesting() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Traverse up to 5 levels to find .env file
	for i := 0; i < 5; i++ {
		envPath := filepath.Join(wd, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}
		wd = filepath.Dir(wd)
	}

	return fmt.Errorf(".env file not found")
}
