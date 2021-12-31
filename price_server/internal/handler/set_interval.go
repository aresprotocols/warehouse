package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
)

func HandleSetInterval(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	var setIntervalReq vo.SetIntervalReq
	err := context.ShouldBind(&setIntervalReq)
	if err != nil {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARSE_PARAM_ERROR
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if len(setIntervalReq.Symbol) == 0 {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if setIntervalReq.Interval == 0 {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	updateIntervalService := service.Svc.UpdateInterval()

	err = updateIntervalService.SetUpdateInterval(setIntervalReq.Symbol, setIntervalReq.Interval)
	if err != nil {
		logger.WithError(err).Errorf("set update interval occur error")
		if err != nil {
			response.Code = constant.SET_UPDATE_INTERVAL_ERROR
			response.Message = err.Error()
			context.JSON(http.StatusInternalServerError, response)
			return
		}
	}

	context.JSON(http.StatusOK, response)
}
