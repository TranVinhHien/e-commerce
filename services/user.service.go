package services

import (
	"context"
	"fmt"
	util_assets "new-project/assets/util"
	assets_services "new-project/services/assets"
	services "new-project/services/entity"
	"time"

	"github.com/rs/zerolog/log"
)

func (s *service) GetInfo(ctx context.Context, userName string) (*services.Users, *assets_services.ServiceError) {
	// kiểm tra thông tin abc xyz'
	user, err := s.repository.GetUser(ctx, userName)
	return &user, assets_services.NewError(400, err)
}
func (s *service) UpdatePassword(ctx context.Context, userName, oldPassword, newPassword string) (*services.Users, *assets_services.ServiceError) {
	//
	user, err := s.repository.GetUser(ctx, userName)
	if err != nil {
		return nil, assets_services.NewError(400, err)
	}
	err = util_assets.CheckPassword(oldPassword, user.Password)
	if err != nil {
		log.Debug().Msg("user password: " + user.Password + " oldPassword: " + oldPassword)
		return nil, assets_services.NewError(400, fmt.Errorf("oldPassword is incorrect"))
	}
	user, err = s.repository.UpdateUser(ctx, services.Users{
		Username: user.Username,
		Password: newPassword,
	})
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("update failed for user: "+user.Username))
	}
	return &user, nil
}
func (s *service) Login(ctx context.Context, userName, password string) (accessToken, refreshToken string, info *services.Users, error_sv *assets_services.ServiceError) {
	user, err := s.repository.GetUser(ctx, userName)
	if err != nil {
		return "", "", nil, assets_services.NewError(400, err)
	}
	err = util_assets.CheckPassword(password, user.Password)
	if err != nil {
		return "", "", nil, assets_services.NewError(400, fmt.Errorf("password is incorrect"))
	}
	_, accessToken, err = s.jwt.CreateToken(userName, s.env.AccessTokenDuration)
	if err != nil {
		return "", "", nil, assets_services.NewError(400, fmt.Errorf("error create access token: %v", err))
	}
	_, refreshToken, err = s.jwt.CreateToken(userName, s.env.RefershTokenDuration)
	if err != nil {
		return "", "", nil, assets_services.NewError(400, fmt.Errorf("error create refersh token: %v", err))
	}
	return accessToken, refreshToken, &user, nil
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
func (s *service) Register(ctx context.Context, userName, password, fullName string) *assets_services.ServiceError {
	//
	password, err := util_assets.HashPassword(password)
	if err != nil {
		return assets_services.NewError(400, err)
	}
	_, err = s.repository.InsertUser(ctx, services.Users{
		Username: userName,
		Password: password,
		FullName: fullName,
		CreateAt: time.Now().UTC(),
	})
	if err != nil {
		return assets_services.NewError(400, err)
	}

	return nil
}
func (s *service) NewAccessToken(ctx context.Context, refreshToken string) (*string, *assets_services.ServiceError) {
	payload, err := s.jwt.VerifyToken(refreshToken)
	if err != nil {
		log.Info().Msg(fmt.Sprintln("s.env.AccessTokenDuration", s.env.AccessTokenDuration))
		log.Info().Msg(fmt.Sprintln("s.env.RefershTokenDuration", s.env.RefershTokenDuration))
		log.Info().Msg(fmt.Sprintln("refreshToken", refreshToken))
		return nil, assets_services.NewError(400, fmt.Errorf("error when verifyToken: %v", err))
	}
	_, accessToken, err := s.jwt.CreateToken(payload.Sub, s.env.AccessTokenDuration)
	if err != nil {
		return nil, assets_services.NewError(400, fmt.Errorf("error when create new accessToken: %v", err))
	}
	return &accessToken, nil
}
