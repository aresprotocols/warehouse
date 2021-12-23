package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	conf "price_api/price_server/config"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
	"strconv"
)

func HandleGetLocalPrices(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	index, exist := context.GetQuery("index")
	if !exist {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		response.Code = constant.PARSE_PARAM_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusBadRequest, response)
		return
	}

	start := idx * int(conf.GCfg.PageSize)
	end := start + int(conf.GCfg.PageSize)

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	priceService := service.Svc.Price()
	retData := priceService.GetLocalPrices(start, end, symbol)
	response.Data = retData
	context.JSON(http.StatusOK, response)
}
