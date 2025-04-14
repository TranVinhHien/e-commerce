package controllers

import (
	"net/http"
	assets_api "new-project/assets/api"
	"new-project/assets/token"
	controllers_model "new-project/controllers/models"
	services "new-project/services/entity"

	"github.com/gin-gonic/gin"
)

func (api *apiController) callbackMoMo() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req services.TransactionMoMO
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusNoContent, assets_api.SimpSuccessResponse("Hello", nil))
			return
		}
		api.service.CallBackMoMo(ctx, req)
		// fmt.Printf("momo tra ve : %s", req)
		ctx.JSON(http.StatusNoContent, assets_api.SimpSuccessResponse("Hello", nil))
	}
}
func (api *apiController) createOrder() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayload).(*token.Payload)
		var req controllers_model.OrderParams
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, assets_api.ResponseError(http.StatusBadRequest, err.Error()))
			return
		}
		NumOfProducts := make([]services.AmountProdduct, 0, len(req.NumOfProducts))
		for _, i := range req.NumOfProducts {
			NumOfProducts = append(NumOfProducts, services.AmountProdduct{
				Product_sku_id: i.Product_sku_id,
				Amount:         i.Amount,
			})
		}
		payment, errorr := api.service.CreateOrder(ctx, authPayload.Sub, &services.CreateOrderParams{
			NumOfProducts: NumOfProducts,
			Discount_Id:   req.Discount_Id,
			Address_id:    req.Address_id,
			Payment_id:    req.Payment_id,
		})

		if errorr != nil {
			ctx.JSON(errorr.Code, assets_api.ResponseError(errorr.Code, errorr.Error()))
			return
		}

		ctx.JSON(http.StatusOK, assets_api.SimpSuccessResponse("create order success", payment))
	}
}
