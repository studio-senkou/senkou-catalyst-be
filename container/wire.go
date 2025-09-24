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
	mailerUtil "senkou-catalyst-be/utils/mailer"
	"senkou-catalyst-be/utils/queue"

	"github.com/google/wire"
)

var DatabaseSet = wire.NewSet(
	config.GetDB,
)

var RepositorySet = wire.NewSet(
	repositories.NewUserRepository,
	repositories.NewMerchantRepository,
	repositories.NewEmailActivationRepository,
	repositories.NewProductRepository,
	repositories.NewProductInteractionRepository,
	repositories.NewCategoryRepository,
	repositories.NewPredefinedCategoryRepository,
	repositories.NewAuthRepository,
	repositories.NewOAuthRepository,
	repositories.NewSubscriptionRepository,
	repositories.NewSubscriptionPlanRepository,
	repositories.NewSubscriptionOrderRepository,
	repositories.NewPaymentTransactionRepository,
)

var ServiceSet = wire.NewSet(
	services.NewUserService,
	services.NewMerchantService,
	services.NewProductService,
	services.NewProductInteractionService,
	services.NewCategoryService,
	services.NewPredefinedCategoryService,
	services.NewAuthService,
	services.NewSubscriptionService,
	services.NewSubscriptionOrderService,
	services.NewPaymentMethodsService,
	services.NewPaymentService,
	mailerUtil.NewMailerService,
)

var ControllerSet = wire.NewSet(
	controllers.NewUserController,
	controllers.NewMerchantController,
	controllers.NewProductController,
	controllers.NewCategoryController,
	controllers.NewPredefinedCategoryController,
	controllers.NewAuthController,
	controllers.NewOAuthController,
	controllers.NewSubscriptionController,
	controllers.NewPaymentMethodsController,
	controllers.NewPaymentController,
	controllers.NewStorageController,
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

func ProvideQueueService() (*queue.QueueService, error) {
	cfg := queue.DefaultQueueConfig()
	return queue.NewQueueService(cfg)
}

var QueueSet = wire.NewSet(
	ProvideQueueService,
)

func InitializeUserController() (*controllers.UserController, error) {
	wire.Build(
		DatabaseSet,
		RepositorySet,
		ServiceSet,
		ControllerSet,
		QueueSet,
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
		QueueSet,
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
		QueueSet,
	)
	return nil, nil
}

func InitializeOAuthController() (*controllers.OAuthController, error) {
	wire.Build(
		DatabaseSet,
		ControllerSet,
		ServiceSet,
		RepositorySet,
		UtilSet,
		QueueSet,
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
		QueueSet,
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
		QueueSet,
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

func InitializeProductInteractionService() (services.ProductInteractionService, func(), error) {
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
		QueueSet,
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
	oauthController *controllers.OAuthController,
	subscriptionController *controllers.SubscriptionController,
	paymentMethodsController *controllers.PaymentMethodsController,
	paymentController *controllers.PaymentController,
	storageController *controllers.StorageController,
	userService services.UserService,
	productService services.ProductService,
	queueService *queue.QueueService,
) *Container {
	return &Container{
		UserController:               userController,
		MerchantController:           merchantController,
		ProductController:            productController,
		CategoryController:           categoryController,
		PredefinedCategoryController: predefinedCategoryController,
		AuthController:               authController,
		OAuthController:              oauthController,
		SubscriptionController:       subscriptionController,
		PaymentMethodsController:     paymentMethodsController,
		PaymentController:            paymentController,
		StorageController:            storageController,
		UserService:                  userService,
		ProductService:               productService,
		QueueService:                 queueService,
	}
}
