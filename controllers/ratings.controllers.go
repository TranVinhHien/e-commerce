package controllers

import (
	"net/http"
	assets_api "new-project/assets/api"
	"new-project/assets/token"
	controllers_model "new-project/controllers/models"
	services "new-project/services/entity"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *apiController) listRating() func(ctx *gin.Context) {
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
		products_spu_id := ctx.Query("products_spu_id")
		if products_spu_id == "" {
			ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("you need provider products_spu_id", nil))
			return
		}
		order := &services.OrderBy{}
		if order_by != "" {
			order.Field = order_by
			order.Field = order_option
		} else {
			order = nil
		}
		ratings, err := api.service.ListRating(ctx, products_spu_id, services.NewQueryFilter(pageInt, pageSizeInt, nil, order))
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get categories successfull", ratings))
	}
}

func (api *apiController) createRating() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		var req controllers_model.RatingParams
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}

		err := api.service.CreateRating(ctx, authPayload.Sub, services.Ratings{
			Comment:       services.Narg[string]{Data: req.Comment, Valid: req.Comment != ""},
			Star:          int32(req.Star),
			ProductsSpuID: req.ProductsSpuID,
		})
		if err != nil {
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("create comment successfull", nil))

	}
}
