package service

import (
	conf "price_api/price_server/config"
	"price_api/price_server/internal/cache"
)

type RequestPriceConfService struct {
	gRequestPriceConfs cache.GlobalRequestPriceConfs
}

func newRequestPriceConfService(svc *service) *RequestPriceConfService {
	return &RequestPriceConfService{
		gRequestPriceConfs: svc.globalRequestPriceConfs,
	}
}
func (s *RequestPriceConfService) SetConfs(conf map[string][]conf.ExchangeConfig) {
	s.gRequestPriceConfs.SetConfs(conf)

}

func (s *RequestPriceConfService) GetConfs() map[string][]conf.ExchangeConfig {
	return s.gRequestPriceConfs.GetConfs()
}
