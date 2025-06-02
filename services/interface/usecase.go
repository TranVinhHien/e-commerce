package iservices

import (
	"context"
	"mime/multipart"
	assets_services "new-project/services/assets"
	services "new-project/services/entity"
)

type UserUseCase interface {
	// without login
	//GetInfo(ctx context.Context, customer_id string) (map[string]interface{}, *assets_services.ServiceError)
	UpdatePassword(ctx context.Context, customer_id, oldPassword, newPassword string) *assets_services.ServiceError
	Login(ctx context.Context, username, password, token string) (accessToken, refreshToken string, info map[string]interface{}, err *assets_services.ServiceError)
	Logout(ctx context.Context, refreshToken string) *assets_services.ServiceError
	Register(ctx context.Context, customer_id, password string, userInfo *services.Customers) *assets_services.ServiceError
	NewAccessToken(ctx context.Context, refreshToken string) (token *string, err *assets_services.ServiceError)

	// with login
	UpdadateInfo(ctx context.Context, customer_id string, info *services.Customers) *assets_services.ServiceError
	UpdadateAvatar(ctx context.Context, customer_id string, file *multipart.FileHeader) (err *assets_services.ServiceError)
	CreateCustomerAddress(ctx context.Context, customer_id string, address_info *services.CustomerAddress) (err *assets_services.ServiceError)
	UpdateCustomerAddress(ctx context.Context, customer_id string, address_info *services.CustomerAddress) (err *assets_services.ServiceError)
	InfoUser(ctx context.Context, customer_id string) (info map[string]interface{}, err *assets_services.ServiceError)
	ListAddress(ctx context.Context, customer_id string) (info map[string]interface{}, err *assets_services.ServiceError)
}

type Media interface {
	RenderImage(ctx context.Context, filename string) string
	RenderProductImages(ctx context.Context, filename string) string
}
type Categories interface {
	GetCategoris(ctx context.Context, userName string) (map[string]interface{}, *assets_services.ServiceError)
}
type Discounts interface {
	ListDiscount(ctx context.Context, query services.QueryFilter) (map[string]interface{}, *assets_services.ServiceError)
}
type Payments interface {
	ListPayment(ctx context.Context) (map[string]interface{}, *assets_services.ServiceError)
}
type Order interface {
	ListOrderByUserID(ctx context.Context, user_id string, query services.QueryFilter) (map[string]interface{}, *assets_services.ServiceError)
	CallBackMoMo(ctx context.Context, tran services.TransactionMoMO)
	CreateOrder(ctx context.Context, user_id string, order *services.CreateOrderParams) (payment_url map[string]interface{}, err *assets_services.ServiceError)
	GetURLOrderMoMOAgain(ctx context.Context, user_id string) (map[string]interface{}, *assets_services.ServiceError)
	RemoveOrderOnline(ctx context.Context, orderIDs string)
	CancelOrder(ctx context.Context, userID, order_id string) *assets_services.ServiceError
}
type RaTings interface {
	ListRating(ctx context.Context, products_spu_id string, query services.QueryFilter) (map[string]interface{}, *assets_services.ServiceError)
	CreateRating(ctx context.Context, userID string, rating services.Ratings) *assets_services.ServiceError
}

type Products interface {
	GetAllProductSimple(ctx context.Context, query services.QueryFilter) (map[string]interface{}, *assets_services.ServiceError)
	GetDetailProduct(ctx context.Context, productSpuID string) (map[string]interface{}, *assets_services.ServiceError)
}
