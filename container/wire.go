//go:build wireinject
// +build wireinject

package container

import (
	"senkou-catalyst-be/config"
	"senkou-catalyst-be/controllers"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/services"

	"github.com/google/wire"
)

var DatabaseSet = wire.NewSet(
	config.GetDB,
)

var RepositorySet = wire.NewSet(
	repositories.NewUserRepository,
	repositories.NewMerchantRepository,
	repositories.NewProductRepository,
	repositories.NewCategoryRepository,
	repositories.NewPredefinedCategoryRepository,
	repositories.NewAuthRepository,
)

var ServiceSet = wire.NewSet(
	services.NewUserService,
	services.NewMerchantService,
	services.NewProductService,
	services.NewCategoryService,
	services.NewPredefinedCategoryService,
	services.NewAuthService,
)

var ControllerSet = wire.NewSet(
	controllers.NewUserController,
	controllers.NewMerchantController,
	controllers.NewProductController,
	controllers.NewCategoryController,
	controllers.NewPredefinedCategoryController,
	controllers.NewAuthController,
)

// Wire functions for controllers
func InitializeUserController() (*controllers.UserController, error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
	)
	return nil, nil
}

func InitializeMerchantController() (*controllers.MerchantController, error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
	)
	return nil, nil
}

func InitializeProductController() (*controllers.ProductController, error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
	)
	return nil, nil
}

func InitializeCategoryController() (*controllers.CategoryController, error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
	)
	return nil, nil
}

func InitializePredefinedCategoryController() (*controllers.PredefinedCategoryController, error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
	)
	return nil, nil
}

func InitializeAuthController() (*controllers.AuthController, error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
	)
	return nil, nil
}

// Initialize services for middleware
func InitializeUserService() (services.UserService, func(), error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
	)
	return nil, nil, nil
}

func InitializeProductService() (services.ProductService, func(), error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
	)
	return nil, nil, nil
}

// Initialize Container with all dependencies
func InitializeContainer() (*Container, error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
		NewContainer,
	)
	return nil, nil
}

// NewContainer creates a new container with all dependencies
func NewContainer(
	userController *controllers.UserController,
	merchantController *controllers.MerchantController,
	productController *controllers.ProductController,
	categoryController *controllers.CategoryController,
	predefinedCategoryController *controllers.PredefinedCategoryController,
	authController *controllers.AuthController,
	userService services.UserService,
	productService services.ProductService,
) *Container {
	return &Container{
		UserController:               userController,
		MerchantController:           merchantController,
		ProductController:            productController,
		CategoryController:           categoryController,
		PredefinedCategoryController: predefinedCategoryController,
		AuthController:               authController,
		UserService:                  userService,
		ProductService:               productService,
	}
}
