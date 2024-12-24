package services

import (
	"context"
	assets_services "new-project/services/assets"
	services "new-project/services/entity"
)

type UserUseCase interface {
	GetInfo(ctx context.Context, userName string) (*services.Users, *assets_services.ServiceError)
	UpdatePassword(ctx context.Context, userName, oldPassword, newPassword string) (*services.Users, *assets_services.ServiceError)
	Login(ctx context.Context, userName, password string) (accessToken, refreshToken string, info *services.Users, err *assets_services.ServiceError)
	Logout(ctx context.Context, refreshToken string) *assets_services.ServiceError
	Register(ctx context.Context, userName, password, fullName string) *assets_services.ServiceError
	NewAccessToken(ctx context.Context, refreshToken string) (token *string, err *assets_services.ServiceError)
}

type UserRepository interface {
	GetUser(ctx context.Context, userName string) (services.Users, error)
	InsertUser(ctx context.Context, user services.Users) (services.Users, error)
	UpdateUser(ctx context.Context, user services.Users) (services.Users, error)
}
type TempRepository interface {
	// GetUser(ctx context.Context, orgText, source, dest string) (services.Users, error)
	// InsertUser(ctx context.Context, user services.Users) (services.Users, error)
	// UpdateUser(ctx context.Context, user services.Users) (services.Users, error)
}
type ServicesRepository interface {
	UserRepository
	TempRepository
}

type ServicesRedis interface {
	CheckExistsFromBlackList(ctx context.Context, token string, exprid float64) bool
	RemoveTokenExp(zsetKey string)
	AddTokenToBlackList(ctx context.Context, token string, exprid float64) error
}
