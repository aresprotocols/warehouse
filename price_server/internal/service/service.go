package service

import (
	"github.com/jmoiron/sqlx"
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
}

type service struct {
	db          *sqlx.DB
	globalCache *cache.GlobalPriceInfoCache
}

// New init service
func New(db *sqlx.DB) Service {
	globalCahe := cache.NewGlobalPriceInfoCache()
	return &service{
		db:          db,
		globalCache: globalCahe,
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
