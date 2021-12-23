package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
)

func HandleSetWeight(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	var setWeightReq vo.SetWeightReq
	err := context.ShouldBind(&setWeightReq)
	if err != nil {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if len(setWeightReq.Symbol) == 0 {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if len(setWeightReq.Exchange) == 0 {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	weightService := service.Svc.WeightInfo()

	err = weightService.SetWeight(setWeightReq)
	if err != nil {
		logger.WithError(err).Errorf("set weight occur error")
		if err != nil {
			response.Code = constant.SET_WEIGHT_ERROR
			response.Message = err.Error()
			context.JSON(http.StatusInternalServerError, response)
			return
		}
	}

	context.JSON(http.StatusOK, response)
}
