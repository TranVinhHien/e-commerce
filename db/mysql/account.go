package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	db "new-project/db/sqlc"
	services "new-project/services/entity"
	"strings"
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
func (s *SQLStore) GetCustomer(ctx context.Context, customer_id string) (services.Customers, error) {
	u, err := s.Queries.GetCustomer(ctx, customer_id)
	if err != nil {
		log.Fatal("error when get user by customer_id: ", customer_id, err)
		return services.Customers{}, err
	}
	return u.Convert(), nil
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

func (s *SQLStore) CreateCustomerAddresses(ctx context.Context, customer_id string, address *services.CustomerAddress) (err error) {
	return s.Queries.CreateCustomerAddress(ctx, db.CreateCustomerAddressParams{
		Address:     address.Address,
		PhoneNumber: address.PhoneNumber,
		IDAddress:   address.IDAddress,
		CustomerID:  address.CustomerID,
	})
}
func (s *SQLStore) UpdateCustomerAddresses(ctx context.Context, customer_id string, address *services.CustomerAddress) (err error) {
	return s.Queries.UpdateCustomerAddress(ctx, db.UpdateCustomerAddressParams{
		Address:     sql.NullString{String: address.Address, Valid: address.Address != ""},
		PhoneNumber: sql.NullString{String: address.PhoneNumber, Valid: address.PhoneNumber != ""},
		IDAddress:   address.IDAddress,
	})

}
func (s *SQLStore) DeleteCustomerAddresses(ctx context.Context, customer_id string, address_id string) (err error) {
	return s.Queries.DeleteCustomerAddress(ctx, address_id)
}
func (s *SQLStore) ListCustomerAddresses(ctx context.Context, customer_id string) (addresss []services.CustomerAddress, err error) {
	list, err := s.Queries.ListCustomerAddresses(ctx, customer_id)
	// Cấp phát dung lượng slice kết quả
	items := make([]services.CustomerAddress, len(list))
	if err != nil {
		return nil, err
	}
	// Duyệt và chuyển đổi
	for i, item := range list {
		items[i] = item.Convert()
	}
	return items, nil
}
func (s *SQLStore) CustomerAddresses(ctx context.Context, customer_id, address_id string) (info services.Customers, address services.CustomerAddress, err error) {
	item, err := s.Queries.GetCustomerAddressByAddressAndCustomer(ctx, db.GetCustomerAddressByAddressAndCustomerParams{
		IDAddress:  address_id,
		CustomerID: customer_id,
	})
	if err != nil {
		return services.Customers{}, services.CustomerAddress{}, err
	}
	infoDB := db.Customers{
		CustomerID: item.CustomerID,
		Name:       item.Name,
		Email:      item.Name,
		Image:      item.Image,
		Dob:        item.Dob,
		Gender:     item.Gender,
		AccountID:  item.Name,
		CreateDate: item.CreateDate_2,
		UpdateDate: item.UpdateDate_2,
	}
	addressDB := db.CustomerAddress{
		IDAddress:   item.IDAddress,
		CustomerID:  item.CustomerID_2,
		Address:     item.Address,
		PhoneNumber: item.PhoneNumber,
		CreateDate:  item.CreateDate,
		UpdateDate:  item.UpdateDate,
	}
	return infoDB.Convert(), addressDB.Convert(), err
}

func (s *SQLStore) buildGetCustomerAddress(ctx context.Context, ids []string) ([]db.CustomerAddress, error) {

	const querySQL = `-- name: GetAddressID :many
	SELECT id_address, customer_id,address, phone_number, create_date, update_date FROM customer_address
	WHERE id_address IN (%s)
	`

	placeholders := make([]string, len(ids))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf(querySQL, strings.Join(placeholders, ","))
	fmt.Println("query", query)
	// Chuyển đổi danh sách thành các tham số cho câu lệnh SQL
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := s.connPool.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []db.CustomerAddress
	for rows.Next() {
		var i db.CustomerAddress
		if err := rows.Scan(
			&i.IDAddress,
			&i.CustomerID,
			&i.Address,
			&i.PhoneNumber,
			&i.CreateDate,
			&i.UpdateDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *SQLStore) GetCustomerAddresss(ctx context.Context, address_id []string) (is []services.CustomerAddress, err error) {
	items, err := s.buildGetCustomerAddress(ctx, address_id)
	if err != nil {
		log.Fatal("error when get GetCustomerAddresss: ", err)
		return nil, err
	}
	is = make([]services.CustomerAddress, len(items))

	// Duyệt và chuyển đổi
	for i, item := range items {
		is[i] = item.Convert()
	}
	return
}

func (s *SQLStore) buildGetCustomersByID(ctx context.Context, ids []string) ([]db.Customers, error) {

	const querySQL = `-- name: GetCustomersByID :many
	SELECT customer_id, name,email, image, dob, gender,account_id,create_date,update_date FROM customers
	WHERE customer_id IN (%s)
	`

	placeholders := make([]string, len(ids))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf(querySQL, strings.Join(placeholders, ","))
	// Chuyển đổi danh sách thành các tham số cho câu lệnh SQL
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := s.connPool.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []db.Customers
	for rows.Next() {
		var i db.Customers
		if err := rows.Scan(
			&i.CustomerID,
			&i.Name,
			&i.Email,
			&i.Image,
			&i.Dob,
			&i.Gender,
			&i.AccountID,
			&i.CreateDate,
			&i.UpdateDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
