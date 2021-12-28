package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	conf "price_api/price_server/config"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
	"strconv"
	"strings"
)

func HandleGetHttpErrorInfo(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")
	symbol = strings.ToLower(symbol)
	if !strings.Contains(symbol, "-") {
		if strings.HasSuffix(symbol, "usdt") {
			symbol = strings.ReplaceAll(symbol, "usdt", "-usdt")
		}
	}

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

	httpErrorService := service.Svc.HttpError()

	total, httpErrorInfos, err := httpErrorService.GetHttpErrorsByPage(idx, int(conf.GCfg.PageSize), symbol)
	if err != nil {
		response.Code = constant.GET_HTTP_ERROR_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Data = vo.Pagination{
		CurPage:  idx,
		TotalNum: total,
		Items:    httpErrorInfos,
	}
	context.JSON(http.StatusOK, response)
}
