package migration

import (
	"senkou-catalyst-be/app/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.UserHasToken{},
		&models.Merchant{},
		&models.Category{},
		&models.Product{},
		&models.ProductMetric{},
		&models.Subscription{},
		&models.SubscriptionPlan{},
		&models.PredefinedCategory{},
	)

	if err != nil {
		panic("Migration failed: " + err.Error())
	}

	err = db.AutoMigrate(
		&models.PaymentTransaction{},
	)

	if err != nil {
		panic("Migration failed: " + err.Error())
	}

	err = db.AutoMigrate(
		&models.UserSubscription{},
		&models.SubscriptionOrder{},
	)

	if err != nil {
		panic("Migration failed: " + err.Error())
	}

	createIndexes(db)
}

func createIndexes(db *gorm.DB) {
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payment_transactions_transaction_id ON payment_transactions(transaction_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payment_transactions_status ON payment_transactions(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_payment_transactions_created_at ON payment_transactions(created_at)")

	db.Exec("CREATE INDEX IF NOT EXISTS idx_subscription_orders_user_id ON subscription_orders(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_subscription_orders_subscription_id ON subscription_orders(subscription_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_subscription_orders_payment_transaction_id ON subscription_orders(payment_transaction_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_subscription_orders_status ON subscription_orders(status)")
}
