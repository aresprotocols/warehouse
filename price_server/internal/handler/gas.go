package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/util"
	"price_api/price_server/internal/vo"
	"strconv"
)

func CalGasFee(ctx *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}
	gasStr := ctx.Query("gas")
	if gasStr == "" {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	gas, err := strconv.ParseInt(gasStr, 10, 64)
	if err != nil {
		logger.WithError(err).Errorln("parse to int occur err")
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = err.Error()
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	gasService := service.Svc.Gas()
	aresGasFee := gasService.CalGasFeeToAres(gas)
	ctx.JSON(http.StatusOK, util.KeepValidDecimals(aresGasFee, 7))
}
