package controllers

import (
	"new-project/assets/token"
	"new-project/services"

	"github.com/gin-gonic/gin"
)

type apiController struct {
	service services.ServiceUseCase
	jwt     token.Maker
}

func NewAPIController(s services.ServiceUseCase, jwt token.Maker) apiController {
	return apiController{service: s, jwt: jwt}
}

func (api apiController) SetUpRoute(group *gin.RouterGroup) {
	group.OPTIONS("/*any", api.optionss())
	user := group.Group("/user")
	{
		// user.GET("/getinfo/:username", api.getInfo())
		user.POST("/updatePassword", api.updatePassword())
		user.POST("/login", api.login())
		user.POST("/logout", api.logout())
		user.POST("/register", api.register())
		user.POST("/new_access_token", api.newAccessToken())
		//
		user_auth := user.Group("/info").Use(authorization(api.jwt))
		{
			user_auth.PATCH("/update_customer", api.updateCustomer())
			user_auth.PATCH("/update_avatar_customer", api.updateCustomerAvatar())
			user_auth.PATCH("/update_customeraddress", api.updateCustomerAddress())
			user_auth.POST("/create_customeraddress", api.createCustomerAddress())
			user_auth.GET("/get", api.infoUser())
			user_auth.GET("/get_address", api.listAddress())
		}
	}

	dalogin := group.Group("/dalogin").Use(authorization(api.jwt))
	{
		dalogin.GET("/ghi", api.checkAuth())
		// dalogin.OPTIONS("/ghi", api.optionss())
	}
	media := group.Group("/media")
	{
		media.GET("/images/:id", api.renderImages())
	}
	categories := group.Group("/categories")
	{
		categories.GET("/get", api.listCategories())
	}
}
