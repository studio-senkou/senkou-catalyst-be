//go:build wireinject
// +build wireinject

package container

import (
	"senkou-catalyst-be/app/controllers"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/integrations/midtrans"
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
	repositories.NewSubscriptionOrderRepository,
	repositories.NewPaymentTransactionRepository,
)

var ServiceSet = wire.NewSet(
	services.NewUserService,
	services.NewMerchantService,
	services.NewProductService,
	services.NewCategoryService,
	services.NewPredefinedCategoryService,
	services.NewAuthService,
	services.NewSubscriptionService,
	services.NewSubscriptionOrderService,
	services.NewPaymentMethodsService,
	services.NewPaymentService,
)

var ControllerSet = wire.NewSet(
	controllers.NewUserController,
	controllers.NewMerchantController,
	controllers.NewProductController,
	controllers.NewCategoryController,
	controllers.NewPredefinedCategoryController,
	controllers.NewAuthController,
	controllers.NewSubscriptionController,
	controllers.NewPaymentMethodsController,
	controllers.NewPaymentController,
)

func ProvideJWTManager() (*authUtil.JWTManager, error) {
	secret := configUtil.MustGetEnv("AUTH_SECRET")
	return authUtil.NewJWTManager(secret)
}

var UtilSet = wire.NewSet(
	ProvideJWTManager,
)

func ProvideMidtransClient() (*midtrans.MidtransClient, error) {
	return midtrans.NewMidtransClient(), nil
}

func ProvideMidtransBuilder(client *midtrans.MidtransClient) *midtrans.PaymentBuilder {
	return midtrans.NewPaymentBuilder(client)
}

var MidtransSet = wire.NewSet(
	ProvideMidtransClient,
	ProvideMidtransBuilder,
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
		MidtransSet,
	)
	return nil, nil
}

func InitializePaymentController() (*controllers.PaymentController, error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
		MidtransSet,
	)
	return nil, nil
}

func InitializePaymentMethodsController() (*controllers.PaymentMethodsController, error) {
	wire.Build(
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

func InitializeSubscriptionOrderService() (services.SubscriptionOrderService, func(), error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
	)
	return nil, nil, nil
}

func InitializePaymentMethodsService() (services.PaymentMethodsService, func(), error) {
	wire.Build(
		ServiceSet,
	)
	return nil, nil, nil
}

func InitializePaymentService() (services.PaymentService, func(), error) {
	wire.Build(
		DatabaseSet,
		MidtransSet,
		ServiceSet,
		RepositorySet,
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
		MidtransSet,
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
	paymentMethodsController *controllers.PaymentMethodsController,
	paymentController *controllers.PaymentController,
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
		PaymentMethodsController:     paymentMethodsController,
		PaymentController:            paymentController,
		UserService:                  userService,
		ProductService:               productService,
	}
}
