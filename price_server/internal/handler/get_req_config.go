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
	intervalService := service.Svc.UpdateInterval()

	data := make(map[string]vo.ReqConfigResp)
	for symbol, confList := range requestPriceConfService.GetConfs() {
		weightConfigs := make([]vo.EXCHANGE_WEIGHT_INFO, 0)
		for _, confTemp := range confList {
			weightConfigs = append(weightConfigs, vo.EXCHANGE_WEIGHT_INFO{Exchange: confTemp.Name, Weight: confTemp.Weight})
		}

		updateInterval := intervalService.GetIntervalFromCache(symbol)
		data[symbol] = vo.ReqConfigResp{
			Weight:   weightConfigs,
			Interval: updateInterval,
		}
	}

	response.Data = data
	context.JSON(http.StatusOK, response)
}
