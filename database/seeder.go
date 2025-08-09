package database

import (
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/utils/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	adminPasswordStr := config.GetEnv("ADMIN_PASSWORD", "admin123")
	adminPassword, err := bcrypt.GenerateFromPassword(
		[]byte(adminPasswordStr), bcrypt.DefaultCost,
	)

	if err != nil {
		panic("Failed to hash admin password: " + err.Error())
	}

	administrator := new(models.User)
	db.Where("email = ?", config.GetEnv("ADMIN_EMAIL", "studio.senkou@example.com")).FirstOrCreate(&administrator, models.User{
		Name:     config.GetEnv("ADMIN_NAME", "Catalyst Admin"),
		Email:    config.GetEnv("ADMIN_EMAIL", "studio.senkou@example.com"),
		Phone:    config.GetEnv("ADMIN_PHONE", "1234567890"),
		Password: adminPassword,
		Role:     "admin",
	})

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to hash password: " + err.Error())
	}

	firstUser := new(models.User)
	db.Where("email = ?", "agus.prasetyo@senkou.co.id").FirstOrCreate(&firstUser, models.User{
		Name:     "Agus Prasetyo",
		Email:    "agus.prasetyo@senkou.co.id",
		Phone:    "6281234567890",
		Password: hashedPassword,
		Role:     "user",
	})

	secondUser := new(models.User)
	db.Where("email = ?", "budi.santoso@senkou.co.id").FirstOrCreate(&secondUser, models.User{
		Name:     "Budi Santoso",
		Email:    "budi.santoso@senkou.co.id",
		Phone:    "6289876543210",
		Password: hashedPassword,
		Role:     "user",
	})
}
