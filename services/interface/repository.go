package iservices

import (
	"context"
	services "new-project/services/entity"
)

type UserRepository interface {
	GetCustomer(ctx context.Context, customer_id string) (services.Customers, error)
	GetCustomerByAccountID(ctx context.Context, accoutn_id string) (services.Customers, error)

	UpdateCustomers(ctx context.Context, user services.Customers, fn func() error) error

	GetAccountByUserName(ctx context.Context, username string) (services.Accounts, error)
	UpdateAccount(ctx context.Context, account services.Accounts) error

	Register(ctx context.Context, account *services.Accounts, userInfo *services.Customers, roleID string) error
	Login(ctx context.Context, username string) (account services.Accounts, role string, err error)

	ListCustomerAddresses(ctx context.Context, customer_id string) (addresss []services.CustomerAddress, err error)
	CreateCustomerAddresses(ctx context.Context, customer_id string, address *services.CustomerAddress) (err error)
	UpdateCustomerAddresses(ctx context.Context, customer_id string, address *services.CustomerAddress) (err error)
	DeleteCustomerAddresses(ctx context.Context, customer_id string, address_id string) (err error)
}
type CategoriesRepository interface {
	ListCategories(ctx context.Context) ([]services.Categorys, error)
	ListCategoriesByID(ctx context.Context, cate_id string) ([]services.Categorys, error)
}
type HandleRepository interface {
	// ExecTran(ctx context.Context, fn func() error) error
	// GetUser(ctx context.Context, orgText, source, dest string) (services.Customers, error)
	// InsertUser(ctx context.Context, user services.Customers) (services.Customers, error)
	// UpdateUser(ctx context.Context, user services.Customers) (services.Customers, error)
}
