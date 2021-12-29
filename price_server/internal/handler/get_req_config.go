package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
)

func HandleGetReqConfig(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	requestPriceConfService := service.Svc.RequestPriceConf()

	data := make(map[string][]vo.EXCHANGE_WEIGHT_INFO)
	for symbol, confList := range requestPriceConfService.GetConfs() {
		for _, confTemp := range confList {
			data[symbol] = append(data[symbol], vo.EXCHANGE_WEIGHT_INFO{Exchange: confTemp.Name, Weight: confTemp.Weight})
		}
	}

	response.Data = data
	context.JSON(http.StatusOK, response)
}
