package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"price_api/price_server/internal/cache"
	"price_api/price_server/internal/repository"
	"price_api/price_server/internal/util"
	"price_api/price_server/internal/vo"
	"strings"
	"time"
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
		updatePriceRepo:    repository.NewUpdatePriceRepository(svc.db),
		coinHistoryRepo:    repository.NewCoinHistoryRepository(svc.db),
	}
}

func (s *CoinHistoryService) GetUpdatePriceHeartbeat(symbol string, interval int64) (vo.HEARTBEAT_INFO, error) {

	latestInfos := s.gPriceInfosCache.GetLatestPriceInfos(symbol)

	var symbolPriceInfo = latestInfos.PriceInfos
	if len(symbolPriceInfo) == 0 {
		return vo.HEARTBEAT_INFO{}, errors.New("not found symbol")
	}
	tokenSymbol := strings.ReplaceAll(symbol, "usdt", "-usdt")
	exchangeConfs := s.gReqeustPriceConfs.GetConfsBySymbol(tokenSymbol)
	return vo.HEARTBEAT_INFO{
		ExpectResources: len(exchangeConfs),
		ActualResources: len(symbolPriceInfo),
		LatestTimestamp: symbolPriceInfo[0].TimeStamp,
		Interval:        interval,
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

		bFind, partyPriceData := util.PartyPrice(infos, true)

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
func (s *CoinHistoryService) GetUpdatePriceHistoryForChart(symbol string) ([]vo.UpdatePriceHistoryResp, error) {

	now := time.Now()
	oldTime := now.AddDate(0, 0, -1)

	histories, err := s.updatePriceRepo.GetUpdatePriceHistoryByInterval(oldTime.Unix(), symbol)
	if err != nil {
		logger.WithError(err).Errorf("get history by interval occur error,symbol:%s", symbol)
		return nil, err
	}
	historyResps := make([]vo.UpdatePriceHistoryResp, 0)
	for _, history := range histories {
		infos, err := s.coinHistoryRepo.GetHistoryBySymbolAndTimestamp(history.Symbol, history.Timestamp)
		if err != nil {
			logger.WithError(err).Errorf("get history by interval occur error,symbol:%s", symbol)
			return nil, err
		}

		bFind, partyPriceData := util.PartyPrice(infos, true)

		if !bFind {
			logger.Infoln("partyPrice error, symbol:", symbol)

			return nil, err
		}

		historyResp := vo.UpdatePriceHistoryResp{
			Timestamp: history.Timestamp,
			Symbol:    history.Symbol,
			Price:     partyPriceData.Price,
			Infos:     infos,
		}
		historyResps = append(historyResps, historyResp)
	}
	return historyResps, nil
}

func (s *CoinHistoryService) DeleteOld() error {
	logger.Info("start delete old coin history")
	now := time.Now()
	oldTime := now.AddDate(0, -2, 0)
	oldTimestamp := oldTime.Unix()
	err := s.coinHistoryRepo.DeleteOldLogs(oldTimestamp)
	if err != nil {
		logger.WithError(err).Errorf("delete old coin history occur err")
		return err
	}
	logger.Info("start delete old update price history")
	err = s.updatePriceRepo.DeleteOldLogs(oldTimestamp)
	if err != nil {
		logger.WithError(err).Errorf("delete old update price occur err")
		return err
	}
	return nil
}
