package repository

import (
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	"price_api/price_server/internal/vo"
	"strconv"
	"time"
)

//go:generate mockgen -destination mock/http_error_mock.go price_api/price_server/internal/repository HttpErrorRepository

type HttpErrorRepository interface {
	InsertHttpError(url string, symbol string, errorInfo string) error
	GetHttpErrorInfo(idx, pageSize int, symbol string) ([]vo.HTTP_ERROR_INFO, error)
	GetTotalHttpErrorInfo(symbol string) (int, error)
}

func NewHttpErrorRepository(db *sqlx.DB) HttpErrorRepository {
	return &httpErrorRepository{db}
}

type httpErrorRepository struct {
	db *sqlx.DB
}

func (r *httpErrorRepository) InsertHttpError(url string, symbol string, errorInfo string) error {
	insertSql := "insert into " + TABLE_HTTP_ERROR + " (url,symbol,error,timestamp)" +
		" values(?,?,?,?)"

	_, err := r.db.Exec(insertSql, url, symbol, errorInfo, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func (r *httpErrorRepository) GetHttpErrorInfo(idx, pageSize int, symbol string) ([]vo.HTTP_ERROR_INFO, error) {
	var infos = make([]vo.HTTP_ERROR_INFO, 0)
	querySql := "select url,symbol,error,timestamp from " +
		TABLE_HTTP_ERROR + " where symbol = ? order by id desc limit ?,?;"
	logger.Infoln("sql:", querySql, "symbol: ", symbol, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	err := r.db.Select(&infos, querySql, symbol, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return nil, err
	}

	return infos, nil
}
func (r *httpErrorRepository) GetTotalHttpErrorInfo(symbol string) (int, error) {
	var total int
	querySql := "select count(1) from " + TABLE_HTTP_ERROR + " where symbol = ?;"
	logger.Infoln("sql:", querySql, "symbol: ", symbol)
	err := r.db.QueryRow(querySql, symbol).Scan(&total)
	if err != nil {
		return total, err
	}
	return total, nil
}
