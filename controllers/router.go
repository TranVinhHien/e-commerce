package controllers

import (
	"new-project/assets/token"
	"new-project/services"

	"github.com/gin-gonic/gin"
)

type apiController struct {
	service services.UserUseCase
	jwt     token.Maker
}

func NewAPIController(s services.UserUseCase, jwt token.Maker) apiController {
	return apiController{service: s, jwt: jwt}
}

func (api apiController) SetUpRoute(group *gin.RouterGroup) {

	user := group.Group("/user")
	{
		user.GET("/getinfo/:username", api.getInfo())
		user.POST("/updatePassword", api.updatePassword())
		user.POST("/login", api.login())
		user.POST("/logout", api.logout())
		user.POST("/register", api.register())
		user.OPTIONS("/login", api.optionss())
		user.OPTIONS("/new_access_token", api.optionss())
		user.POST("/new_access_token", api.newAccessToken())
	}
	dalogin := group.Group("/dalogin").Use(authorization(api.jwt))
	{
		dalogin.GET("/ghi", api.checkAuth())
		dalogin.OPTIONS("/ghi", api.checkAuth())
	}
}
