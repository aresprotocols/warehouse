package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
)

func HandleGetBulkPrices(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	priceService := service.Svc.Price()

	mSymbolPriceInfo := priceService.GetBulkPrices(symbol)
	response.Data = mSymbolPriceInfo
	context.JSON(http.StatusOK, response)
}
