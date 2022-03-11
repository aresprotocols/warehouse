package routers

import (
	"github.com/gin-gonic/gin"
	conf "price_api/price_server/internal/config"
	"price_api/price_server/internal/handler"
	"price_api/price_server/internal/middleware"
)

func NewRouter(conf conf.Config) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(middleware.Cors())

	router.GET("/api/getPrice/*name", middleware.Check(), handler.HandleGetPrice)
	router.GET("/api/getPartyPrice/:symbol", middleware.Check(), handler.HandleGetPartyPrice)
	router.GET("/api/getPriceAll/:symbol", middleware.Check(), handler.HandleGetPriceAll)
	router.GET("/api/getHistoryPrice/:symbol", handler.HandleGetHistoryPrice)
	router.GET("/api/getBulkPrices", middleware.Check(), handler.HandleGetBulkPrices)
	router.GET("/api/getBulkCurrencyPrices", middleware.Check(), handler.HandleGetBulkCurrencyPrices)
	router.GET("/api/getReqConfig", handler.HandleGetReqConfig)
	router.GET("/api/getRequestInfo", middleware.JWTAuthMiddleware(), handler.HandleGetRequestInfo)
	router.GET("/api/getRequestInfoBySymbol", handler.HandleGetRequestInfoBySymbol)
	router.GET("/api/getHttpErrorInfo/:symbol", handler.HandleGetHttpErrorInfo)
	router.GET("/api/getLocalPrices", middleware.Check(), handler.HandleGetLocalPrices)
	router.GET("/api/getUpdatePriceHistory", handler.HandleGetUpdatePriceHistory)
	router.POST("/api/setWeight", middleware.JWTAuthMiddleware(), middleware.Check(), handler.HandleSetWeight)
	router.GET("/api/getAresAll", handler.HandleGetAresAll)
	router.GET("/api/getDexPrice", handler.HandleGetDexPrice)
	router.POST("/api/auth", handler.HandleAuth)
	router.GET("/api/getUpdatePriceHeartbeat/:symbol", middleware.Check(), handler.HandleGetUpdatePriceHeartbeat)
	router.GET("/api/getBulkSymbolsState", middleware.Check(), handler.HandleGetBulkSymbolsState)
	router.POST("/api/setInterval", middleware.JWTAuthMiddleware(), middleware.Check(), handler.HandleSetInterval)

	if !conf.RunByDocker {
		router.GET("/api/gas/cal", middleware.Check(), handler.CalGasFee)
	}
	return router
}
