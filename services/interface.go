package services

import (
	"context"
	services "new-project/services/entity"
	iservices "new-project/services/interface"
)

type ServicesRepository interface {
	iservices.UserRepository
	iservices.HandleRepository
	iservices.CategoriesRepository
}
type ServiceUseCase interface {
	iservices.UserUseCase
	iservices.Media
	iservices.Categories
}

type ServicesRedis interface {
	CheckExistsFromBlackList(ctx context.Context, token string, exprid float64) bool
	RemoveTokenExp(zsetKey string)
	AddTokenToBlackList(ctx context.Context, token string, exprid float64) error

	AddCategories(ctx context.Context, cates []services.Categorys) error
	RemoveCategories(ctx context.Context) error
	GetCategoryTree(ctx context.Context, rootID string) ([]services.Categorys, error)
}
