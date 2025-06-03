package services

import (
	"context"
	"fmt"
	"mime/multipart"
	util_assets "new-project/assets/util"
	assets_services "new-project/services/assets"
	services "new-project/services/entity"

	"github.com/google/uuid"
)

func (s *service) GetInfo(ctx context.Context, userName string) (map[string]interface{}, *assets_services.ServiceError) {
	// kiểm tra thông tin abc xyz'
	user, err := s.repository.GetCustomerByAccountID(ctx, userName)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, assets_services.NewError(400, err)
	}
	result, err := assets_services.HideFields(user, "", "account_id", "customer_id")
	if err != nil {
		fmt.Println("Error:", err)
		return nil, assets_services.NewError(400, err)
	}

	return result, assets_services.NewError(400, err)
}
func (s *service) UpdatePassword(ctx context.Context, userName, oldPassword, newPassword string) *assets_services.ServiceError {
	//
	account, err := s.repository.GetAccountByUserName(ctx, userName)
	if err != nil {
		return assets_services.NewError(400, err)
	}
	err = util_assets.CheckPassword(oldPassword, account.Password)
	if err != nil {
		return assets_services.NewError(400, fmt.Errorf("oldPassword is incorrect"))
	}
	newPassword, err = util_assets.HashPassword(newPassword)
	if err != nil {
		return assets_services.NewError(400, err)
	}
	err = s.repository.UpdateAccount(ctx, services.Accounts{
		Username: userName,
		Password: newPassword,
	})
	if err != nil {
		return assets_services.NewError(400, fmt.Errorf("update failed for user: "+userName))
	}
	return nil
}

func (s *service) Login(ctx context.Context, userName, password, token string) (accessToken, refreshToken string, info map[string]interface{}, error_s *assets_services.ServiceError) {
	account, roleID, err := s.repository.Login(ctx, userName)
	fmt.Println("token:", token)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", "", nil, assets_services.NewError(404, fmt.Errorf("không tìm thấy tài khoảng"))
		}
		return "", "", nil, assets_services.NewError(400, err)
	}
	err = util_assets.CheckPassword(password, account.Password)
	if err != nil {
		return "", "", nil, assets_services.NewError(400, fmt.Errorf("password is incorrect"))
	}

	userID := ""
	result123 := map[string]interface{}{}
	// get info customer
	if roleID == s.env.Customer {
		user, addrs, err := s._getInfoCustomer(ctx, "", account.AccountID)
		if err != nil {
			return "", "", nil, assets_services.NewError(400, err)
		}
		userID = user.CustomerID
		result123, err = assets_services.HideFields(user, "", "account_id", "customer_id")
		if err != nil {
			return "", "", nil, assets_services.NewError(400, err)
		}
		result123["addrs"] = addrs
		// get info emp
	} else {
		user, err := s.repository.GetCustomerByAccountID(ctx, account.AccountID)
		if err != nil {
			return "", "", nil, assets_services.NewError(400, err)
		}
		userID = user.CustomerID
		result123, err = assets_services.HideFields(user, "", "account_id", "customer_id")
		if err != nil {
			return "", "", nil, assets_services.NewError(400, err)
		}
	}

	// update token mobile
	if token != "" {
		dulieucapnhat := services.Customers{
			CustomerID:              userID,
			DeviceRegistrationToken: services.Narg[string]{Data: token, Valid: true},
		}
		s.repository.UpdateCustomers(ctx, dulieucapnhat, func() error { return nil })
	}

	_, accessToken, err = s.jwt.CreateToken(userID, s.env.AccessTokenDuration)
	if err != nil {
		return "", "", nil, assets_services.NewError(400, fmt.Errorf("error create access token: %v", err))
	}
	_, refreshToken, err = s.jwt.CreateToken(userID, s.env.RefershTokenDuration)
	if err != nil {
		return "", "", nil, assets_services.NewError(400, fmt.Errorf("error create refersh token: %v", err))
	}
	return accessToken, refreshToken, result123, nil

}

