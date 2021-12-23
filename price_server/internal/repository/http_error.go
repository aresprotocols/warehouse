package repository

import (
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	"price_api/price_server/internal/vo"
	"strconv"
	"time"
)

type HttpErrorRepository struct {
	DB *sqlx.DB
}

func (r *HttpErrorRepository) InsertHttpError(url string, symbol string, errorInfo string) error {
	insertSql := "insert into " + TABLE_HTTP_ERROR + " (url,symbol,error,timestamp)" +
		" values(?,?,?,?)"

	_, err := r.DB.Exec(insertSql, url, symbol, errorInfo, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

func (r *HttpErrorRepository) GetHttpErrorInfo(idx int, symbol string, pageSize int) ([]vo.HTTP_ERROR_INFO, error) {
	var infos = make([]vo.HTTP_ERROR_INFO, 0)
	querySql := "select url,symbol,error,timestamp from " +
		TABLE_HTTP_ERROR + " where symbol = ? order by id desc limit ?,?;"
	logger.Infoln("sql:", querySql, "symbol: ", symbol, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	err := r.DB.Select(&infos, querySql, symbol, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return nil, err
	}

	return infos, nil
}
func (r *HttpErrorRepository) GetTotalHttpErrorInfo(symbol string) (int, error) {
	var total int
	querySql := "select count(1) from " + TABLE_HTTP_ERROR + " where symbol = ?;"
	logger.Infoln("sql:", querySql, "symbol: ", symbol)
	err := r.DB.QueryRow(querySql, symbol).Scan(&total)
	if err != nil {
		return total, err
	}
	return total, nil
}
