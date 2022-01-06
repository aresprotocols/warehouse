package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"price_api/price_server/internal/config"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/exchange"
	"price_api/price_server/internal/vo"
)

func HandleGetAresAll(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	aresShowInfo, err := exchange.GetGateAresInfo(conf.GCfg.Proxy)
	aresShowInfo.Rank = handle.fetcher.GetCMCInfo().Rank
	aresShowInfo.MarketCap = handle.fetcher.GetCMCInfo().MarketCap
	aresShowInfo.Volume = handle.fetcher.GetCMCInfo().Volume

	if err != nil {
		logger.WithError(err).Errorf("get gate ares info occur error")
		response.Code = constant.GET_ARES_INFO_ERROR
		response.Message = constant.MSG_GET_ARES_ERROR
		context.JSON(http.StatusBadRequest, response)
	}

	response.Data = aresShowInfo
	context.JSON(http.StatusOK, response)
}
