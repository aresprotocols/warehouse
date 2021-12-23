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

func HandleGetRequestInfo(context *gin.Context) {
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

	requestInfoService := service.Svc.RequestInfo()
	logInfos, err := requestInfoService.GetLogInfos(idx, int(conf.GCfg.PageSize))

	if err != nil {
		response.Code = constant.GET_LOG_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusBadRequest, response)
		return
	}

	response.Data = logInfos
	context.JSON(http.StatusOK, response)
}
