package seeder

import (
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/utils/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {

	userRepository := repositories.NewUserRepository(db)
	merchantRepository := repositories.NewMerchantRepository(db)
	userService := services.NewUserService(userRepository, merchantRepository)

	adminPasswordStr := config.GetEnv("SEEDER_ADMIN_PASSWORD", "admin123")
	adminPassword, err := bcrypt.GenerateFromPassword(
		[]byte(adminPasswordStr), bcrypt.DefaultCost,
	)

	userService.Create(&models.User{
		Name:     config.GetEnv("SEEDER_ADMIN_NAME", "Catalyst Admin"),
		Email:    config.GetEnv("SEEDER_ADMIN_EMAIL", "studio.senkou@example.com"),
		Phone:    config.GetEnv("SEEDER_ADMIN_PHONE", "1234567890"),
		Password: adminPassword,
		Role:     "admin",
	}, nil)

	if err != nil {
		panic("Failed to hash admin password: " + err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to hash password: " + err.Error())
	}

	userService.Create(&models.User{
		Name:     "Agus Prasetyo",
		Email:    "agus.prasetyo@senkou.co.id",
		Phone:    "6281234567890",
		Password: hashedPassword,
		Role:     "user",
	}, &models.Merchant{
		Name:     "Agus's Store",
		Username: "agus-prasetyo",
	})

	userService.Create(&models.User{
		Name:     "Budi Santoso",
		Email:    "budi.santoso@senkou.co.id",
		Phone:    "6289876543210",
		Password: hashedPassword,
		Role:     "user",
	}, &models.Merchant{
		Name:     "Budi's Store",
		Username: "budi-santoso",
	})

	return nil
}
