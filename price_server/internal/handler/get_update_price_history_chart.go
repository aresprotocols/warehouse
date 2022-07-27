package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
)

func HandleGetUpdatePriceHistoryForChart(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}
	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	coinHistoryService := service.Svc.CoinHistory()

	historyResps, err := coinHistoryService.GetUpdatePriceHistoryForChart(symbol)
	if err != nil {
		logger.WithError(err).Errorf("get update price history for chart occur error,symbol:%s", symbol)
		response.Code = constant.GET_LOG_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = historyResps
	context.JSON(http.StatusOK, response)
}
