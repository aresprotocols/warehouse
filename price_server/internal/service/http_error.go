package service

import (
	logger "github.com/sirupsen/logrus"
	"price_api/price_server/internal/repository"
	"price_api/price_server/internal/vo"
)

type HttpErrorService struct {
	httpErrorRepo repository.HttpErrorRepository
}

func newHttpError(svc *service) *HttpErrorService {
	return &HttpErrorService{httpErrorRepo: repository.NewHttpErrorRepository(svc.db)}
}

func (s *HttpErrorService) GetHttpErrorsByPage(idx, pageSize int, symbol string) (int, []vo.HTTP_ERROR_INFO, error) {

	total, err := s.httpErrorRepo.GetTotalHttpErrorInfo(symbol)
	if err != nil {
		logger.WithError(err).Error("GetTotalHttpErrorInfo occur error")
		return 0, nil, err
	}

	httpErrorInfos, err := s.httpErrorRepo.GetHttpErrorInfo(idx, pageSize, symbol)
	if err != nil {
		logger.WithError(err).Error("GetHttpErrorInfo occur error")
		return 0, nil, err
	}
	return total, httpErrorInfos, nil
}

func (s *HttpErrorService) InsertHttpError(url string, symbol string, errorInfo string) error {
	return s.httpErrorRepo.InsertHttpError(url, symbol, errorInfo)
}
