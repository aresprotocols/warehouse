package service

import (
	"github.com/jmoiron/sqlx"
	gocache "github.com/patrickmn/go-cache"
	"price_api/price_server/internal/cache"
)

// Svc global var
var Svc Service

type Service interface {
	Auth() *AuthService
	Price() *PriceService
	HttpError() *HttpErrorService
	RequestInfo() *RequestInfoService
	CoinHistory() *CoinHistoryService
	WeightInfo() *WeightInfoService
	RequestPriceConf() *RequestPriceConfService
	UpdateInterval() *UpdateIntervalService
	Gas() *GasService
}

type service struct {
	db                        *sqlx.DB
	globalCache               cache.GlobalPriceInfoCache
	globalRequestPriceConfs   cache.GlobalRequestPriceConfs
	globalUpdateIntervalCache cache.GlobalUpdateIntervalCache
	goCache                   *gocache.Cache
}

// New init service
func New(db *sqlx.DB, goCache *gocache.Cache) Service {
	globalCahe := cache.NewGlobalPriceInfoCache()
	globalRequestPriceConfs := cache.NewGlobalRequestPriceConfs()
	globalUpdateIntervalCache := cache.NewGlobalUpdateIntervalCache()
	return &service{
		db:                        db,
		globalCache:               globalCahe,
		globalRequestPriceConfs:   globalRequestPriceConfs,
		globalUpdateIntervalCache: globalUpdateIntervalCache,
		goCache:                   goCache,
	}
}

func (s *service) Auth() *AuthService {
	return newAuth()
}

func (s *service) Price() *PriceService {
	return newPrice(s)
}

func (s *service) HttpError() *HttpErrorService {
	return newHttpError(s)
}

func (s *service) RequestInfo() *RequestInfoService {
	return newRequestInfo(s)
}

func (s *service) CoinHistory() *CoinHistoryService {
	return newCoinHistory(s)
}

func (s *service) WeightInfo() *WeightInfoService {
	return newWeightInfo(s)
}

func (s *service) RequestPriceConf() *RequestPriceConfService {
	return newRequestPriceConfService(s)
}

func (s *service) UpdateInterval() *UpdateIntervalService {
	return newUpdateInterval(s)
}

func (s *service) Gas() *GasService {
	return newGas(s)
}
