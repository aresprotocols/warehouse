package repository

import (
	"github.com/jmoiron/sqlx"
	"strings"
)

type UpdateIntervalRepository struct {
	DB *sqlx.DB
}

func (r *UpdateIntervalRepository) CheckUpdateInterval(symbol string, interval int) (int, error) {
	var intervalDb int
	querySql := "select interval_second from " + TABLE_UPDATE_INTERVAL + " where symbol = ?"
	err := r.DB.Get(&intervalDb, querySql, symbol)
	if err != nil {
		//no result, insert interval to db
		if strings.Contains(err.Error(), "no rows in result set") {
			insertSql := "insert into " + TABLE_UPDATE_INTERVAL + " (symbol,interval_second)" +
				" values(?,?)"
			_, err := r.DB.Exec(insertSql, symbol, interval)
			if err != nil {
				return interval, err
			} else {
				return interval, nil
			}
		} else {
			return interval, err
		}
	}

	return intervalDb, nil
}

func (r *UpdateIntervalRepository) SetUpdateInterval(symbol string, interval int) error {
	var intervalDb int64
	querySql := "select interval_second from " + TABLE_UPDATE_INTERVAL + " where symbol = ?"
	err := r.DB.Get(&intervalDb, querySql, symbol)
	if err != nil {
		return err
	}

	updateSql := "update " + TABLE_UPDATE_INTERVAL + " set interval_second = ? where symbol = ?"
	_, err = r.DB.Exec(updateSql, interval, symbol)
	if err != nil {
		return err
	}

	return nil
}
