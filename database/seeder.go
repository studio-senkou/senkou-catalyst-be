package database

import (
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	adminPasswordStr := utils.GetEnv("ADMIN_PASSWORD", "admin123")
	adminPassword, err := bcrypt.GenerateFromPassword(
		[]byte(adminPasswordStr), bcrypt.DefaultCost,
	)

	if err != nil {
		panic("Failed to hash admin password: " + err.Error())
	}

	administrator := new(models.User)
	db.Where("email = ?", utils.GetEnv("ADMIN_EMAIL", "studio.senkou@example.com")).FirstOrCreate(&administrator, models.User{
		Name:     utils.GetEnv("ADMIN_NAME", "Catalyst Admin"),
		Email:    utils.GetEnv("ADMIN_EMAIL", "studio.senkou@example.com"),
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
		Password: hashedPassword,
		Role:     "user",
	})

	secondUser := new(models.User)
	db.Where("email = ?", "budi.santoso@senkou.co.id").FirstOrCreate(&secondUser, models.User{
		Name:     "Budi Santoso",
		Email:    "budi.santoso@senkou.co.id",
		Password: hashedPassword,
		Role:     "user",
	})
}
