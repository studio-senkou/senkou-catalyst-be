package container

import (
	"senkou-catalyst-be/app/controllers"
	"senkou-catalyst-be/app/services"
)

type Container struct {
	UserController               *controllers.UserController
	MerchantController           *controllers.MerchantController
	ProductController            *controllers.ProductController
	CategoryController           *controllers.CategoryController
	PredefinedCategoryController *controllers.PredefinedCategoryController
	AuthController               *controllers.AuthController
	SubscriptionController       *controllers.SubscriptionController
	PaymentMethodsController     *controllers.PaymentMethodsController
	PaymentController            *controllers.PaymentController
	StorageController            *controllers.StorageController
	UserService                  services.UserService
	ProductService               services.ProductService
}
