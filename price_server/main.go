package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	conf "price_api/price_server/config"
	exchange "price_api/price_server/exchange"
	"price_api/price_server/sql"
	"strings"

	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var gPriceInfosCache conf.PriceInfosCache
var m *sync.RWMutex
var gCfg conf.Config
var gRequestPriceConfs map[string][]conf.ExchangeConfig

func main() {
	m = new(sync.RWMutex)
	//gRequestPriceConfs = make(map[string][]conf.ExchangeConfig)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	cfg, err := conf.GetConfig()
	if err != nil {
		log.Println(err)
		return
	}

	gCfg = cfg
	log.Println("config load over:", cfg)

	gRequestPriceConfs = exchange.InitRequestPriceConf(cfg)

	err = sql.InitMysqlDB(cfg)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("mysql init over")

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
	router.GET("/api/getReqConfig", HandleGetReqConfig)
	router.GET("/api/getRequestInfo", HandleGetRequestInfo)
	router.GET("/api/getHttpErrorInfo", HandleGetHttpErrorInfo)
	router.GET("/api/getLocalPrices", Check(), HandleGetLocalPrices)
	router.GET("/api/getAresAll", HandleGetAresAll)

	go updatePrice(cfg, gRequestPriceConfs)
	router.Run(":" + strconv.Itoa(int(cfg.Port)))
}

func updatePrice(cfg conf.Config, reqConf map[string][]conf.ExchangeConfig) {
	idx := 0

	for {
		infos, err := exchange.GetExchangePrice(reqConf, cfg)
		if err != nil {
			log.Println(err)
		} else {
			idx++
			m.Lock()
			gPriceInfosCache.PriceInfosCache = append(gPriceInfosCache.PriceInfosCache, infos)
			if len(gPriceInfosCache.PriceInfosCache) > int(cfg.MaxMemTime) {
				gPriceInfosCache.PriceInfosCache = gPriceInfosCache.PriceInfosCache[1:]
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

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
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

		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		startTime := time.Now().Format("2006-01-02 15:04:05")
		c.Next()

		responseBody := bodyLogWriter.body.String()

		endTime := time.Now().Format("2006-01-02 15:04:05")

		if c.Request.Method == "POST" {
			c.Request.ParseForm()
		}

		accessLogMap := make(map[string]string)

		accessLogMap["request_time"] = startTime
		//accessLogMap["request_method"] = c.Request.Method
		accessLogMap["request_uri"] = c.Request.RequestURI
		//accessLogMap["request_proto"] = c.Request.Proto
		accessLogMap["request_ua"] = c.Request.UserAgent()
		//accessLogMap["request_referer"] = c.Request.Referer()
		//accessLogMap["request_post_data"] = c.Request.PostForm.Encode()
		accessLogMap["request_client_ip"] = c.ClientIP()

		accessLogMap["response_time"] = endTime
		accessLogMap["response"] = responseBody

		accessLogJson, _ := json.Marshal(accessLogMap)

		log.Println(string(accessLogJson))

		if !strings.Contains(accessLogMap["request_uri"], "getRequestInfo") {
			err := sql.InsertLogInfo(accessLogMap)
			if err != nil {
				log.Println(err)
			}
		}
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
