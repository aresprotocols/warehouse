package service

import (
	logger "github.com/sirupsen/logrus"
	conf "price_api/price_server/config"
	"price_api/price_server/internal/repository"
	"price_api/price_server/internal/vo"
)

type HttpErrorService struct {
	httpErrorRepo repository.HttpErrorRepository
}

func newHttpError(svc *service) *HttpErrorService {
	return &HttpErrorService{httpErrorRepo: repository.HttpErrorRepository{DB: svc.db}}
}

func (s *HttpErrorService) GetHttpErrorsByPage(idx int, symbol string) (int, []vo.HTTP_ERROR_INFO, error) {

	total, err := s.httpErrorRepo.GetTotalHttpErrorInfo(symbol)
	if err != nil {
		logger.WithError(err).Error("GetTotalHttpErrorInfo occur error")
		return 0, nil, err
	}

	httpErrorInfos, err := s.httpErrorRepo.GetHttpErrorInfo(idx, symbol, int(conf.GCfg.PageSize))
	if err != nil {
		logger.WithError(err).Error("GetHttpErrorInfo occur error")
		return 0, nil, err
	}
	return total, httpErrorInfos, nil
}

func (s *HttpErrorService) InsertHttpError(url string, symbol string, errorInfo string) error {
	return s.httpErrorRepo.InsertHttpError(url, symbol, errorInfo)
}
