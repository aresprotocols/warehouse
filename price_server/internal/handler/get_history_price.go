package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
	"strconv"
)

func HandleGetHistoryPrice(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")
	timestampStr, exist := context.GetQuery("timestamp")
	if !exist {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	priceService := service.Svc.Price()

	bFind, partyPriceData := priceService.GetHistoryPrice(symbol, timestamp, true)

	if !bFind {
		logger.Infoln("symbol or exchange not find, symbol:", symbol)
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusNotFound, response)
		return
	}

	response.Data = partyPriceData
	context.JSON(http.StatusOK, response)
}
