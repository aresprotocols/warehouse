package service

import (
	logger "github.com/sirupsen/logrus"
	"price_api/price_server/internal/cache"
	"price_api/price_server/internal/repository"
)

type UpdateIntervalService struct {
	updateIntervalRepo  repository.UpdateIntervalRepository
	updateIntervalCache cache.GlobalUpdateIntervalCache
}

func newUpdateInterval(svc *service) *UpdateIntervalService {
	return &UpdateIntervalService{
		updateIntervalRepo:  repository.UpdateIntervalRepository{DB: svc.db},
		updateIntervalCache: svc.globalUpdateIntervalCache,
	}
}

func (s *UpdateIntervalService) SetUpdateInterval(symbol string, interval int) error {

	err := s.updateIntervalRepo.SetUpdateInterval(symbol, interval)
	if err != nil {
		logger.WithError(err).Errorf("set update interval occur error,symbol:%s,interval:%d", symbol, interval)
		return err
	}

	s.updateIntervalCache.UpdateSymbolInterval(symbol, interval)
	return nil
}

func (s *UpdateIntervalService) CheckUpdateInterval(symbol string, interval int) (int, error) {
	return s.updateIntervalRepo.CheckUpdateInterval(symbol, interval)
}

func (s *UpdateIntervalService) GetIntervalFromCache(symbol string) int {
	return s.updateIntervalCache.GetSymbolInterval(symbol)
}