func (s *service) Logout(ctx context.Context, refreshToken string) *assets_services.ServiceError {
	// add token
	payload, err := s.jwt.VerifyToken(refreshToken)
	if err != nil {
		return assets_services.NewError(400, fmt.Errorf("error when verifyToken: %v", err))
	}
	s.redis.AddTokenToBlackList(ctx, refreshToken, float64(payload.Exp))
	return nil
}
func (s *service) Register(ctx context.Context, userName, password string, userInfo *services.Customers) *assets_services.ServiceError {
	//
	password, err := util_assets.HashPassword(password)
	if err != nil {
		return assets_services.NewError(400, err)
	}

	accountID := uuid.New().String()
	err = s.repository.Register(ctx, &services.Accounts{
		AccountID: accountID,
		Username:  userName,
		Password:  password,
	}, &services.Customers{
		CustomerID: uuid.New().String(),
		AccountID:  accountID,
		Name:       userInfo.Name,
		Email:      userInfo.Email,
		Dob:        userInfo.Dob,
		Gender:     userInfo.Gender,
	}, s.env.Customer)
	if err != nil {
		return assets_services.NewError(400, err)
	}

	return nil
}
func (s *service) NewAccessToken(ctx context.Context, refreshToken string) (*string, *assets_services.ServiceError) {
	payload, err := s.jwt.VerifyToken(refreshToken)
	if err != nil {

		return nil, assets_services.NewError(400, fmt.Errorf("error when verifyToken: %v", err))
	}
	_, accessToken, err := s.jwt.CreateToken(payload.Sub, s.env.AccessTokenDuration)
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("error when create new accessToken: %v", err))
	}
	return &accessToken, nil
}

// // withlogin//////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//	//
//
// // withlogin//////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *service) UpdadateInfo(ctx context.Context, customer_id string, info *services.Customers) *assets_services.ServiceError {

	info.CustomerID = customer_id
	err := s.repository.UpdateCustomers(ctx, *info, func() error { return nil })
	if err != nil {
		return assets_services.NewError(400, err)
	}
	return nil
}
func (s *service) UpdadateAvatar(ctx context.Context, customer_id string, file *multipart.FileHeader) (err *assets_services.ServiceError) {
	// create fiile

	filePathS := fmt.Sprintf("%s%s%s", s.env.ImagePath, customer_id, file.Filename)
	//save path to db
	fmt.Println("file", filePathS)
	errors := s.repository.UpdateCustomers(ctx, services.Customers{
		CustomerID: customer_id,
		Image:      services.Narg[string]{Data: filePathS, Valid: true},
	}, func() error {
		err := assets_services.SaveUploadedFile(file, filePathS)
		return err
	})
	if errors != nil {
		return assets_services.NewError(400, errors)
	}
	// save file

	// if errors != nil {
	// 	return assets_services.NewError(400, errors)
	// }
	return nil
}

func (s *service) CreateCustomerAddress(ctx context.Context, customer_id string, address_info *services.CustomerAddress) (err *assets_services.ServiceError) {
	//save path to db
	errors := s.repository.CreateCustomerAddresses(ctx, customer_id, &services.CustomerAddress{
		CustomerID:  customer_id,
		Address:     address_info.Address,
		PhoneNumber: address_info.PhoneNumber,
		IDAddress:   uuid.New().String(),
	})
	if errors != nil {
		return assets_services.NewError(400, errors)
	}

	return nil
}
func (s *service) UpdateCustomerAddress(ctx context.Context, customer_id string, address_info *services.CustomerAddress) (err *assets_services.ServiceError) {
	//save path to db
	errors := s.repository.UpdateCustomerAddresses(ctx, customer_id, &services.CustomerAddress{
		CustomerID:  customer_id,
		Address:     address_info.Address,
		PhoneNumber: address_info.PhoneNumber,
		IDAddress:   address_info.IDAddress,
	})
	if errors != nil {
		return assets_services.NewError(400, errors)
	}

	return nil
}
func (s *service) InfoUser(ctx context.Context, customer_id string) (info map[string]interface{}, err *assets_services.ServiceError) {
	user, addrs, error_s := s._getInfoCustomer(ctx, customer_id, "")
	if error_s != nil {
		return info, assets_services.NewError(400, err)
	}
	info, error_s = assets_services.HideFields(user, "", "account_id", "customer_id")
	if error_s != nil {
		return info, assets_services.NewError(400, err)
	}
	info["addrs"] = addrs
	return info, nil
}
func (s *service) ListAddress(ctx context.Context, customer_id string) (info map[string]interface{}, err *assets_services.ServiceError) {
	addrs, error_s := s.repository.ListCustomerAddresses(ctx, customer_id)
	if error_s != nil {
		return info, assets_services.NewError(400, error_s)
	}
	info, error_s = assets_services.HideFields(addrs, "addrs")
	if error_s != nil {
		return info, assets_services.NewError(400, error_s)
	}
	return info, nil
}

func (s *service) _getInfoCustomer(ctx context.Context, customer_id string, acount_id string) (customer services.Customers, addresss []services.CustomerAddress, err error) {
	// get info by customer_id
	if customer_id != "" {
		customer, err = s.repository.GetCustomer(ctx, customer_id)
	} else if acount_id != "" {
		customer, err = s.repository.GetCustomerByAccountID(ctx, acount_id)
	}
	if err != nil {
		return services.Customers{}, nil, err
	}
	addresss, err = s.repository.ListCustomerAddresses(ctx, customer.CustomerID)

	return customer, addresss, err
}
