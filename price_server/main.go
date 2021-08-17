package main

import (
	"log"
	"net/http"
	conf "price_api/price_server/config"
	exchange "price_api/price_server/exchange"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var gPriceInfosCache conf.PriceInfosCache
var m *sync.RWMutex

func main() {
	m = new(sync.RWMutex)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	cfg, err := conf.GetConfig()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("config load over:", cfg)

	router := gin.Default()

	router.Use(Cors())
	router.GET("/index", HandleHello)
	router.GET("/api/getprice/*name", Check(), HandleGetPrice)
	router.GET("/api/getPartyPrice/:symbol", Check(), HandleGetPartyPrice)
	router.GET("/api/getPriceAll/:symbol", Check(), HandleGetPriceAll)

	go updatePrice(cfg)
	router.Run(":" + strconv.Itoa(int(cfg.Port)))
}

func updatePrice(cfg conf.Config) {
	for {
		infos, err := exchange.GetExchangePrice(cfg)
		if err != nil {
			log.Println(err)
		} else {
			m.Lock()
			gPriceInfosCache.PriceInfosCache = append(gPriceInfosCache.PriceInfosCache, infos)
			m.Unlock()
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
			response.Message = "price not ready"
			context.JSON(http.StatusOK, response)
			context.Abort()
			return
		}

		context.Next()
	}
}

func HandleHello(context *gin.Context) {
	context.String(http.StatusOK, "Hello, world")
}

func HandleGetPrice(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	lastIndex := strings.LastIndex(context.Param("name")[1:], "/")
	if lastIndex == -1 {
		log.Println("not true param name", context.Param("name")[1:])
		response.Code = -1
		response.Message = "url not find"
		context.JSON(http.StatusOK, response)
		return
	}

	symbol := context.Param("name")[1 : lastIndex+1]
	exchange := context.Param("name")[lastIndex+2:]

	type RspData struct {
		Timestamps int64
		Price      float64
	}

	var rspData RspData
	bFind := false

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	for _, info := range latestInfos.PriceInfos {
		if strings.ToLower(info.Symbol) == strings.ToLower(symbol) &&
			strings.ToLower(info.PriceOrigin) == strings.ToLower(exchange) {
			bFind = true
			rspData.Price = info.Price
			rspData.Timestamps = info.TimeStamps
		}
	}
	m.RUnlock()

	if !bFind {
		log.Println("symbol or exchange not find, symbol:", symbol, " exchange:", exchange)
		response.Code = -1
		response.Message = "url not find"
		context.JSON(http.StatusOK, response)
		return
	}

	response.Data = rspData
	context.JSON(http.StatusOK, response)
}

func HandleGetPartyPrice(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")

	type RspData struct {
		Price      float64
		Timestamps int64
	}

	var rspData RspData
	bFind := false

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	totalPrice := 0.0
	totalWeight := int64(0)
	for _, info := range latestInfos.PriceInfos {
		if strings.ToLower(info.Symbol) == strings.ToLower(symbol) {
			bFind = true
			totalPrice += info.Price * float64(info.Weight)
			totalWeight += info.Weight
			rspData.Timestamps = info.TimeStamps
		}
	}
	m.RUnlock()

	if !bFind {
		log.Println("symbol or exchange not find, symbol:", symbol)
		response.Code = -1
		response.Message = "url not find"
		context.JSON(http.StatusOK, response)
		return
	}

	rspData.Price = totalPrice / float64(totalWeight)
	response.Data = rspData
	context.JSON(http.StatusOK, response)
}

func HandleGetPriceAll(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")

	bFind := false

	type PriceAllInfo struct {
		Name       string
		Symbol     string
		Price      float64
		Timestamps int64
	}

	var priceAll []PriceAllInfo

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]

	for _, info := range latestInfos.PriceInfos {
		if strings.ToLower(info.Symbol) == strings.ToLower(symbol) {
			bFind = true
			priceAllInfo := PriceAllInfo{Name: info.PriceOrigin,
				Symbol:     info.Symbol,
				Price:      info.Price,
				Timestamps: info.TimeStamps,
			}
			priceAll = append(priceAll, priceAllInfo)
		}
	}
	m.RUnlock()

	if !bFind {
		log.Println("symbol or exchange not find, symbol:", symbol)
		response.Code = -1
		response.Message = "url not find"
		context.JSON(http.StatusOK, response)
		return
	}

	response.Data = priceAll
	context.JSON(http.StatusOK, response)
}

type RESPONSE struct {
	Code    int
	Message string
	Data    interface{}
}
