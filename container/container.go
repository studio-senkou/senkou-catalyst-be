package container

import (
	"log"
	"senkou-catalyst-be/app/controllers"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/utils/queue"
)

type Container struct {
	UserController               *controllers.UserController
	MerchantController           *controllers.MerchantController
	ProductController            *controllers.ProductController
	CategoryController           *controllers.CategoryController
	PredefinedCategoryController *controllers.PredefinedCategoryController
	AuthController               *controllers.AuthController
	OAuthController              *controllers.OAuthController
	SubscriptionController       *controllers.SubscriptionController
	PaymentMethodsController     *controllers.PaymentMethodsController
	PaymentController            *controllers.PaymentController
	StorageController            *controllers.StorageController
	UserService                  services.UserService
	ProductService               services.ProductService
	QueueService                 *queue.QueueService
}

func (c *Container) StartQueueService() {
	if c.QueueService != nil {

		// Register handlers
		c.QueueService.RegisterEmailHandlers()

		go func() {
			if err := c.QueueService.Start(); err != nil {
				log.Printf("Queue service failed: %v", err)
			}
		}()

		log.Println("Queue service started with email handlers")
	}
}
