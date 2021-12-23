package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
	"strings"
)

func HandleGetPrice(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	lastIndex := strings.LastIndex(context.Param("name")[1:], "/")
	if lastIndex == -1 {
		logger.Infoln("not true param name", context.Param("name")[1:])
		response.Code = constant.NO_MATCH_FORMAT_ERROR
		response.Message = constant.MSG_URL_NOT_FIND
		context.JSON(http.StatusBadRequest, response)
		return
	}

	symbol := context.Param("name")[1 : lastIndex+1]
	exchange := context.Param("name")[lastIndex+2:]

	priceService := service.Svc.Price()
	bFind, rspData := priceService.GetPrice(symbol, exchange)

	if !bFind {
		logger.Infoln("symbol or exchange not find, symbol:", symbol, " exchange:", exchange)
		response.Code = constant.NO_MATCH_FORMAT_ERROR
		response.Message = constant.MSG_URL_NOT_FIND
		context.JSON(http.StatusNotFound, response)
		return
	}

	response.Data = rspData
	context.JSON(http.StatusOK, response)
}
