package services

import (
	config_assets "new-project/assets/config"
	"new-project/assets/token"
)

type service struct {
	repository ServicesRepository
	redis      ServicesRedis
	jwt        token.Maker
	env        config_assets.ReadENV
}

func NewService(repo ServicesRepository, jwt token.Maker, env config_assets.ReadENV, redis ServicesRedis) ServiceUseCase {
	return &service{repository: repo, jwt: jwt, env: env, redis: redis}
}
