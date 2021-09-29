package main

import (
	"encoding/json"
	"log"
	"net/http"
	conf "price_api/price_server/config"
	"price_api/price_server/exchange"
	"price_api/price_server/sql"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const MSG_URL_NOT_FIND = "url not find"
const MSG_PRICE_NOT_READY = "price not ready"
const MSG_PARAM_NOT_TRUE = "param not true"
const MSG_GET_ARES_ERROR = "get ares info error"
const MSG_PARSE_PARAM_ERROR = "parse param error"
const MSG_GET_LOG_INFO_ERROR = "get log info error"
const MSG_CHECK_USER_ERROR = "user and password not match"

const (
	ERROR = iota - 1000
	NO_MATCH_FORMAT_ERROR
	PARAM_NOT_TRUE_ERROR
	GET_ARES_INFO_ERROR
	PARSE_PARAM_ERROR
	GET_LOG_INFO_ERROR
	GET_HTTP_ERROR_ERROR
	CHECK_USER_ERROR
	SET_WEIGHT_ERROR
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

	var rspData PRICE_INFO
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

type PriceAllInfo struct {
	Name      string  `json:"name"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

func HandleGetPriceAll(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")

	bFind := false

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

type PartyPrice struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
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

func HandleGetBulkPrices(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	symbols := strings.Split(symbol, "_")

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	m.RUnlock()

	mSymbolPriceInfo := make(map[string]PartyPrice)
	for _, symbol := range symbols {
		bFind, partyPriceData := partyPrice(latestInfos.PriceInfos, symbol, true)
		if !bFind {
			mSymbolPriceInfo[symbol] = PartyPrice{Price: 0, Timestamp: 0}
		} else {
			mSymbolPriceInfo[symbol] = PartyPrice{Price: partyPriceData.Price, Timestamp: partyPriceData.Timestamp}
		}
	}

	response.Data = mSymbolPriceInfo
	context.JSON(http.StatusOK, response)
}

func HandleGetReqConfig(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	type EXCHANGE_WEIGHT_INFO struct {
		Exchange string `json:"exchange"`
		Weight   int64  `json:"weight"`
	}

	data := make(map[string][]EXCHANGE_WEIGHT_INFO)
	for symbol, confList := range gRequestPriceConfs {
		for _, conf := range confList {
			data[symbol] = append(data[symbol], EXCHANGE_WEIGHT_INFO{Exchange: conf.Name, Weight: conf.Weight})
		}
	}

	response.Data = data
	context.JSON(http.StatusOK, response)
}

func HandleGetRequestInfo(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	index, exist := context.GetQuery("index")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		response.Code = PARSE_PARAM_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusOK, response)
		return
	}

	user, exist := context.GetQuery("user")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	password, exist := context.GetQuery("password")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	if user != gCfg.User || password != gCfg.Password {
		response.Code = CHECK_USER_ERROR
		response.Message = MSG_CHECK_USER_ERROR
		context.JSON(http.StatusOK, response)
		return
	}

	logInfos, err := sql.GetLogInfo(idx, int(gCfg.PageSize))
	if err != nil {
		response.Code = GET_LOG_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusOK, response)
		return
	}

	response.Data = logInfos
	context.JSON(http.StatusOK, response)
}

func HandleGetRequestInfoBySymbol(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	index, exist := context.GetQuery("index")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		response.Code = PARSE_PARAM_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusOK, response)
		return
	}

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	logInfos, err := sql.GetLogInfoBySymbol(idx, int(gCfg.PageSize), symbol)
	if err != nil {
		response.Code = GET_LOG_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusOK, response)
		return
	}

	response.Data = parseLogInfos(logInfos, symbol)
	context.JSON(http.StatusOK, response)
}

type PRICE_INFO struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type PRICE_EXCHANGE_INFO struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
	Exchange  string  `json:"exchange"`
}

func parseLogInfos(logInfos []sql.REQ_RSP_LOG_INFO, symbol string) map[string][]interface{} {
	log.Println("loginfos:", logInfos)
	retPriceInfos := make(map[string][]interface{})
	//tag
	for _, logInfo := range logInfos {
		var rsp RESPONSE
		err := json.Unmarshal([]byte(logInfo.Response), &rsp)
		if err != nil {
			log.Println(err)
			continue
		}

		if strings.Contains(logInfo.ReqUrl, "getPriceAll") {
			priceInfoLists := rsp.Data.([]interface{})
			for _, priceInfo := range priceInfoLists {
				info := priceInfo.(map[string]interface{})
				retPriceInfos["getPriceAll"] = append(retPriceInfos["getPriceAll"], PRICE_EXCHANGE_INFO{Price: info["price"].(float64),
					Exchange: info["name"].(string), Timestamp: int64(info["timestamp"].(float64))})
			}
		} else if strings.Contains(logInfo.ReqUrl, "getPrice") {
			mapPriceInfo := rsp.Data.(map[string]interface{})

			retPriceInfos["getPrice"] = append(retPriceInfos["getPrice"], PRICE_INFO{Price: mapPriceInfo["price"].(float64), Timestamp: int64(mapPriceInfo["timestamp"].(float64))})
			//retPriceInfos["getPrice"] = append(retPriceInfos["getPrice"], rsp.Data.(PRICE_INFO))
		} else if strings.Contains(logInfo.ReqUrl, "getPartyPrice") {
			mapPriceInfo := rsp.Data.(map[string]interface{})
			retPriceInfos["getPartyPrice"] = append(retPriceInfos["getPartyPrice"], PRICE_INFO{Price: mapPriceInfo["price"].(float64), Timestamp: int64(mapPriceInfo["timestamp"].(float64))})
		} else if strings.Contains(logInfo.ReqUrl, "getHistoryPrice") {
			mapPriceInfo := rsp.Data.(map[string]interface{})
			timestamp := int64(mapPriceInfo["timestamp"].(float64))
			retPriceInfos["getHistoryPrice"] = append(retPriceInfos["getHistoryPrice"], PRICE_EXCHANGE_INFO{Price: mapPriceInfo["price"].(float64), Timestamp: timestamp})

			priceInfoLists := mapPriceInfo["infos"].([]interface{})
			for _, priceInfo := range priceInfoLists {
				info := priceInfo.(map[string]interface{})
				retPriceInfos["getHistoryPrice"] = append(retPriceInfos["getHistoryPrice"], PRICE_EXCHANGE_INFO{Price: info["price"].(float64),
					Exchange: info["exchangeName"].(string), Timestamp: timestamp})
			}
		} else if strings.Contains(logInfo.ReqUrl, "getBulkPrices") {
			mapPriceInfo := rsp.Data.(map[string]interface{})
			symbolPriceInfo := mapPriceInfo[symbol].(map[string]interface{})
			retPriceInfos["getBulkPrices"] = append(retPriceInfos["getBulkPrices"], PRICE_INFO{Price: symbolPriceInfo["price"].(float64), Timestamp: int64(symbolPriceInfo["timestamp"].(float64))})
		} else {
			//log.Println("unknow logInfo", logInfo)
			continue
		}
	}
	return retPriceInfos
}

func HandleGetHttpErrorInfo(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	index, exist := context.GetQuery("index")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		response.Code = PARSE_PARAM_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusOK, response)
		return
	}

	httpErrorInfos, err := sql.GetHttpErrorInfo(idx, int(gCfg.PageSize))
	if err != nil {
		response.Code = GET_HTTP_ERROR_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusOK, response)
		return
	}

	response.Data = httpErrorInfos
	context.JSON(http.StatusOK, response)
}

func HandleGetLocalPrices(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	index, exist := context.GetQuery("index")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		response.Code = PARSE_PARAM_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusOK, response)
		return
	}

	start := idx * int(gCfg.PageSize)
	end := start + int(gCfg.PageSize)

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	tmpRetData := conf.PriceInfosCache{}
	m.RLock()
	if start < len(gPriceInfosCache.PriceInfosCache) {
		if end < len(gPriceInfosCache.PriceInfosCache) {
			tmpRetData.PriceInfosCache = gPriceInfosCache.PriceInfosCache[start:end]
		} else {
			tmpRetData.PriceInfosCache = gPriceInfosCache.PriceInfosCache[start:]
		}
	}
	m.RUnlock()

	retData := conf.PriceInfosCache{}
	for _, infosCache := range tmpRetData.PriceInfosCache {
		var retPriceInfos conf.PriceInfos
		for _, priceInfo := range infosCache.PriceInfos {
			if priceInfo.Symbol == symbol {
				retPriceInfos.PriceInfos = append(retPriceInfos.PriceInfos, priceInfo)
			}
		}
		if len(retPriceInfos.PriceInfos) != 0 {
			retData.PriceInfosCache = append(retData.PriceInfosCache, retPriceInfos)
		}
	}

	response.Data = retData
	context.JSON(http.StatusOK, response)
}

func HandleSetWeight(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	weightStr, exist := context.GetQuery("weight")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	weight, err := strconv.Atoi(weightStr)
	if err != nil {
		response.Code = PARSE_PARAM_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusOK, response)
		return
	}

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	user, exist := context.GetQuery("user")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	password, exist := context.GetQuery("password")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	if user != gCfg.User || password != gCfg.Password {
		response.Code = CHECK_USER_ERROR
		response.Message = MSG_CHECK_USER_ERROR
		context.JSON(http.StatusOK, response)
		return
	}

	exchange, exist := context.GetQuery("exchange")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusOK, response)
		return
	}

	err = sql.SetWeight(symbol, exchange, weight)
	if err != nil {
		response.Code = SET_WEIGHT_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusOK, response)
		return
	}

	m.Lock()
	for i, conf := range gRequestPriceConfs[symbol] {
		if conf.Name == exchange {
			gRequestPriceConfs[symbol][i].Weight = int64(weight)
			break
		}
	}
	m.Unlock()

	context.JSON(http.StatusOK, response)
}

func HandleGetAresAll(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}
	aresShowInfo, err := exchange.GetAresInfo(gCfg.Proxy)
	if err != nil {
		response.Code = GET_ARES_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusOK, response)
	}

	response.Data = aresShowInfo
	context.JSON(http.StatusOK, response)
}
