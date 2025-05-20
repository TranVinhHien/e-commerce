package services

import (
	"context"
	services "new-project/services/entity"
	iservices "new-project/services/interface"
	"time"
)

type ServicesRepository interface {
	iservices.UserRepository
	iservices.ProductRepository
	iservices.CategoriesRepository
	iservices.ĐiscountRepository
	iservices.PaymentRepository
	iservices.OrderRepository
	iservices.RatingRepository
	iservices.ProductsRepository
}
type ServiceUseCase interface {
	iservices.UserUseCase
	iservices.Media
	iservices.Categories
	iservices.Discounts
	iservices.Order
	iservices.Payments
	iservices.RaTings
	iservices.Products
}

type ServicesRedis interface {
	StartExpirationListenerOrderOnline(func(ctx context.Context, orderID string))

	CheckExistsFromBlackList(ctx context.Context, token string, exprid float64) bool
	RemoveTokenExp(zsetKey string)
	AddTokenToBlackList(ctx context.Context, token string, exprid float64) error
	// category
	AddCategories(ctx context.Context, cates []services.Categorys) error
	RemoveCategories(ctx context.Context) error
	GetCategoryTree(ctx context.Context, rootID string) ([]services.Categorys, error)

	// orderOnline
	AddOrderOnline(ctx context.Context, user_id string, payload services.CombinedDataPayLoadMoMo, duration time.Duration) error
	GetOrderOnline(ctx context.Context, user_id string) (payload *services.CombinedDataPayLoadMoMo, err error)
	DeleteOrderOnline(ctx context.Context, orderID string) error
}
