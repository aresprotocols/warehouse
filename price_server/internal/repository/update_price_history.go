package repository

import (
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	"price_api/price_server/internal/vo"
	"strconv"
)

//go:generate mockgen -destination mock/update_price_history_mock.go price_api/price_server/internal/repository UpdatePriceRepository

type UpdatePriceRepository interface {
	GetTotalUpdatePriceHistoryBySymbol(symbol string) (int, error)
	GetUpdatePriceHistoryBySymbol(idx int, pageSize int, symbol string) ([]vo.UpdatePirceHistory, error)
	DeleteOldLogs(timestamp int64) error
	GetUpdatePriceHistoryByInterval(beforeTimestamp int64, symbol string) ([]vo.UpdatePirceHistory, error)
}

func NewUpdatePriceRepository(db *sqlx.DB) UpdatePriceRepository {
	return &updatePriceRepository{db}
}

type updatePriceRepository struct {
	DB *sqlx.DB
}

func (r *updatePriceRepository) GetTotalUpdatePriceHistoryBySymbol(symbol string) (int, error) {
	var total int
	querySql := "select count(1) from `" + TABLE_UPDATE_PRICE_HISTORY + "` where symbol = ?;"
	logger.Infoln("sql:", querySql, "symbol", symbol)
	err := r.DB.QueryRow(querySql, symbol).Scan(&total)
	if err != nil {
		return total, err
	}
	return total, nil
}

func (r *updatePriceRepository) GetUpdatePriceHistoryBySymbol(idx int, pageSize int, symbol string) ([]vo.UpdatePirceHistory, error) {
	var histories []vo.UpdatePirceHistory
	querySql := "select symbol, timestamp  from `" + TABLE_UPDATE_PRICE_HISTORY + "` where symbol = ? order by timestamp desc limit ?,? ;"
	logger.Infoln("sql:", querySql, "symbol", symbol, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))

	err := r.DB.Select(&histories, querySql, symbol, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return histories, err
	}
	return histories, nil
}

func (r *updatePriceRepository) GetUpdatePriceHistoryByInterval(beforeTimestamp int64, symbol string) ([]vo.UpdatePirceHistory, error) {
	var histories []vo.UpdatePirceHistory
	querySql := "select b.timestamp,b.symbol from " +
		"(SELECT min(timestamp) as timestamp,symbol,floor(timestamp / (10 * 60)) * (10 * 60) as intervals from `" +
		TABLE_UPDATE_PRICE_HISTORY +
		"` where timestamp >= ? and symbol = ? group by intervals,symbol) as b order by b.timestamp"
	logger.Infoln("sql:", querySql, "symbol", symbol)
	err := r.DB.Select(&histories, querySql, beforeTimestamp, symbol)
	if err != nil {
		return histories, err
	}
	return histories, nil
}

func (r *updatePriceRepository) DeleteOldLogs(timestamp int64) error {
	argsArr := make([]interface{}, 0)
	querySql := "delete from " + TABLE_UPDATE_PRICE_HISTORY + " where timestamp <= ? "
	argsArr = append(argsArr, timestamp)
	logger.Infoln("sql:", querySql, " args:", argsArr)
	result, err := r.DB.Exec(querySql, argsArr...)
	if err != nil {
		logger.WithError(err).Errorf("exec delete update price history occur err")
		return err
	}
	effectRows, _ := result.RowsAffected()
	logger.Infof("exec delete update price history success,%d rows", effectRows)
	return nil
}
