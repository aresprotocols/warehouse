package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	conf "price_api/price_server/config"
	"price_api/price_server/exchange"
	"price_api/price_server/jwt"
	"price_api/price_server/sql"
	"price_api/price_server/util"
	"sort"
	"strconv"
	"strings"
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

var (
	handle *Handle
)

type Handle struct {
	fetcher *exchange.Fetcher
}

func InitHandle(cfg conf.Config) *Handle {
	handle = &Handle{
		fetcher: exchange.InitFetcher(cfg),
	}

	handle.fetcher.Start()

	return handle
}

func (h *Handle) Stop() {
	h.fetcher.Stop()
}

func HandleGetPrice(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	lastIndex := strings.LastIndex(context.Param("name")[1:], "/")
	if lastIndex == -1 {
		logger.Infoln("not true param name", context.Param("name")[1:])
		response.Code = NO_MATCH_FORMAT_ERROR
		response.Message = MSG_URL_NOT_FIND
		context.JSON(http.StatusBadRequest, response)
		return
	}

	symbol := context.Param("name")[1 : lastIndex+1]
	exchange := context.Param("name")[lastIndex+2:]

	var rspData PRICE_INFO
	bFind := false

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	for _, info := range latestInfos.PriceInfos {
		if strings.EqualFold(info.Symbol, symbol) &&
			strings.EqualFold(info.PriceOrigin, exchange) {
			bFind = true
			rspData.Price = info.Price
			rspData.Timestamp = info.TimeStamp
		}
	}
	m.RUnlock()

	if !bFind {
		logger.Infoln("symbol or exchange not find, symbol:", symbol, " exchange:", exchange)
		response.Code = NO_MATCH_FORMAT_ERROR
		response.Message = MSG_URL_NOT_FIND
		context.JSON(http.StatusNotFound, response)
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
		logger.Infoln("symbol or exchange not find, symbol:", symbol)
		response.Code = NO_MATCH_FORMAT_ERROR
		response.Message = MSG_URL_NOT_FIND
		context.JSON(http.StatusNotFound, response)
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
	Weight    int64   `json:"weight"`
}

func HandleGetPriceAll(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")

	bFind := false

	var priceAll []PriceAllInfo

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]

	for _, info := range latestInfos.PriceInfos {
		if strings.EqualFold(info.Symbol, symbol) {
			bFind = true
			priceAllInfo := PriceAllInfo{Name: info.PriceOrigin,
				Symbol:    info.Symbol,
				Price:     info.Price,
				Timestamp: info.TimeStamp,
				Weight:    info.Weight,
			}
			priceAll = append(priceAll, priceAllInfo)
		}
	}
	m.RUnlock()

	if !bFind {
		logger.Infoln("symbol or exchange not find, symbol:", symbol)
		response.Code = NO_MATCH_FORMAT_ERROR
		response.Message = MSG_URL_NOT_FIND
		context.JSON(http.StatusNotFound, response)
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

//@param bAverage     get average not cointain lowest and highest
//@return bool     symbol find?
func partyPrice(infos []conf.PriceInfo, symbol string, bAverage bool) (bool, PartyPriceInfo) {
	var symbolPriceInfo []conf.PriceInfo
	for _, info := range infos {
		if strings.EqualFold(info.Symbol, symbol) {
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
		context.JSON(http.StatusBadRequest, response)
		return
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	bFind, partyPriceData := getHistoryPrice(symbol, timestamp, true)

	if !bFind {
		logger.Infoln("symbol or exchange not find, symbol:", symbol)
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusNotFound, response)
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
	if len(gPriceInfosCache.PriceInfosCache) != 0 {
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
		logger.WithError(err).Errorf("get history by symbol timestamp error,symbol:%s", symbol)
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
		context.JSON(http.StatusBadRequest, response)
		return
	}

	symbols := strings.Split(symbol, "_")

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	m.RUnlock()

	mSymbolPriceInfo := make(map[string]PRICE_INFO)
	for _, symbol := range symbols {
		bFind, partyPriceData := partyPrice(latestInfos.PriceInfos, symbol, true)
		if !bFind {
			mSymbolPriceInfo[symbol] = PRICE_INFO{Price: 0, Timestamp: 0}
		} else {
			mSymbolPriceInfo[symbol] = PRICE_INFO{Price: partyPriceData.Price, Timestamp: partyPriceData.Timestamp}
		}
	}

	response.Data = mSymbolPriceInfo
	context.JSON(http.StatusOK, response)
}

func HandleGetBulkCurrencyPrices(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}
	currency, exist := context.GetQuery("currency")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	symbols := strings.Split(symbol, "_")

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	m.RUnlock()

	mSymbolPriceInfo := make(map[string]PRICE_INFO)
	for _, symbol := range symbols {
		token := symbol + currency
		bFind, partyPriceData := partyPrice(latestInfos.PriceInfos, token, true)
		if !bFind {
			mSymbolPriceInfo[token] = PRICE_INFO{Price: 0, Timestamp: 0}
		} else {
			mSymbolPriceInfo[token] = PRICE_INFO{Price: partyPriceData.Price, Timestamp: partyPriceData.Timestamp}
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
		context.JSON(http.StatusBadRequest, response)
		return
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		response.Code = PARSE_PARAM_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusBadRequest, response)
		return
	}

	logInfos, err := sql.GetLogInfo(idx, int(gCfg.PageSize))
	if err != nil {
		response.Code = GET_LOG_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusBadRequest, response)
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
		context.JSON(http.StatusBadRequest, response)
		return
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		response.Code = PARSE_PARAM_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusBadRequest, response)
		return
	}

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	logInfos, err := sql.GetLogInfoBySymbol(idx, int(gCfg.PageSize), symbol)
	if err != nil {
		logger.WithError(err).Errorf("get log info by symbol occur error,symbol:%s,index:%d", symbol, idx)
		response.Code = GET_LOG_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}
	total, err := sql.GetTotalLogInfoBySymbol(symbol)
	if err != nil {
		logger.WithError(err).Errorf("get total log info by symbol occur error,symbol:%s", symbol)
		response.Code = GET_LOG_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	items := parseLogInfos(logInfos, symbol)
	response.Data = Pagination{
		CurPage:  idx,
		TotalNum: total,
		Items:    items,
	}
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
	Weight    int64   `json:"weight"`
}

type PRICE_EXCHANGE_WEIGHT_INFO struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
	Exchange  string  `json:"exchange"`
	Weight    int     `json:"weight"`
}

type CLIENT_INFO struct {
	Ip               string `json:"ip"`
	RequestTime      string `json:"request_time"`
	RequestTimestamp int64  `json:"request_timestamp"`
}

type CLIENT_PRICE_INFO struct {
	Client    CLIENT_INFO `json:"client"`
	PriceInfo PRICE_INFO  `json:"price_info"`
}

type CLIENT_PRICEALL_INFO struct {
	Client     CLIENT_INFO           `json:"client"`
	PriceInfos []PRICE_EXCHANGE_INFO `json:"price_infos"`
}

type PARTY_PRICE_INFO struct {
	Type       string                       `json:"type"`
	Client     CLIENT_INFO                  `json:"client"`
	PriceInfo  PRICE_INFO                   `json:"price_info"`
	PriceInfos []PRICE_EXCHANGE_WEIGHT_INFO `json:"price_infos"`
}

func parseLogInfos(logInfos []sql.REQ_RSP_LOG_INFO, symbol string) []PARTY_PRICE_INFO {
	//retPriceInfos := make(map[string][]interface{})
	retPriceInfos := make([]PARTY_PRICE_INFO, 0)

	for _, logInfo := range logInfos {
		var rsp RESPONSE
		err := json.Unmarshal([]byte(logInfo.Response), &rsp)
		if err != nil {
			logger.WithError(err).Errorf("unmarshal logInfo response occur error")
			continue
		}

		if strings.Contains(logInfo.ReqUrl, "getPriceAll") {
			var historyPriceInfo PARTY_PRICE_INFO
			historyPriceInfo.Type = "getPriceAll"

			priceInfoLists := rsp.Data.([]interface{})

			historyPriceInfo.Client = CLIENT_INFO{Ip: logInfo.Ip, RequestTime: logInfo.RequestTime, RequestTimestamp: logInfo.RequestTimestamp}

			for index, priceInfo := range priceInfoLists {
				info := priceInfo.(map[string]interface{})
				if index == 0 {
					historyPriceInfo.PriceInfo = PRICE_INFO{Price: info["price"].(float64), Timestamp: int64(info["timestamp"].(float64))}
				}
				weightExchangeInfo := PRICE_EXCHANGE_WEIGHT_INFO{Price: info["price"].(float64),
					Exchange: info["name"].(string), Timestamp: int64(info["timestamp"].(float64)), Weight: int(info["weight"].(float64))}
				historyPriceInfo.PriceInfos = append(historyPriceInfo.PriceInfos, weightExchangeInfo)

			}

			retPriceInfos = append(retPriceInfos, historyPriceInfo)

		} else if strings.Contains(logInfo.ReqUrl, "getPrice") {
			var historyPriceInfo PARTY_PRICE_INFO
			historyPriceInfo.Type = "getPrice"

			mapPriceInfo := rsp.Data.(map[string]interface{})

			historyPriceInfo.Client = CLIENT_INFO{Ip: logInfo.Ip, RequestTime: logInfo.RequestTime, RequestTimestamp: logInfo.RequestTimestamp}
			historyPriceInfo.PriceInfo = PRICE_INFO{Price: mapPriceInfo["price"].(float64), Timestamp: int64(mapPriceInfo["timestamp"].(float64))}
			historyPriceInfo.PriceInfos = make([]PRICE_EXCHANGE_WEIGHT_INFO, 0)

			retPriceInfos = append(retPriceInfos, historyPriceInfo)
		} else if strings.Contains(logInfo.ReqUrl, "getPartyPrice") {
			mapPriceInfo := rsp.Data.(map[string]interface{})
			var historyPriceInfo PARTY_PRICE_INFO

			historyPriceInfo.Type = "getPartyPrice"

			timestamp := int64(mapPriceInfo["timestamp"].(float64))
			historyPriceInfo.Client = CLIENT_INFO{Ip: logInfo.Ip, RequestTime: logInfo.RequestTime, RequestTimestamp: logInfo.RequestTimestamp}
			historyPriceInfo.PriceInfo = PRICE_INFO{Price: mapPriceInfo["price"].(float64), Timestamp: timestamp}

			priceInfoLists := mapPriceInfo["infos"].([]interface{})
			for _, priceInfo := range priceInfoLists {
				info := priceInfo.(map[string]interface{})
				weightExchangeInfo := PRICE_EXCHANGE_WEIGHT_INFO{Price: info["price"].(float64),
					Exchange: info["exchangeName"].(string), Timestamp: timestamp, Weight: int(info["weight"].(float64))}
				historyPriceInfo.PriceInfos = append(historyPriceInfo.PriceInfos, weightExchangeInfo)
			}

			retPriceInfos = append(retPriceInfos, historyPriceInfo)
		} else if strings.Contains(logInfo.ReqUrl, "getHistoryPrice") {
			mapPriceInfo := rsp.Data.(map[string]interface{})
			var historyPriceInfo PARTY_PRICE_INFO

			historyPriceInfo.Type = "getHistoryPrice"

			timestamp := int64(mapPriceInfo["timestamp"].(float64))
			historyPriceInfo.Client = CLIENT_INFO{Ip: logInfo.Ip, RequestTime: logInfo.RequestTime, RequestTimestamp: logInfo.RequestTimestamp}
			historyPriceInfo.PriceInfo = PRICE_INFO{Price: mapPriceInfo["price"].(float64), Timestamp: timestamp}

			priceInfoLists := mapPriceInfo["infos"].([]interface{})
			for _, priceInfo := range priceInfoLists {
				info := priceInfo.(map[string]interface{})
				weightExchangeInfo := PRICE_EXCHANGE_WEIGHT_INFO{Price: info["price"].(float64),
					Exchange: info["exchangeName"].(string), Timestamp: timestamp, Weight: int(info["weight"].(float64))}
				historyPriceInfo.PriceInfos = append(historyPriceInfo.PriceInfos, weightExchangeInfo)
			}
			retPriceInfos = append(retPriceInfos, historyPriceInfo)
		} else if strings.Contains(logInfo.ReqUrl, "getBulkPrices") {
			var historyPriceInfo PARTY_PRICE_INFO
			historyPriceInfo.Type = "getBulkPrices"

			mapPriceInfo := rsp.Data.(map[string]interface{})
			symbolPriceInfo := mapPriceInfo[symbol].(map[string]interface{})

			historyPriceInfo.Client = CLIENT_INFO{Ip: logInfo.Ip, RequestTime: logInfo.RequestTime}
			historyPriceInfo.PriceInfo = PRICE_INFO{Price: symbolPriceInfo["price"].(float64), Timestamp: int64(symbolPriceInfo["timestamp"].(float64))}
			historyPriceInfo.PriceInfos = make([]PRICE_EXCHANGE_WEIGHT_INFO, 0)
			retPriceInfos = append(retPriceInfos, historyPriceInfo)
		} else {
			logger.Infoln("unknow logInfo", logInfo)
			continue
		}
	}
	return retPriceInfos
}

func HandleGetHttpErrorInfo(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")
	symbol = strings.ToLower(symbol)
	if !strings.Contains(symbol, "-") {
		if strings.HasSuffix(symbol, "usdt") {
			symbol = strings.ReplaceAll(symbol, "usdt", "-usdt")
		}
	}

	index, exist := context.GetQuery("index")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		response.Code = PARSE_PARAM_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusBadRequest, response)
		return
	}

	total, err := sql.GetTotalHttpErrorInfo(symbol)
	if err != nil {
		response.Code = GET_HTTP_ERROR_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	httpErrorInfos, err := sql.GetHttpErrorInfo(idx, symbol, int(gCfg.PageSize))
	if err != nil {
		response.Code = GET_HTTP_ERROR_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}
	response.Data = Pagination{
		CurPage:  idx,
		TotalNum: total,
		Items:    httpErrorInfos,
	}
	context.JSON(http.StatusOK, response)
}

