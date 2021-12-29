package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	conf "price_api/price_server/config"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
)

func HandleGetUpdatePriceHeartbeat(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")

	updatePriceService := service.Svc.CoinHistory()
	heartbeatInfo, err := updatePriceService.GetUpdatePriceHeartbeat(symbol, conf.GCfg.Interval)

	if err != nil {
		logger.WithError(err).Error("get update price heart beat occur err")
		response.Code = constant.GET_LOG_INFO_ERROR
		response.Message = constant.MSG_GET_LOG_INFO_ERROR
		context.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Data = heartbeatInfo
	context.JSON(http.StatusOK, response)
}
