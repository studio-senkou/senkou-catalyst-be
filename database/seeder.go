package database

import (
	"fmt"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	adminPassword, err := bcrypt.GenerateFromPassword(
		fmt.Appendf(nil, "%s", utils.GetEnv("ADMIN_PASSWORD", "")), bcrypt.DefaultCost,
	)

	if err != nil {
		panic("Failed to hash admin password: " + err.Error())
	}

	db.FirstOrCreate(&models.User{
		Name:     utils.GetEnv("ADMIN_NAME", "Catalyst Admin"),
		Email:    utils.GetEnv("ADMIN_EMAIL", "studio.senkou@example.com"),
		Password: adminPassword,
		Role:     "admin",
	})
}
