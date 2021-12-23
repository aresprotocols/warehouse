package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	conf "price_api/price_server/config"
	"price_api/price_server/internal/vo"
)

func HandleGetReqConfig(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	type EXCHANGE_WEIGHT_INFO struct {
		Exchange string `json:"exchange"`
		Weight   int64  `json:"weight"`
	}

	data := make(map[string][]EXCHANGE_WEIGHT_INFO)
	for symbol, confList := range conf.GRequestPriceConfs {
		for _, confTemp := range confList {
			data[symbol] = append(data[symbol], EXCHANGE_WEIGHT_INFO{Exchange: confTemp.Name, Weight: confTemp.Weight})
		}
	}

	response.Data = data
	context.JSON(http.StatusOK, response)
}
