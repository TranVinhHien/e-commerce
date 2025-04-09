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
	result, err := assets_services.HideFields(user, "account_id", "customer_id")
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

func (s *service) Login(ctx context.Context, userName, password string) (accessToken, refreshToken string, info map[string]interface{}, error *assets_services.ServiceError) {
	account, roleID, err := s.repository.Login(ctx, userName)
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

	if roleID == s.env.Customer {
		user, err := s.repository.GetCustomerByAccountID(ctx, account.AccountID)
		if err != nil {
			return "", "", nil, assets_services.NewError(400, err)
		}
		userID = user.CustomerID
		result123, err = assets_services.HideFields(user, "account_id", "customer_id")
		if err != nil {
			return "", "", nil, assets_services.NewError(400, err)
		}
	} else {
		user, err := s.repository.GetCustomerByAccountID(ctx, account.AccountID)
		if err != nil {
			return "", "", nil, assets_services.NewError(400, err)
		}
		userID = user.CustomerID
		result123, err = assets_services.HideFields(user, "account_id", "customer_id")
		if err != nil {
			return "", "", nil, assets_services.NewError(400, err)
		}
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
func (s *service) UpdadateInfo(ctx context.Context, user_id string, info *services.Customers) *assets_services.ServiceError {

	info.CustomerID = user_id
	err := s.repository.UpdateCustomers(ctx, *info, nil)
	if err != nil {
		return assets_services.NewError(400, err)
	}
	return nil
}
func (s *service) UpdadateAvatar(ctx context.Context, user_id string, file *multipart.FileHeader) (err *assets_services.ServiceError) {
	// create fiile

	filePathS := fmt.Sprintf("%s%s%s", s.env.ImagePath, user_id, file.Filename)
	//save path to db
	errors := s.repository.UpdateCustomers(ctx, services.Customers{
		CustomerID: user_id,
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
