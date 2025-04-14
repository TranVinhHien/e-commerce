package controllers

import (
	"net/http"
	assets_api "new-project/assets/api"
	services "new-project/services/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *apiController) listDiscount() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		page := ctx.DefaultQuery("page", "1")
		limit := ctx.DefaultQuery("limit", "10")
		pageInt, errors := strconv.Atoi(page)
		if errors != nil {
			ctx.JSON(402, assets_api.ResponseError(402, errors.Error()))
			return
		}
		pageSizeInt, errors := strconv.Atoi(limit)
		if errors != nil {
			ctx.JSON(402, assets_api.ResponseError(402, errors.Error()))
			return
		}
		order_by := ctx.Query("order_by")
		order_option := ctx.Query("order_option")

		order := &services.OrderBy{}
		if order_by != "" {
			order.Field = order_by
			order.Field = order_option
		} else {
			order = nil
		}
		addrs, err := api.service.ListDiscount(ctx, services.QueryFilter{
			Page:     pageInt,
			PageSize: pageSizeInt,
			OrderBy:  order,
		})
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get categories successfull", addrs))
	}
}
