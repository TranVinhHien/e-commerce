package controllers

import (
	"fmt"
	"net/http"
	assets_api "new-project/assets/api"

	"github.com/gin-gonic/gin"
)

func (api *apiController) listPayment() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		addrs, err := api.service.ListPayment(ctx)
		if err != nil {
			fmt.Print(err)
			ctx.JSON(err.Code, assets_api.ResponseError(err.Code, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("get categories successfull", addrs))
	}
}
