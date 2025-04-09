package iservices

import (
	"context"
	services "new-project/services/entity"
)

type UserRepository interface {
	GetCustomerByAccountID(ctx context.Context, idUser string) (services.Customers, error)
	UpdateCustomers(ctx context.Context, user services.Customers, fn func() error) error

	GetAccountByUserName(ctx context.Context, username string) (services.Accounts, error)
	UpdateAccount(ctx context.Context, account services.Accounts) error

	Register(ctx context.Context, account *services.Accounts, userInfo *services.Customers, roleID string) error
	Login(ctx context.Context, username string) (account services.Accounts, role string, err error)
}
type HandleRepository interface {
	// ExecTran(ctx context.Context, fn func() error) error
	// GetUser(ctx context.Context, orgText, source, dest string) (services.Customers, error)
	// InsertUser(ctx context.Context, user services.Customers) (services.Customers, error)
	// UpdateUser(ctx context.Context, user services.Customers) (services.Customers, error)
}
