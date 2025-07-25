//go:build wireinject
// +build wireinject

package container

import (
	"senkou-catalyst-be/app/controllers"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/platform/config"
	"senkou-catalyst-be/repositories"

	authUtil "senkou-catalyst-be/utils/auth"
	configUtil "senkou-catalyst-be/utils/config"

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
	repositories.NewSubscriptionRepository,
	repositories.NewSubscriptionPlanRepository,
)

var ServiceSet = wire.NewSet(
	services.NewUserService,
	services.NewMerchantService,
	services.NewProductService,
	services.NewCategoryService,
	services.NewPredefinedCategoryService,
	services.NewAuthService,
	services.NewSubscriptionService,
)

var ControllerSet = wire.NewSet(
	controllers.NewUserController,
	controllers.NewMerchantController,
	controllers.NewProductController,
	controllers.NewCategoryController,
	controllers.NewPredefinedCategoryController,
	controllers.NewAuthController,
	controllers.NewSubscriptionController,
)

func ProvideJWTManager() (*authUtil.JWTManager, error) {
	secret := configUtil.MustGetEnv("AUTH_SECRET")
	return authUtil.NewJWTManager(secret)
}

var UtilSet = wire.NewSet(
	ProvideJWTManager,
)

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
		UtilSet,
	)
	return nil, nil
}

func InitializeSubscriptionController() (*controllers.SubscriptionController, error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
	)
	return nil, nil
}

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

func InitializeContainer() (*Container, error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
		UtilSet,
		NewContainer,
	)
	return nil, nil
}

func NewContainer(
	userController *controllers.UserController,
	merchantController *controllers.MerchantController,
	productController *controllers.ProductController,
	categoryController *controllers.CategoryController,
	predefinedCategoryController *controllers.PredefinedCategoryController,
	authController *controllers.AuthController,
	subscriptionController *controllers.SubscriptionController,
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
		SubscriptionController:       subscriptionController,
		UserService:                  userService,
		ProductService:               productService,
	}
}
