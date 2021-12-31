package service

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"price_api/price_server/internal/repository"
	"price_api/price_server/internal/vo"
	"strings"
)

type RequestInfoService struct {
	logInfoRepo repository.LogInfoRepository
}

func newRequestInfo(svc *service) *RequestInfoService {
	return &RequestInfoService{logInfoRepo: repository.NewLogInfoRepository(svc.db)}
}

func (s *RequestInfoService) GetLogInfos(idx, pageSize int) (vo.LOG_INFOS, error) {

	return s.logInfoRepo.GetLogInfo(idx, pageSize)
}

func (s *RequestInfoService) GetRequestInfoBySymbol(idx, pageSize int, symbol string, ip string) (int, []vo.PARTY_PRICE_INFO, error) {
	logInfos, err := s.logInfoRepo.GetLogInfoBySymbol(idx, pageSize, symbol, ip)
	if err != nil {
		logger.WithError(err).Errorf("get log info by symbol occur error,symbol:%s,index:%d", symbol, idx)
		return 0, nil, err
	}
	total, err := s.logInfoRepo.GetTotalLogInfoBySymbol(symbol, ip)
	if err != nil {
		logger.WithError(err).Errorf("get total log info by symbol occur error,symbol:%s", symbol)
		return 0, nil, err
	}
	items := s.parseLogInfos(logInfos, symbol)
	return total, items, nil
}

