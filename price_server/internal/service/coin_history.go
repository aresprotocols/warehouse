package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	conf "price_api/price_server/config"
	"price_api/price_server/internal/cache"
	"price_api/price_server/internal/repository"
	"price_api/price_server/internal/util"
	"price_api/price_server/internal/vo"
	"strings"
)

type CoinHistoryService struct {
	gPriceInfosCache   cache.GlobalPriceInfoCache
	gReqeustPriceConfs cache.GlobalRequestPriceConfs
	updatePriceRepo    repository.UpdatePriceRepository
	coinHistoryRepo    repository.CoinHistoryRepository
}

func newCoinHistory(svc *service) *CoinHistoryService {
	return &CoinHistoryService{
		gPriceInfosCache:   svc.globalCache,
		gReqeustPriceConfs: svc.globalRequestPriceConfs,
		updatePriceRepo:    repository.UpdatePriceRepository{DB: svc.db},
		coinHistoryRepo:    repository.NewCoinHistoryRepository(svc.db),
	}
}

func (s *CoinHistoryService) GetUpdatePriceHeartbeat(symbol string) (vo.HEARTBEAT_INFO, error) {

	latestInfos := s.gPriceInfosCache.GetLatestPriceInfos()

	var symbolPriceInfo = make([]conf.PriceInfo, 0)
	for _, info := range latestInfos.PriceInfos {
		if strings.EqualFold(info.Symbol, symbol) {
			symbolPriceInfo = append(symbolPriceInfo, info)
		}
	}

	if len(symbolPriceInfo) == 0 {
		return vo.HEARTBEAT_INFO{}, errors.New("not found symbol")
	}
	tokenSymbol := strings.ReplaceAll(symbol, "usdt", "-usdt")
	exchangeConfs := s.gReqeustPriceConfs.GetConfsBySymbol(tokenSymbol)
	return vo.HEARTBEAT_INFO{
		ExpectResources: len(exchangeConfs),
		ActualResources: len(symbolPriceInfo),
		LatestTimestamp: symbolPriceInfo[0].TimeStamp,
		Interval:        conf.GCfg.Interval,
	}, nil
}

func (s *CoinHistoryService) GetUpdatePriceHistory(idx, pageSize int, symbol string) (int, []vo.UpdatePriceHistoryResp, error) {

	histories, err := s.updatePriceRepo.GetUpdatePriceHistoryBySymbol(idx, pageSize, symbol)
	if err != nil {
		logger.WithError(err).Errorf("get history by symbol occur error,symbol:%s,index:%d", symbol, idx)
		return 0, nil, err
	}
	total, err := s.updatePriceRepo.GetTotalUpdatePriceHistoryBySymbol(symbol)
	if err != nil {
		logger.WithError(err).Errorf("get total history by symbol occur error,symbol:%s", symbol)
		return 0, nil, err
	}

	historyResps := make([]vo.UpdatePriceHistoryResp, 0)
	for _, history := range histories {
		infos, err := s.coinHistoryRepo.GetHistoryBySymbolAndTimestamp(history.Symbol, history.Timestamp)
		if err != nil {
			logger.WithError(err).Errorf("get history by symbol and timestamp occur error,symbol:%s", symbol)
			return 0, nil, err
		}

		bFind, partyPriceData := util.PartyPrice(infos, symbol, true)

		if !bFind {
			logger.Infoln("partyPrice error, symbol:", symbol)

			return 0, nil, err
		}

		historyResp := vo.UpdatePriceHistoryResp{
			Timestamp: history.Timestamp,
			Symbol:    history.Symbol,
			Price:     partyPriceData.Price,
			Infos:     infos,
		}
		historyResps = append(historyResps, historyResp)
	}
	return total, historyResps, nil
}
