package controllers

import (
	"fmt"
	"net/http"
	assets_api "new-project/assets/api"
	controllers_model "new-project/controllers/models"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// create a json to test this controller
func (api *apiController) optionss() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		fmt.Println(ctx)
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

		result, err := api.service.UpdatePassword(ctx.Request.Context(), param.UserName, param.OldPassword, param.NewPassword)

		if err != nil {
			ctx.AbortWithStatusJSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("Update successful", result))
	}
}

func (api *apiController) getInfo() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		username := ctx.Param("username")

		result, err := api.service.GetInfo(ctx.Request.Context(), username)

		if err != nil {
			ctx.AbortWithStatusJSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get successful", result))
	}
}
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
			"user":          user.FullName,
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
		err := api.service.Register(ctx, req.Username, req.Password, req.FullName)
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
		log.Info().Msg("đoạn 1")
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}

		log.Info().Msg("đoạn 2")
		token, err := api.service.NewAccessToken(ctx, req.RefreshToken)
		if err != nil {
			log.Error().Err(err).Msg("lỗi service")

			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		log.Info().Msg("đoạn 3")
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get new Token successful", map[string]interface{}{
			"access_token": token,
		}))
	}
}
