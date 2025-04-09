package iservices

import (
	"context"
	"mime/multipart"
	assets_services "new-project/services/assets"
	services "new-project/services/entity"
)

type UserUseCase interface {
	// without login
	GetInfo(ctx context.Context, userName string) (map[string]interface{}, *assets_services.ServiceError)
	UpdatePassword(ctx context.Context, userName, oldPassword, newPassword string) *assets_services.ServiceError
	Login(ctx context.Context, userName, password string) (accessToken, refreshToken string, info map[string]interface{}, err *assets_services.ServiceError)
	Logout(ctx context.Context, refreshToken string) *assets_services.ServiceError
	Register(ctx context.Context, userName, password string, userInfo *services.Customers) *assets_services.ServiceError
	NewAccessToken(ctx context.Context, refreshToken string) (token *string, err *assets_services.ServiceError)

	// with login
	UpdadateInfo(ctx context.Context, user_id string, info *services.Customers) *assets_services.ServiceError
	UpdadateAvatar(ctx context.Context, user_id string, file *multipart.FileHeader) (err *assets_services.ServiceError)
}

type Media interface {
	RenderImage(ctx context.Context, filename string) string
}
