package container

import (
	"senkou-catalyst-be/controllers"
	"senkou-catalyst-be/services"
)

type Container struct {
	UserController               *controllers.UserController
	MerchantController           *controllers.MerchantController
	ProductController            *controllers.ProductController
	CategoryController           *controllers.CategoryController
	PredefinedCategoryController *controllers.PredefinedCategoryController
	AuthController               *controllers.AuthController

	UserService    services.UserService
	ProductService services.ProductService
}
