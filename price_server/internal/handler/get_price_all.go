package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
)

func HandleGetPriceAll(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}
	symbol := context.Param("symbol")

	priceService := service.Svc.Price()
	bFind, priceAll := priceService.GetPriceAll(symbol)

	if !bFind {
		logger.Infoln("symbol or exchange not find, symbol:", symbol)
		response.Code = constant.NO_MATCH_FORMAT_ERROR
		response.Message = constant.MSG_URL_NOT_FIND
		context.JSON(http.StatusNotFound, response)
		return
	}

	response.Data = priceAll
	context.JSON(http.StatusOK, response)
}
