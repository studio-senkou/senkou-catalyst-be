package seeder

import (
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/utils/config"
	"time"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {

	userRepository := repositories.NewUserRepository(db)
	merchantRepository := repositories.NewMerchantRepository(db)
	emailActivationRepo := repositories.NewEmailActivationRepository(db)

	userService := services.NewUserService(userRepository, merchantRepository, emailActivationRepo)

	adminPasswordStr := config.GetEnv("SEEDER_ADMIN_PASSWORD", "admin123")

	activeNow := time.Now()

	userService.Create(&models.User{
		Name:            config.GetEnv("SEEDER_ADMIN_NAME", "Catalyst Admin"),
		Email:           config.GetEnv("SEEDER_ADMIN_EMAIL", "studio.senkou@example.com"),
		Phone:           config.GetEnv("SEEDER_ADMIN_PHONE", "1234567890"),
		Password:        []byte(adminPasswordStr),
		Role:            "admin",
		EmailVerifiedAt: &activeNow,
	}, nil)

	userService.Create(&models.User{
		Name:            "Agus Prasetyo",
		Email:           "agus.prasetyo@senkou.co.id",
		Phone:           "6281234567890",
		Password:        []byte("password"),
		Role:            "user",
		EmailVerifiedAt: &activeNow,
	}, &models.Merchant{
		Name:     "Agus's Store",
		Username: "agus-prasetyo",
	})

	userService.Create(&models.User{
		Name:            "Budi Santoso",
		Email:           "budi.santoso@senkou.co.id",
		Phone:           "6289876543210",
		Password:        []byte("password"),
		Role:            "user",
		EmailVerifiedAt: &activeNow,
	}, &models.Merchant{
		Name:     "Budi's Store",
		Username: "budi-santoso",
	})

	return nil
}
