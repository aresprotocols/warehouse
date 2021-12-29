package service

import (
	logger "github.com/sirupsen/logrus"
	"price_api/price_server/internal/cache"
	"price_api/price_server/internal/repository"
	"price_api/price_server/internal/vo"
)

type WeightInfoService struct {
	weightInfoRepo   repository.WeightInfoRepository
	gPriceInfosCache cache.GlobalPriceInfoCache
}

func newWeightInfo(svc *service) *WeightInfoService {
	return &WeightInfoService{
		weightInfoRepo:   repository.WeightInfoRepository{DB: svc.db},
		gPriceInfosCache: svc.globalCache,
	}
}

func (s *WeightInfoService) SetWeight(setWeightReq vo.SetWeightReq) error {

	err := s.weightInfoRepo.SetWeight(setWeightReq.Symbol, setWeightReq.Exchange, setWeightReq.Weight)
	if err != nil {
		logger.WithError(err).Errorf("set weight occur error,symbol:%s,exchange:%s,weight:%d", setWeightReq.Symbol, setWeightReq.Exchange, setWeightReq.Weight)
		return err
	}

	s.gPriceInfosCache.UpdateSymbolWeight(setWeightReq.Symbol, setWeightReq.Exchange, setWeightReq.Weight)
	return nil
}

func (s *WeightInfoService) CheckUpdateWeight(symbol, exchangeName string, weight int64) (int64, error) {
	return s.weightInfoRepo.CheckUpdateWeight(symbol, exchangeName, weight)
}
