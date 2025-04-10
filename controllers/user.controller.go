package controllers

import (
	"net/http"
	assets_api "new-project/assets/api"
	"new-project/assets/token"
	controllers_model "new-project/controllers/models"
	services "new-project/services/entity"

	"github.com/gin-gonic/gin"
)

// create a json to test this controller
func (api *apiController) optionss() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Lỗi zo options", nil))
	}
}

func (api *apiController) updatePassword() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		var param struct {
			UserName    string `json:"username"`
			OldPassword string `json:"old_password"`
			NewPassword string `json:"new_password"`
		}

		if err := ctx.ShouldBindJSON(&param); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}

		err := api.service.UpdatePassword(ctx.Request.Context(), param.UserName, param.OldPassword, param.NewPassword)

		if err != nil {
			ctx.AbortWithStatusJSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Update successful", nil))
	}
}

// func (api *apiController) getInfo() func(c *gin.Context) {
// 	return func(ctx *gin.Context) {
// 		username := ctx.Param("username")

// 		result, err := api.service.InfoUser(ctx.Request.Context(), username)

//			if err != nil {
//				ctx.AbortWithStatusJSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
//				return
//			}
//			ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get successful", result))
//		}
//	}
func (api *apiController) login() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req controllers_model.LoginParams
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}
		accessToken, refershToken, user, err := api.service.Login(ctx, req.Username, req.Password)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("login successful", map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refershToken,
			"user":          user,
		}))
	}
}
func (api *apiController) logout() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req controllers_model.LogOutParams
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}
		err := api.service.Logout(ctx, req.RefreshToken)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("logout successful", nil))
	}
}
func (api *apiController) register() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req controllers_model.RegisterParams
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}
		err := api.service.Register(ctx, req.Username, req.Password, &services.Customers{
			Name:   req.Name,
			Email:  req.Email,
			Dob:    req.Dob,
			Gender: req.Gender,
		})
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("register successful", nil))
	}
}
func (api *apiController) newAccessToken() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req controllers_model.LogOutParams
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}

		token, err := api.service.NewAccessToken(ctx, req.RefreshToken)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get new Token successful", map[string]interface{}{
			"access_token": token,
		}))
	}
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
func (api *apiController) updateCustomer() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		var req controllers_model.Customers
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}

		errorr := api.service.UpdadateInfo(ctx, authPayload.Sub, &services.Customers{
			Name:   req.Name,
			Email:  req.Email,
			Dob:    req.Dob,
			Gender: req.Gender,
		})

		if errorr != nil {
			ctx.JSON(errorr.Code, assets_api.ResponseError(errorr.Code, errorr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("update account success", nil))
	}
}
func (api *apiController) updateCustomerAddress() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		var req controllers_model.CustomersAddress
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}
		errorr := api.service.UpdateCustomerAddress(ctx, authPayload.Sub, &services.CustomerAddress{
			PhoneNumber: req.PhoneNumber,
			Address:     req.Address,
			IDAddress:   req.Address_id,
		})

		if errorr != nil {
			ctx.JSON(errorr.Code, assets_api.ResponseError(errorr.Code, errorr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("update account success", nil))
	}
}
func (api *apiController) createCustomerAddress() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		var req controllers_model.CustomersAddress
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}

		errorr := api.service.CreateCustomerAddress(ctx, authPayload.Sub, &services.CustomerAddress{
			PhoneNumber: req.PhoneNumber,
			Address:     req.Address,
		})

		if errorr != nil {
			ctx.JSON(errorr.Code, assets_api.ResponseError(errorr.Code, errorr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("create account_address success", nil))
	}
}
func (api *apiController) updateCustomerAvatar() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		file, err := ctx.FormFile("avatar")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		errorr := api.service.UpdadateAvatar(ctx, authPayload.Sub, file)
		if errorr != nil {
			ctx.JSON(errorr.Code, assets_api.ResponseError(errorr.Code, errorr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("update account_address success", nil))
	}
}
func (api *apiController) infoUser() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

		info, err := api.service.InfoUser(ctx, authPayload.Sub)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get new Token successful", map[string]interface{}{
			"user": info,
		}))
	}
}

func (api *apiController) listAddress() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)

		addrs, err := api.service.ListAddress(ctx, authPayload.Sub)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get new Token successful", addrs))
	}
}
