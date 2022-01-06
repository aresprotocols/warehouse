package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"price_api/price_server/internal/config"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
	"strconv"
)

func HandleGetUpdatePriceHistory(context *gin.Context) {
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

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	coinHistoryService := service.Svc.CoinHistory()

	total, historyResps, err := coinHistoryService.GetUpdatePriceHistory(idx, int(conf.GCfg.PageSize), symbol)
	if err != nil {
		logger.WithError(err).Errorf("get update price history occur error,symbol:%s", symbol)
		response.Code = constant.GET_LOG_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = vo.Pagination{
		CurPage:  idx,
		TotalNum: total,
		Items:    historyResps,
	}
	context.JSON(http.StatusOK, response)
}
