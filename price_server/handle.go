package main

import (
	"log"
	"net/http"
	conf "price_api/price_server/config"
	"price_api/price_server/sql"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const MSG_URL_NOT_FIND = "url not find"
const MSG_PRICE_NOT_READY = "price not ready"
const MSG_PARAM_NOT_TRUE = "param not true"

const (
	ERROR = iota - 1000
	NO_MATCH_FORMAT_ERROR
	PARAM_NOT_TRUE_ERROR
)

func HandleHello(context *gin.Context) {
	context.String(http.StatusOK, "Hello, world")
}

func HandleGetPrice(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	lastIndex := strings.LastIndex(context.Param("name")[1:], "/")
	if lastIndex == -1 {
		log.Println("not true param name", context.Param("name")[1:])
		response.Code = NO_MATCH_FORMAT_ERROR
		response.Message = MSG_URL_NOT_FIND
		context.JSON(http.StatusOK, response)
		return
	}

	symbol := context.Param("name")[1 : lastIndex+1]
	exchange := context.Param("name")[lastIndex+2:]

	type RspData struct {
		Timestamp int64   `json:"timestamp"`
		Price     float64 `json:"price"`
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
			rspData.Timestamp = info.TimeStamp
		}
	}
	m.RUnlock()

	if !bFind {
		log.Println("symbol or exchange not find, symbol:", symbol, " exchange:", exchange)
		response.Code = NO_MATCH_FORMAT_ERROR
		response.Message = MSG_URL_NOT_FIND
		context.JSON(http.StatusOK, response)
		return
	}

	response.Data = rspData
	context.JSON(http.StatusOK, response)
}

func HandleGetPartyPrice(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	m.RUnlock()

	bFind, partyPriceData := partyPrice(latestInfos.PriceInfos, symbol, true)

	if !bFind {
		log.Println("symbol or exchange not find, symbol:", symbol)
		response.Code = NO_MATCH_FORMAT_ERROR
		response.Message = MSG_URL_NOT_FIND
		context.JSON(http.StatusOK, response)
		return
	}

	response.Data = partyPriceData
	context.JSON(http.StatusOK, response)
}

func HandleGetPriceAll(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")

	bFind := false

	type PriceAllInfo struct {
		Name      string  `json:"name"`
		Symbol    string  `json:"symbol"`
		Price     float64 `json:"price"`
		Timestamp int64   `json:"timestamp"`
	}

	var priceAll []PriceAllInfo

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]

	for _, info := range latestInfos.PriceInfos {
		if strings.ToLower(info.Symbol) == strings.ToLower(symbol) {
			bFind = true
			priceAllInfo := PriceAllInfo{Name: info.PriceOrigin,
				Symbol:    info.Symbol,
				Price:     info.Price,
				Timestamp: info.TimeStamp,
			}
			priceAll = append(priceAll, priceAllInfo)
		}
	}
	m.RUnlock()

	if !bFind {
		log.Println("symbol or exchange not find, symbol:", symbol)
		response.Code = NO_MATCH_FORMAT_ERROR
		response.Message = MSG_URL_NOT_FIND
		context.JSON(http.StatusOK, response)
		return
	}

	response.Data = priceAll
	context.JSON(http.StatusOK, response)
}

func HandleGetConfigWeight(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	type WeightInfo struct {
		ExchangeName string `json:"exchangeName"`
		Weight       int64  `json:"weight"`
	}

	type ExchangesWeightInfo struct {
		WeightInfos []WeightInfo `json:"weightInfos"`
	}

	var exchangesWeightInfo ExchangesWeightInfo

	for _, info := range gCfg.Exchanges {
		exchangesWeightInfo.WeightInfos = append(exchangesWeightInfo.WeightInfos, WeightInfo{ExchangeName: info.Name, Weight: info.Weight})
	}

	response.Data = exchangesWeightInfo
	context.JSON(http.StatusOK, response)
}

type WeightInfo struct {
	Price        float64 `json:"price"`
	Weight       int64   `json:"weight"`
	ExchangeName string  `json:"exchangeName"`
}

type PartyPriceInfo struct {
	Price     float64      `json:"price"`
	Timestamp int64        `json:"timestamp"`
	Infos     []WeightInfo `json:"infos"`
}

//@param bAverage     get average not cointain lowest and highest
//@return bool     symbol find?
func partyPrice(infos []conf.PriceInfo, symbol string, bAverage bool) (bool, PartyPriceInfo) {
	var symbolPriceInfo []conf.PriceInfo
	for _, info := range infos {
		if strings.ToLower(info.Symbol) == strings.ToLower(symbol) {
			symbolPriceInfo = append(symbolPriceInfo, info)
		}
	}

	infosLen := len(symbolPriceInfo)
	if infosLen == 0 {
		return false, PartyPriceInfo{}
	}

	sort.Slice(symbolPriceInfo, func(i, j int) bool {
		if symbolPriceInfo[i].Price > infos[j].Price {
			return true
		} else {
			return false
		}
	})

	if infosLen > 2 && bAverage {
		symbolPriceInfo = symbolPriceInfo[1 : infosLen-1]
	}

	var partyPriceInfo PartyPriceInfo
	totalPrice := 0.0
	totalWeight := int64(0)
	for _, info := range symbolPriceInfo {
		totalPrice += info.Price * float64(info.Weight)
		totalWeight += info.Weight

		partyPriceInfo.Infos = append(partyPriceInfo.Infos, WeightInfo{Price: info.Price, Weight: info.Weight, ExchangeName: info.PriceOrigin})
	}
	partyPriceInfo.Price = totalPrice / float64(totalWeight)
	partyPriceInfo.Timestamp = symbolPriceInfo[0].TimeStamp

	return true, partyPriceInfo
}

func HandleGetHistoryPrice(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")
	timestampStr, exist := context.GetQuery("timestamp")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	bFind, partyPriceData := getHistoryPrice(symbol, timestamp, true)

	if !bFind {
		log.Println("symbol or exchange not find, symbol:", symbol)
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	response.Data = partyPriceData
	context.JSON(http.StatusOK, response)
}

func getHistoryPrice(symbol string, timestamp int64, bAverage bool) (bool, PartyPriceInfo) {
	//first find in memory
	bMemory := false
	var cacheInfo conf.PriceInfos
	m.RLock()
	//latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	if len(gPriceInfosCache.PriceInfosCache) == 0 {
		//nothing todo
		//just find db
	} else {
		for i := len(gPriceInfosCache.PriceInfosCache) - 1; i >= 0; i-- {
			info := gPriceInfosCache.PriceInfosCache[i]
			if len(info.PriceInfos) == 0 {
				continue
			}
			if info.PriceInfos[0].TimeStamp < timestamp {
				//use memory
				bMemory = true
				cacheInfo = gPriceInfosCache.PriceInfosCache[i]
			}
		}
	}
	m.RUnlock()

	if bMemory {
		return partyPrice(cacheInfo.PriceInfos, symbol, bAverage)
	}

	dbPriceInfos, err := sql.GetHistoryBySymbolTimestamp(symbol, timestamp)
	if err != nil {
		log.Println(err)
		return false, PartyPriceInfo{}
	}

	return partyPrice(dbPriceInfos, symbol, bAverage)
}