func (s *RequestInfoService) parseLogInfos(logInfos []vo.REQ_RSP_LOG_INFO, symbol string) []vo.PARTY_PRICE_INFO {
	//retPriceInfos := make(map[string][]interface{})
	retPriceInfos := make([]vo.PARTY_PRICE_INFO, 0)

	for _, logInfo := range logInfos {
		var rsp vo.RESPONSE
		err := json.Unmarshal([]byte(logInfo.Response), &rsp)
		if err != nil {
			logger.WithError(err).Errorf("unmarshal logInfo response occur error")
			continue
		}

		if strings.Contains(logInfo.ReqUrl, "getPriceAll") {
			var historyPriceInfo vo.PARTY_PRICE_INFO
			historyPriceInfo.Type = "getPriceAll"

			priceInfoLists := rsp.Data.([]interface{})

			historyPriceInfo.Client = vo.CLIENT_INFO{Ip: logInfo.Ip, RequestTime: logInfo.RequestTime, RequestTimestamp: logInfo.RequestTimestamp}

			for index, priceInfo := range priceInfoLists {
				info := priceInfo.(map[string]interface{})
				if index == 0 {
					historyPriceInfo.PriceInfo = vo.PRICE_INFO{Price: info["price"].(float64), Timestamp: int64(info["timestamp"].(float64))}
				}
				weightExchangeInfo := vo.PRICE_EXCHANGE_WEIGHT_INFO{Price: info["price"].(float64),
					Exchange: info["name"].(string), Timestamp: int64(info["timestamp"].(float64)), Weight: int(info["weight"].(float64))}
				historyPriceInfo.PriceInfos = append(historyPriceInfo.PriceInfos, weightExchangeInfo)

			}

			retPriceInfos = append(retPriceInfos, historyPriceInfo)

		} else if strings.Contains(logInfo.ReqUrl, "getPrice") {
			var historyPriceInfo vo.PARTY_PRICE_INFO
			historyPriceInfo.Type = "getPrice"

			mapPriceInfo := rsp.Data.(map[string]interface{})

			historyPriceInfo.Client = vo.CLIENT_INFO{Ip: logInfo.Ip, RequestTime: logInfo.RequestTime, RequestTimestamp: logInfo.RequestTimestamp}
			historyPriceInfo.PriceInfo = vo.PRICE_INFO{Price: mapPriceInfo["price"].(float64), Timestamp: int64(mapPriceInfo["timestamp"].(float64))}
			historyPriceInfo.PriceInfos = make([]vo.PRICE_EXCHANGE_WEIGHT_INFO, 0)

			retPriceInfos = append(retPriceInfos, historyPriceInfo)
		} else if strings.Contains(logInfo.ReqUrl, "getPartyPrice") {
			mapPriceInfo := rsp.Data.(map[string]interface{})
			var historyPriceInfo vo.PARTY_PRICE_INFO

			historyPriceInfo.Type = "getPartyPrice"

			timestamp := int64(mapPriceInfo["timestamp"].(float64))
			historyPriceInfo.Client = vo.CLIENT_INFO{Ip: logInfo.Ip, RequestTime: logInfo.RequestTime, RequestTimestamp: logInfo.RequestTimestamp}
			historyPriceInfo.PriceInfo = vo.PRICE_INFO{Price: mapPriceInfo["price"].(float64), Timestamp: timestamp}

			priceInfoLists := mapPriceInfo["infos"].([]interface{})
			for _, priceInfo := range priceInfoLists {
				info := priceInfo.(map[string]interface{})
				weightExchangeInfo := vo.PRICE_EXCHANGE_WEIGHT_INFO{Price: info["price"].(float64),
					Exchange: info["exchangeName"].(string), Timestamp: timestamp, Weight: int(info["weight"].(float64))}
				historyPriceInfo.PriceInfos = append(historyPriceInfo.PriceInfos, weightExchangeInfo)
			}

			retPriceInfos = append(retPriceInfos, historyPriceInfo)
		} else if strings.Contains(logInfo.ReqUrl, "getHistoryPrice") {
			mapPriceInfo := rsp.Data.(map[string]interface{})
			var historyPriceInfo vo.PARTY_PRICE_INFO

			historyPriceInfo.Type = "getHistoryPrice"

			timestamp := int64(mapPriceInfo["timestamp"].(float64))
			historyPriceInfo.Client = vo.CLIENT_INFO{Ip: logInfo.Ip, RequestTime: logInfo.RequestTime, RequestTimestamp: logInfo.RequestTimestamp}
			historyPriceInfo.PriceInfo = vo.PRICE_INFO{Price: mapPriceInfo["price"].(float64), Timestamp: timestamp}

			priceInfoLists := mapPriceInfo["infos"].([]interface{})
			for _, priceInfo := range priceInfoLists {
				info := priceInfo.(map[string]interface{})
				weightExchangeInfo := vo.PRICE_EXCHANGE_WEIGHT_INFO{Price: info["price"].(float64),
					Exchange: info["exchangeName"].(string), Timestamp: timestamp, Weight: int(info["weight"].(float64))}
				historyPriceInfo.PriceInfos = append(historyPriceInfo.PriceInfos, weightExchangeInfo)
			}
			retPriceInfos = append(retPriceInfos, historyPriceInfo)
		} else if strings.Contains(logInfo.ReqUrl, "getBulkPrices") {
			var historyPriceInfo vo.PARTY_PRICE_INFO
			historyPriceInfo.Type = "getBulkPrices"

			mapPriceInfo := rsp.Data.(map[string]interface{})
			symbolPriceInfo := mapPriceInfo[symbol].(map[string]interface{})

			historyPriceInfo.Client = vo.CLIENT_INFO{Ip: logInfo.Ip, RequestTime: logInfo.RequestTime, RequestTimestamp: logInfo.RequestTimestamp}
			historyPriceInfo.PriceInfo = vo.PRICE_INFO{Price: symbolPriceInfo["price"].(float64), Timestamp: int64(symbolPriceInfo["timestamp"].(float64))}
			historyPriceInfo.PriceInfos = make([]vo.PRICE_EXCHANGE_WEIGHT_INFO, 0)
			retPriceInfos = append(retPriceInfos, historyPriceInfo)
		} else if strings.Contains(logInfo.ReqUrl, "getBulkCurrencyPrices") {
			var historyPriceInfo vo.PARTY_PRICE_INFO
			historyPriceInfo.Type = "getBulkCurrencyPrices"

			mapPriceInfo := rsp.Data.(map[string]interface{})
			symbolPriceInfo := mapPriceInfo[symbol].(map[string]interface{})

			timestamp := int64(symbolPriceInfo["timestamp"].(float64))

			historyPriceInfo.Client = vo.CLIENT_INFO{Ip: logInfo.Ip, RequestTime: logInfo.RequestTime, RequestTimestamp: logInfo.RequestTimestamp}
			historyPriceInfo.PriceInfo = vo.PRICE_INFO{Price: symbolPriceInfo["price"].(float64), Timestamp: int64(symbolPriceInfo["timestamp"].(float64))}

			priceInfoListsValue := symbolPriceInfo["infos"]
			if priceInfoListsValue == nil {
				historyPriceInfo.PriceInfos = make([]vo.PRICE_EXCHANGE_WEIGHT_INFO, 0)
			} else {
				priceInfoLists := priceInfoListsValue.([]interface{})
				for _, priceInfo := range priceInfoLists {
					info := priceInfo.(map[string]interface{})
					weightExchangeInfo := vo.PRICE_EXCHANGE_WEIGHT_INFO{Price: info["price"].(float64),
						Exchange: info["exchangeName"].(string), Timestamp: timestamp, Weight: int(info["weight"].(float64))}
					historyPriceInfo.PriceInfos = append(historyPriceInfo.PriceInfos, weightExchangeInfo)
				}
			}
			retPriceInfos = append(retPriceInfos, historyPriceInfo)
		} else {
			logger.Infoln("unknow logInfo", logInfo)
			continue
		}
	}
	return retPriceInfos
}

func (s *RequestInfoService) InsertLogInfo(mapInfo map[string]interface{}, t int) error {
	return s.logInfoRepo.InsertLogInfo(mapInfo, t)
}
