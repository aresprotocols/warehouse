package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	conf "price_api/price_server/config"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
	"strconv"
)

func HandleGetRequestInfoBySymbol(context *gin.Context) {
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

	requestLogService := service.Svc.RequestInfo()

	total, logInfos, err := requestLogService.GetRequestInfoBySymbol(idx, int(conf.GCfg.PageSize), symbol)

	if err != nil {
		logger.WithError(err).Errorf("get log info by symbol occur error,symbol:%s,index:%d", symbol, idx)
		response.Code = constant.GET_LOG_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = vo.Pagination{
		CurPage:  idx,
		TotalNum: total,
		Items:    logInfos,
	}
	context.JSON(http.StatusOK, response)
}
