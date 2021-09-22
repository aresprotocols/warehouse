package main

import (
	"log"
	"net/http"
	conf "price_api/price_server/config"
	exchange "price_api/price_server/exchange"
	"price_api/price_server/sql"

	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var gPriceInfosCache conf.PriceInfosCache
var m *sync.RWMutex
var gCfg conf.Config

func main() {
	m = new(sync.RWMutex)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	cfg, err := conf.GetConfig()
	if err != nil {
		log.Println(err)
		return
	}

	gCfg = cfg
	log.Println("config load over:", cfg)

	err = sql.InitMysqlDB(cfg)
	if err != nil {
		log.Println(err)
		return
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(Cors())
	router.GET("/index", HandleHello)
	router.GET("/api/getprice/*name", Check(), HandleGetPrice)
	router.GET("/api/getPartyPrice/:symbol", Check(), HandleGetPartyPrice)
	router.GET("/api/getPriceAll/:symbol", Check(), HandleGetPriceAll)
	router.GET("/api/getConfigWeight", HandleGetConfigWeight)
	router.GET("/api/getHistoryPrice/:symbol", HandleGetHistoryPrice)
	router.GET("/api/getBulkPrices", Check(), HandleGetBulkPrices)
	router.GET("/api/getAresAll", HandleGetAresAll)

	go updatePrice(cfg)
	router.Run(":" + strconv.Itoa(int(cfg.Port)))
}

func updatePrice(cfg conf.Config) {
	idx := 0

	for {
		infos, err := exchange.GetExchangePrice(cfg)
		if err != nil {
			log.Println(err)
		} else {
			idx++
			m.Lock()
			gPriceInfosCache.PriceInfosCache = append(gPriceInfosCache.PriceInfosCache, infos)
			if len(gPriceInfosCache.PriceInfosCache) == int(cfg.MaxVolume) {
				gPriceInfosCache.PriceInfosCache = gPriceInfosCache.PriceInfosCache[cfg.MaxVolume/2:]
			}
			m.Unlock()
		}

		if idx >= int(cfg.InsertInterval) {
			err = sql.InsertPriceInfo(infos)
			if err != nil {
				log.Println(err)
			} else {
				idx = 0
			}
		}
		time.Sleep(time.Second * time.Duration(cfg.Interval))
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

func Check() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := RESPONSE{Code: 0, Message: "OK"}

		m.RLock()
		infoLen := len(gPriceInfosCache.PriceInfosCache)
		m.RUnlock()
		if infoLen == 0 {
			response.Code = -1
			response.Message = MSG_PRICE_NOT_READY
			context.JSON(http.StatusOK, response)
			context.Abort()
			return
		}

		context.Next()
	}
}

type RESPONSE struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
