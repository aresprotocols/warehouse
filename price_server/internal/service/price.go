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
	mSymbolPriceInfo := make(map[string]vo.PartyPriceInfo)
	for _, symbolTemp := range symbols {
		token := symbolTemp + currency
		latestInfos := s.gPriceInfosCache.GetLatestPriceInfos(token)
		bFind, partyPriceData := util.PartyPrice(latestInfos.PriceInfos, true)
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
	mSymbolPriceInfo := make(map[string]vo.PRICE_INFO)
	for _, symbolTemp := range symbols {
		latestInfos := s.gPriceInfosCache.GetLatestPriceInfos(symbolTemp)
		bFind, partyPriceData := util.PartyPrice(latestInfos.PriceInfos, true)
		if !bFind {
			mSymbolPriceInfo[symbolTemp] = vo.PRICE_INFO{Price: 0, Timestamp: 0}
		} else {
			mSymbolPriceInfo[symbolTemp] = vo.PRICE_INFO{Price: partyPriceData.Price, Timestamp: partyPriceData.Timestamp}
		}
	}
	return mSymbolPriceInfo
}

func (s *PriceService) GetBulkSymbolsState(symbolStr string, currency string) map[string]bool {
	symbols := strings.Split(symbolStr, "_")

	mSymbolState := make(map[string]bool)
	for _, symbol := range symbols {
		token := symbol + currency
		latestInfos := s.gPriceInfosCache.GetLatestPriceInfos(token)
		var symbolPriceInfo = latestInfos.PriceInfos
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

	bMemory, cacheInfo = s.gPriceInfosCache.GetPriceInfosEqualTimestamp(symbol, timestamp)

	if bMemory {
		return util.PartyPrice(cacheInfo.PriceInfos, bAverage)
	}

	dbPriceInfos, err := s.coinHistoryRepo.GetHistoryBySymbolAndTimestamp(symbol, timestamp)
	if err != nil {
		logger.WithError(err).Errorf("get history by symbol timestamp error,symbol:%s", symbol)
		return false, vo.PartyPriceInfo{}
	}

	return util.PartyPrice(dbPriceInfos, bAverage)

}

func (s *PriceService) GetLocalPrices(start int, end int, symbol string) conf.PriceInfosCache {
	tmpRetData := s.gPriceInfosCache.GetPriceInfosByRange(symbol, start, end)
	return tmpRetData
}

func (s *PriceService) GetPartyPrice(symbol string) (bool, vo.PartyPriceInfo) {
	latestInfos := s.gPriceInfosCache.GetLatestPriceInfos(symbol)
	return util.PartyPrice(latestInfos.PriceInfos, true)
}

func (s *PriceService) GetPrice(symbol, exchange string) (bool, vo.PRICE_INFO) {
	var rspData vo.PRICE_INFO
	bFind := false

	latestInfos := s.gPriceInfosCache.GetLatestPriceInfos(symbol)
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
	latestInfos := s.gPriceInfosCache.GetLatestPriceInfos(symbol)
	for _, info := range latestInfos.PriceInfos {
		bFind = true
		priceAllInfo := vo.PriceAllInfo{Name: info.PriceOrigin,
			Symbol:    info.Symbol,
			Price:     info.Price,
			Timestamp: info.TimeStamp,
			Weight:    info.Weight,
		}
		priceAll = append(priceAll, priceAllInfo)
	}
	return bFind, priceAll
}

func (s *PriceService) GetCacheLength() int {
	return s.gPriceInfosCache.GetCacheLength()
}

func (s *PriceService) UpdateCachePrice(symbol string, infos conf.PriceInfos, maxMemTime int) {
	s.gPriceInfosCache.UpdateCachePrice(symbol, infos, maxMemTime)
}

func (s *PriceService) InsertPriceInfo(cfg conf.PriceInfos) error {
	return s.coinHistoryRepo.InsertPriceInfo(cfg)
}
