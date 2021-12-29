package service

import (
	logger "github.com/sirupsen/logrus"
	conf "price_api/price_server/config"
	"price_api/price_server/internal/cache"
	"price_api/price_server/internal/repository"
	"price_api/price_server/internal/util"
	"price_api/price_server/internal/vo"
	"strings"
)

type PriceService struct {
	gPriceInfosCache   cache.GlobalPriceInfoCache
	coinHistoryRepo    repository.CoinHistoryRepository
	gRequestPriceConfs cache.GlobalRequestPriceConfs
}

func newPrice(svc *service) *PriceService {
	return &PriceService{
		gPriceInfosCache:   svc.globalCache,
		coinHistoryRepo:    repository.NewCoinHistoryRepository(svc.db),
		gRequestPriceConfs: svc.globalRequestPriceConfs,
	}
}

func (s *PriceService) GetBulkCurrencyPrices(symbol string, currency string) map[string]vo.PartyPriceInfo {

	symbols := strings.Split(symbol, "_")

	latestInfos := s.gPriceInfosCache.GetLatestPriceInfos()

	mSymbolPriceInfo := make(map[string]vo.PartyPriceInfo)
	for _, symbolTemp := range symbols {
		token := symbolTemp + currency
		bFind, partyPriceData := util.PartyPrice(latestInfos.PriceInfos, token, true)
		if !bFind {
			mSymbolPriceInfo[token] = partyPriceData
		} else {
			mSymbolPriceInfo[token] = partyPriceData
		}
	}
	return mSymbolPriceInfo
}

func (s *PriceService) GetBulkPrices(symbol string) map[string]vo.PRICE_INFO {
	symbols := strings.Split(symbol, "_")

	latestInfos := s.gPriceInfosCache.GetLatestPriceInfos()

	mSymbolPriceInfo := make(map[string]vo.PRICE_INFO)
	for _, symbol := range symbols {
		bFind, partyPriceData := util.PartyPrice(latestInfos.PriceInfos, symbol, true)
		if !bFind {
			mSymbolPriceInfo[symbol] = vo.PRICE_INFO{Price: 0, Timestamp: 0}
		} else {
			mSymbolPriceInfo[symbol] = vo.PRICE_INFO{Price: partyPriceData.Price, Timestamp: partyPriceData.Timestamp}
		}
	}
	return mSymbolPriceInfo
}

func (s *PriceService) GetBulkSymbolsState(symbolStr string, currency string) map[string]bool {
	symbols := strings.Split(symbolStr, "_")

	latestInfos := s.gPriceInfosCache.GetLatestPriceInfos()

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
		exchangeConfs := s.gRequestPriceConfs.GetConfsBySymbol(tokenSymbol)
		expectResourcesLens := len(exchangeConfs)

		mSymbolState[token] = actualResourcesLens > expectResourcesLens/2
	}

	return mSymbolState

}

func (s *PriceService) GetHistoryPrice(symbol string, timestamp int64, bAverage bool) (bool, vo.PartyPriceInfo) {

	//first find in memory
	bMemory := false
	var cacheInfo conf.PriceInfos

	bMemory, cacheInfo = s.gPriceInfosCache.GetPriceInfosEqualTimestamp(timestamp)

	if bMemory {
		return util.PartyPrice(cacheInfo.PriceInfos, symbol, bAverage)
	}

	dbPriceInfos, err := s.coinHistoryRepo.GetHistoryByTimestamp(timestamp)
	if err != nil {
		logger.WithError(err).Errorf("get history by symbol timestamp error,symbol:%s", symbol)
		return false, vo.PartyPriceInfo{}
	}

	return util.PartyPrice(dbPriceInfos, symbol, bAverage)

}

func (s *PriceService) GetLocalPrices(start int, end int, symbol string) conf.PriceInfosCache {
	tmpRetData := s.gPriceInfosCache.GetPriceInfosByRange(start, end)

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
	return retData
}

func (s *PriceService) GetPartyPrice(symbol string) (bool, vo.PartyPriceInfo) {
	latestInfos := s.gPriceInfosCache.GetLatestPriceInfos()
	return util.PartyPrice(latestInfos.PriceInfos, symbol, true)
}

func (s *PriceService) GetPrice(symbol, exchange string) (bool, vo.PRICE_INFO) {
	var rspData vo.PRICE_INFO
	bFind := false

	latestInfos := s.gPriceInfosCache.GetLatestPriceInfos()
	for _, info := range latestInfos.PriceInfos {
		if strings.EqualFold(info.Symbol, symbol) &&
			strings.EqualFold(info.PriceOrigin, exchange) {
			bFind = true
			rspData.Price = info.Price
			rspData.Timestamp = info.TimeStamp
		}
	}
	return bFind, rspData
}

func (s *PriceService) GetPriceAll(symbol string) (bool, []vo.PriceAllInfo) {
	bFind := false
	var priceAll []vo.PriceAllInfo
	latestInfos := s.gPriceInfosCache.GetLatestPriceInfos()
	for _, info := range latestInfos.PriceInfos {
		if strings.EqualFold(info.Symbol, symbol) {
			bFind = true
			priceAllInfo := vo.PriceAllInfo{Name: info.PriceOrigin,
				Symbol:    info.Symbol,
				Price:     info.Price,
				Timestamp: info.TimeStamp,
				Weight:    info.Weight,
			}
			priceAll = append(priceAll, priceAllInfo)
		}
	}
	return bFind, priceAll
}

func (s *PriceService) GetCacheLength() int {
	return s.gPriceInfosCache.GetCacheLength()
}

func (s *PriceService) UpdateCachePrice(infos conf.PriceInfos, maxMemTime int) {
	s.gPriceInfosCache.UpdateCachePrice(infos, maxMemTime)
}

func (s *PriceService) InsertPriceInfo(cfg conf.PriceInfos) error {
	return s.coinHistoryRepo.InsertPriceInfo(cfg)
}
