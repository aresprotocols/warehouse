package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
)

func Check() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := vo.RESPONSE{Code: 0, Message: "OK"}
		priceService := service.Svc.Price()
		infoLen := priceService.GetCacheLength()
		if infoLen == 0 {
			response.Code = -1
			response.Message = constant.MSG_PRICE_NOT_READY
			context.JSON(http.StatusInternalServerError, response)
			context.Abort()
			return
		}

		context.Next()
	}
}
