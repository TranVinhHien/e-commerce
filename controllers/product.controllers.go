package controllers

import (
	"net/http"
	assets_api "new-project/assets/api"
	controllers_assets "new-project/controllers/assets"
	services "new-project/services/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *apiController) getAllProductSimple() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
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
		search := ctx.Query("search")
		condition, errors := controllers_assets.ParseConditions(search)
		if errors != nil {
			ctx.JSON(400, assets_api.ResponseError(400, errors.Error()))
			return
		}
		order := &services.OrderBy{}
		if order_by != "" {
			order.Field = order_by
			order.Field = order_option
		} else {
			order = nil
		}
		orders, err := api.service.GetAllProductSimple(ctx, services.NewQueryFilter(pageInt, pageSizeInt, condition, order))
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get categories successfull", orders))
	}
}
func (api *apiController) getDetailProduct() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		spu := ctx.Param("id")

		if spu == "" {
			ctx.JSON(402, assets_api.ResponseError(402, "must provide spu_id"))
			return
		}

		product_spu, err := api.service.GetDetailProduct(ctx, spu)
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get categories successfull", product_spu))
	}
}
