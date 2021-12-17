package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	conf "price_api/price_server/config"
	"price_api/price_server/exchange"
	"price_api/price_server/sql"
	"strconv"
	"strings"
	"sync"
	"time"
)

var gPriceInfosCache conf.PriceInfosCache
var m *sync.RWMutex
var gCfg conf.Config
var gRequestPriceConfs map[string][]conf.ExchangeConfig

func init() {
	config := DefaultConfiguration()
	err := InitLogrusLogger(config)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}

}

func main() {
	m = new(sync.RWMutex)
	//gRequestPriceConfs = make(map[string][]conf.ExchangeConfig)

	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	cfg, err := conf.GetConfig()
	if err != nil {
		logger.Errorf("get config occur err:%v", err)
		return
	}

	gCfg = cfg
	logger.Infof("config load over:%v", cfg)

	err = sql.InitMysqlDB(cfg)
	if err != nil {
		logger.Errorf("Init mysql db occur err:%v", err)
		return
	}

	logger.Info("mysql init over")

	handle := InitHandle(cfg)

	gRequestPriceConfs, err = exchange.InitRequestPriceConf(cfg)
	if err != nil {
		logger.Errorf("Init request price conf occur err:%v", err)
		return
	}
	logger.Info("request init over")

	showIgnoreSymbols(cfg, gRequestPriceConfs)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(Cors())

	router.GET("/api/getPrice/*name", Check(), HandleGetPrice)
	router.GET("/api/getPartyPrice/:symbol", Check(), HandleGetPartyPrice)
	router.GET("/api/getPriceAll/:symbol", Check(), HandleGetPriceAll)
	router.GET("/api/getHistoryPrice/:symbol", HandleGetHistoryPrice)
	router.GET("/api/getBulkPrices", Check(), HandleGetBulkPrices)
	router.GET("/api/getBulkCurrencyPrices", Check(), HandleGetBulkCurrencyPrices)
	router.GET("/api/getReqConfig", HandleGetReqConfig)
	router.GET("/api/getRequestInfo", JWTAuthMiddleware(), HandleGetRequestInfo)
	router.GET("/api/getRequestInfoBySymbol", HandleGetRequestInfoBySymbol)
	router.GET("/api/getHttpErrorInfo/:symbol", HandleGetHttpErrorInfo)
	router.GET("/api/getLocalPrices", Check(), HandleGetLocalPrices)
	router.GET("/api/getUpdatePriceHistory", HandleGetUpdatePriceHistory)
	router.POST("/api/setWeight", JWTAuthMiddleware(), Check(), HandleSetWeight)
	router.GET("/api/getAresAll", HandleGetAresAll)
	router.GET("/api/getDexPrice", HandleGetDexPrice)
	router.POST("/api/auth", HandleAuth)
	router.GET("/api/getUpdatePriceHeartbeat/:symbol", Check(), HandleGetUpdatePriceHeartbeat)
	router.GET("/api/getBulkSymbolsState", Check(), HandleGetBulkSymbolsState)

	go updatePrice(cfg, gRequestPriceConfs)
	router.Run(":" + strconv.Itoa(int(cfg.Port)))

	abortChan := make(chan os.Signal, 1)
	signal.Notify(abortChan, os.Interrupt)

	sig := <-abortChan
	handle.Stop()
	logger.Infof("Exiting... signal %v", sig)
}

func updatePrice(cfg conf.Config, reqConf map[string][]conf.ExchangeConfig) {
	idx := 0
	time.Sleep(time.Second * 2) // run update for the first time,  need to sleep , because you have just completed initialization and have already requested data once
	for {
		logger.Infof("start new round update price")
		infos, err := exchange.GetExchangePrice(reqConf, cfg)
		if err != nil {
			logger.WithError(err).Errorf("get exchange price occur error")
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
				logger.Errorf("insert price info occur err:%v", err)
			} else {
				idx = 0
			}
		}
		logger.Infof("end this round update price")
		time.Sleep(time.Second * time.Duration(cfg.Interval))

	}
}

func showIgnoreSymbols(cfg conf.Config, gRequestPriceConfs map[string][]conf.ExchangeConfig) {
	ignoreSymbols := make(map[string][]string)
	for _, symbol := range cfg.Symbols {
		var exchanges []string
		existSymbols, ok := gRequestPriceConfs[symbol]
		if ok {
			for _, exchangeConf := range cfg.Exchanges {
				//check config exchange if have symbol
				bFind := false
				for _, existSymbol := range existSymbols {
					if exchangeConf.Name == existSymbol.Name {
						//find it
						bFind = true
					}
				}
				if !bFind {
					exchanges = append(exchanges, exchangeConf.Name)
				}
			}
		} else {
			for _, exchangeConf := range cfg.Exchanges {
				exchanges = append(exchanges, exchangeConf.Name)
			}
		}
		ignoreSymbols[symbol] = exchanges
	}
	logger.Infof("ignore symbols and exchange:", ignoreSymbols)
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
		startTimeStamp := time.Now().Unix()
		c.Next()

		responseBody := bodyLogWriter.body.String()

		endTime := time.Now().Format("2006-01-02 15:04:05")
		endTimeStamp := time.Now().Unix()

		if c.Request.Method == "POST" {
			c.Request.ParseForm()
		}
		if bodyLogWriter.Status() != http.StatusOK { // not insert log if http status not ok
			return
		}

		accessLogMap := make(map[string]interface{})

		requestUri := c.Request.RequestURI

		accessLogMap["request_time"] = startTime
		accessLogMap["request_uri"] = requestUri
		accessLogMap["request_ua"] = c.Request.UserAgent()
		accessLogMap["request_client_ip"] = c.ClientIP()

		accessLogMap["response_time"] = endTime
		accessLogMap["response"] = responseBody
		accessLogMap["request_timestamp"] = startTimeStamp
		accessLogMap["response_timestamp"] = endTimeStamp

		if strings.Contains(requestUri, "getPrice") ||
			strings.Contains(requestUri, "getPartyPrice") ||
			strings.Contains(requestUri, "getHistoryPrice") ||
			strings.Contains(requestUri, "getBulkPrices") {
			err := sql.InsertLogInfo(accessLogMap, 1)
			if err != nil {
				logger.Errorf("insert log info occur err:%v", err)
			}
		} else {
			err := sql.InsertLogInfo(accessLogMap, 0)
			if err != nil {
				logger.Errorf("insert log info occur err:%v", err)
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
			context.JSON(http.StatusInternalServerError, response)
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

type Pagination struct {
	CurPage  int         `json:"curPage"`
	TotalNum int         `json:"totalNum"`
	Items    interface{} `json:"items"`
}
