package services

import (
	config_assets "new-project/assets/config"
	assets_firebase "new-project/assets/fire-base"
	assets_jobs "new-project/assets/jobs"
	"new-project/assets/token"
)

type service struct {
	repository ServicesRepository
	redis      ServicesRedis
	jwt        token.Maker
	env        config_assets.ReadENV
	firebase   *assets_firebase.FirebaseMessaging
	jobs       *assets_jobs.JobScheduler
}

func NewService(repo ServicesRepository, jwt token.Maker, env config_assets.ReadENV, redis ServicesRedis, firebase *assets_firebase.FirebaseMessaging, jobs *assets_jobs.JobScheduler) ServiceUseCase {
	return &service{repository: repo, jwt: jwt, env: env, redis: redis, firebase: firebase, jobs: jobs}
}
