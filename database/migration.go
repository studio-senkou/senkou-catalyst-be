package database

import (
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.UserHasToken{},
		&models.Merchant{},
		&models.Category{},
		&models.Product{},
		&models.Subscription{},
		&models.SubscriptionPlan{},
		&models.UserSubscription{},
		&models.PredefinedCategory{},
	)

	if err != nil {
		panic("Migration failed: " + err.Error())
	}
}
