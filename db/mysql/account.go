package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	db "new-project/db/sqlc"
	services "new-project/services/entity"
	"time"

	"github.com/google/uuid"
)

// CreateAccount - Truyền toàn bộ services.Accounts
func (s *SQLStore) CreateAccount(ctx context.Context, account services.Accounts) (services.Accounts, error) {
	err := s.Queries.CreateAccount(ctx, db.CreateAccountParams{
		AccountID: account.AccountID,
		Username:  account.Username,
		Password:  account.Password,
	})
	if err != nil {
		return services.Accounts{}, err
	}
	return account, nil
}

// DeleteAccount - Chỉ truyền accountID
func (s *SQLStore) DeleteAccount(ctx context.Context, accountID string) error {
	err := s.Queries.DeleteAccount(ctx, accountID)
	return err
}

// UpdateAccount - Truyền toàn bộ services.Accounts
func (s *SQLStore) UpdateAccount(ctx context.Context, account services.Accounts) error {
	err := s.Queries.UpdateAccount(ctx, db.UpdateAccountParams{
		Username: account.Username,
		Password: sql.NullString{String: account.Password, Valid: account.Password != ""},
		ActiveStatus: db.NullAccountsActiveStatus{
			AccountsActiveStatus: db.AccountsActiveStatus(account.ActiveStatus),
			Valid:                account.ActiveStatus != "",
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// GetAccount - Chỉ truyền accountID
func (s *SQLStore) GetAccountByUserName(ctx context.Context, accountID string) (services.Accounts, error) {
	acc, err := s.Queries.GetAccountByUsername(ctx, accountID)
	if err != nil {
		return services.Accounts{}, err
	}
	return acc.Convert(), nil
}

// ListAccounts - Không cần tham số
func (s *SQLStore) ListAccounts(ctx context.Context) ([]services.Accounts, error) {
	accounts, err := s.Queries.ListAccounts(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]services.Accounts, len(accounts))
	for i, acc := range accounts {
		result[i] = acc.Convert()
	}
	return result, nil
}

// ListAccountsPaged - Truyền limit và offset
func (s *SQLStore) ListAccountsPaged(ctx context.Context, limit, offset int32) ([]services.Accounts, error) {
	accounts, err := s.Queries.ListAccountsPaged(ctx, db.ListAccountsPagedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	result := make([]services.Accounts, len(accounts))
	for i, acc := range accounts {
		result[i] = acc.Convert()
	}
	return result, nil
}
func (s *SQLStore) GetCustomerByAccountID(ctx context.Context, accountID string) (services.Customers, error) {
	u, err := s.Queries.GetCustomerByAccountID(ctx, accountID)
	if err != nil {
		log.Fatal("error when get user by accountID: ", accountID, err)
		return services.Customers{}, err
	}
	return u.Convert(), nil
}

func (s *SQLStore) InsertCustomers(ctx context.Context, user services.Customers) error {

	return nil
}

func (s *SQLStore) UpdateCustomers(ctx context.Context, user services.Customers, fn func() error) error {
	hihi := db.UpdateCustomerParams{
		Name:   sql.NullString{String: user.Name, Valid: user.Name != ""},
		Email:  sql.NullString{String: user.Name, Valid: user.Name != ""},
		Image:  sql.NullString{String: user.Image.Data, Valid: user.Image.Valid},
		Dob:    sql.NullTime{Time: user.Dob, Valid: user.Dob != time.Time{}},
		Gender: db.NullCustomersGender{CustomersGender: db.CustomersGender(user.Gender), Valid: user.Gender != ""},

		CustomerID: user.CustomerID,
	}

	// wrap of tran
	err := s.execTS(ctx, func(q *db.Queries) error {
		err := s.Queries.UpdateCustomer(ctx, hihi)
		if err != nil {
			return err
		}
		// // write upload file here
		return fn()
	})

	return err
}

func (s *SQLStore) InsertAccount(ctx context.Context, account services.Accounts) error {
	// u, err := s.Queries.GetUser(ctx, userName)
	// if err != nil {
	return nil
	// }
	// return u.Convert(), nil
}

func (s *SQLStore) Register(ctx context.Context, account *services.Accounts, userInfo *services.Customers, roleID string) error {
	// check accout chung user name email
	_, err := s.Queries.GetAccount(ctx, account.Username)
	if err == nil {
		return errors.New("account already exists")
	}
	err = s.execTS(ctx, func(q *db.Queries) error {
		err := q.CreateAccount(ctx, db.CreateAccountParams{
			AccountID: account.AccountID,
			Username:  account.Username,
			Password:  account.Password,
		})
		if err != nil {
			return err
		}
		err = q.CreateCustomer(ctx, db.CreateCustomerParams{
			CustomerID: userInfo.CustomerID,
			Name:       userInfo.Name,
			Email:      userInfo.Email,
			Dob:        sql.NullTime{Time: userInfo.Dob, Valid: userInfo.Dob != time.Time{}},
			Gender:     db.NullCustomersGender{CustomersGender: db.CustomersGender(userInfo.Gender), Valid: userInfo.Gender != ""},
			AccountID:  account.AccountID,
		})
		if err != nil {
			return err
		}

		err = q.CreateRoleAccount(ctx, db.CreateRoleAccountParams{
			RoleAccountID: uuid.New().String(),
			AccountID:     account.AccountID,
			RoleID:        roleID, // set default role
		})
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (s *SQLStore) Login(ctx context.Context, username string) (account services.Accounts, role string, err error) {
	acc, err := s.Queries.Login(ctx, username)
	if err != nil {
		return services.Accounts{}, "", err
	}
	account1 := db.Accounts{
		AccountID:    acc.AccountID,
		Username:     acc.Username,
		Password:     acc.Password,
		ActiveStatus: acc.ActiveStatus,
		CreateDate:   acc.CreateDate,
		UpdateDate:   acc.UpdateDate,
	}
	account = account1.Convert()
	return account, acc.RoleID, nil
}

//
