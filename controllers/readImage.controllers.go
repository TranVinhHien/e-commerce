package controllers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func (api *apiController) renderAvatars() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		filename := ctx.Param("id")
		filePath := api.service.RenderImage(ctx, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Nếu file không tồn tại, trả về lỗi 404
			ctx.JSON(400, gin.H{"error": "File not found"})
			return
		}
		ctx.File(filePath)
	}
}

func (api *apiController) renderProductImages() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		filename := ctx.Query("id")
		filePath := api.service.RenderProductImages(ctx, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Nếu file không tồn tại, trả về lỗi 404
			ctx.JSON(400, gin.H{"error": "File not found"})
			return
		}
		ctx.File(filePath)
	}
}
