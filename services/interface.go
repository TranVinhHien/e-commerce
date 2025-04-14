package services

import (
	"context"
	services "new-project/services/entity"
	iservices "new-project/services/interface"
)

type ServicesRepository interface {
	iservices.UserRepository
	iservices.ProductRepository
	iservices.CategoriesRepository
	iservices.ĐiscountRepository
	iservices.PaymentRepository
	iservices.OrderRepository
}
type ServiceUseCase interface {
	iservices.UserUseCase
	iservices.Media
	iservices.Categories
	iservices.Discounts
	iservices.Order
	iservices.Payments
}

type ServicesRedis interface {
	CheckExistsFromBlackList(ctx context.Context, token string, exprid float64) bool
	RemoveTokenExp(zsetKey string)
	AddTokenToBlackList(ctx context.Context, token string, exprid float64) error

	AddCategories(ctx context.Context, cates []services.Categorys) error
	RemoveCategories(ctx context.Context) error
	GetCategoryTree(ctx context.Context, rootID string) ([]services.Categorys, error)
}