func HandleGetLocalPrices(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	index, exist := context.GetQuery("index")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		response.Code = PARSE_PARAM_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusBadRequest, response)
		return
	}

	start := idx * int(gCfg.PageSize)
	end := start + int(gCfg.PageSize)

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
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

func HandleGetUpdatePriceHistory(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	index, exist := context.GetQuery("index")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	idx, err := strconv.Atoi(index)
	if err != nil {
		response.Code = PARSE_PARAM_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusBadRequest, response)
		return
	}

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	infos, err := sql.GetHistoryBySymbol(idx, int(gCfg.PageSize), symbol)
	if err != nil {
		logger.WithError(err).Errorf("get history by symbol occur error,symbol:%s,index:%d", symbol, idx)
		response.Code = GET_LOG_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}
	total, err := sql.GetTotalHistoryBySymbol(symbol)
	if err != nil {
		logger.WithError(err).Errorf("get total history by symbol occur error,symbol:%s", symbol)
		response.Code = GET_LOG_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = Pagination{
		CurPage:  idx,
		TotalNum: total,
		Items:    infos,
	}
	context.JSON(http.StatusOK, response)
}

type SetWeightReq struct {
	Weight   int    `json:"weight"`
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}

func HandleSetWeight(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	var setWeightReq SetWeightReq
	err := context.ShouldBind(&setWeightReq)

	if len(setWeightReq.Symbol) == 0 {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if len(setWeightReq.Exchange) == 0 {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	err = sql.SetWeight(setWeightReq.Symbol, setWeightReq.Exchange, setWeightReq.Weight)
	if err != nil {
		response.Code = SET_WEIGHT_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	m.Lock()
	for i, conf := range gRequestPriceConfs[setWeightReq.Symbol] {
		if conf.Name == setWeightReq.Exchange {
			gRequestPriceConfs[setWeightReq.Symbol][i].Weight = int64(setWeightReq.Weight)
			break
		}
	}
	m.Unlock()

	context.JSON(http.StatusOK, response)
}

type HEARTBEAT_INFO struct {
	ExpectResources int   `json:"expect_resources"`
	ActualResources int   `json:"actual_resources"`
	LatestTimestamp int64 `json:"latest_timestamp"`
	Interval        int64 `json:"interval"`
}

func HandleGetUpdatePriceHeartbeat(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol := context.Param("symbol")

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	m.RUnlock()

	var symbolPriceInfo = make([]conf.PriceInfo, 0)
	for _, info := range latestInfos.PriceInfos {
		if strings.EqualFold(info.Symbol, symbol) {
			symbolPriceInfo = append(symbolPriceInfo, info)
		}
	}

	if len(symbolPriceInfo) == 0 {
		logger.Infoln("symbol not find, symbol:", symbol)
		response.Code = NO_MATCH_FORMAT_ERROR
		response.Message = MSG_URL_NOT_FIND
		context.JSON(http.StatusNotFound, response)
		return
	}

	tokenSymbol := strings.ReplaceAll(symbol, "usdt", "-usdt")

	exchangeConfs := gRequestPriceConfs[tokenSymbol]

	response.Data = HEARTBEAT_INFO{
		ExpectResources: len(exchangeConfs),
		ActualResources: len(symbolPriceInfo),
		LatestTimestamp: symbolPriceInfo[0].TimeStamp,
		Interval:        gCfg.Interval,
	}
	context.JSON(http.StatusOK, response)
}

func HandleGetBulkSymbolsState(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	symbol, exist := context.GetQuery("symbol")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}
	currency, exist := context.GetQuery("currency")
	if !exist {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	symbols := strings.Split(symbol, "_")

	m.RLock()
	latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	m.RUnlock()

	mSymbolState := make(map[string]bool)
	for _, symbol := range symbols {
		token := symbol + currency
		var symbolPriceInfo = make([]conf.PriceInfo, 0)
		for _, info := range latestInfos.PriceInfos {
			if strings.EqualFold(info.Symbol, token) {
				symbolPriceInfo = append(symbolPriceInfo, info)
			}
		}
		actualResourcesLens := len(symbolPriceInfo)

		tokenSymbol := symbol + "-" + currency
		exchangeConfs := gRequestPriceConfs[tokenSymbol]
		expectResourcesLens := len(exchangeConfs)

		mSymbolState[token] = actualResourcesLens > expectResourcesLens/2
	}

	response.Data = mSymbolState
	context.JSON(http.StatusOK, response)
}

func HandleGetAresAll(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	aresShowInfo, err := exchange.GetGateAresInfo(gCfg.Proxy)
	aresShowInfo.Rank = handle.fetcher.GetCMCInfo().Rank

	if err != nil {
		logger.WithError(err).Errorf("get gate ares info occur error")
		response.Code = GET_ARES_INFO_ERROR
		response.Message = err.Error()
		context.JSON(http.StatusBadRequest, response)
	}

	response.Data = aresShowInfo
	context.JSON(http.StatusOK, response)
}

func HandleGetDexPrice(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	aresShowInfo := handle.fetcher.GetDexPrice()

	response.Data = aresShowInfo
	context.JSON(http.StatusOK, response)
}

type AdminUser struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func HandleAuth(context *gin.Context) {
	response := RESPONSE{Code: 0, Message: "OK"}

	var user AdminUser
	err := context.ShouldBind(&user)
	if err != nil {
		logger.WithError(err).Errorf("bind user occur error")
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if len(user.User) == 0 {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if len(user.Password) == 0 {
		response.Code = PARAM_NOT_TRUE_ERROR
		response.Message = MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	md5Password := util.Md5Str(gCfg.Password)

	if user.User != gCfg.User || user.Password != md5Password {
		response.Code = CHECK_USER_ERROR
		response.Message = MSG_CHECK_USER_ERROR
		context.JSON(http.StatusUnauthorized, response)
		return
	}
	authToken, err := jwt.GenToken(user.User, []byte(gCfg.Password))
	if err != nil {
		logger.WithError(err).Error("generate jwt token occur error")
		response.Code = ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
	}
	response.Data = authToken
	context.JSON(http.StatusOK, response)

}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		response := RESPONSE{Code: 0, Message: "OK"}

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Code = CHECK_USER_ERROR
			response.Message = MSG_CHECK_USER_ERROR
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Code = CHECK_USER_ERROR
			response.Message = MSG_CHECK_USER_ERROR
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		mc, err := jwt.ParseToken(parts[1], []byte(gCfg.Password))
		if err != nil {
			response.Code = CHECK_USER_ERROR
			response.Message = MSG_CHECK_USER_ERROR
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		c.Set("username", mc.Username)
		c.Next()
	}
}
