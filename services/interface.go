package services

import (
	"context"
	iservices "new-project/services/interface"
)

type ServicesRepository interface {
	iservices.UserRepository
	iservices.HandleRepository
}
type ServiceUseCase interface {
	iservices.UserUseCase
	iservices.Media
}

type ServicesRedis interface {
	CheckExistsFromBlackList(ctx context.Context, token string, exprid float64) bool
	RemoveTokenExp(zsetKey string)
	AddTokenToBlackList(ctx context.Context, token string, exprid float64) error
}
