package repository

import (
	"github.com/jmoiron/sqlx"
	"strings"
)

type WeightInfoRepository struct {
	DB *sqlx.DB
}

//check symbo exchangeName in db, if not, update. if in. get weight return
func (r *WeightInfoRepository) CheckUpdateWeight(symbol, exchangeName string, weight int64) (int64, error) {
	var weightDb int64
	querySql := "select weight from " + TABLE_WEIGH_INFO + " where symbol = ? and exchange = ?"
	err := r.DB.Get(&weightDb, querySql, symbol, exchangeName)
	if err != nil {
		//no result, insert weight to db
		if strings.Contains(err.Error(), "no rows in result set") {
			insertSql := "insert into " + TABLE_WEIGH_INFO + " (symbol,exchange,weight)" +
				" values(?,?,?)"
			_, err := r.DB.Exec(insertSql, symbol, exchangeName, weight)
			if err != nil {
				return weight, err
			} else {
				return weight, nil
			}
		} else {
			return weight, err
		}
	}

	return weightDb, nil
}

func (r *WeightInfoRepository) SetWeight(symbol, exchangeName string, weight int) error {
	var weightDb int64
	querySql := "select weight from " + TABLE_WEIGH_INFO + " where symbol = ? and exchange = ?"
	err := r.DB.Get(&weightDb, querySql, symbol, exchangeName)
	if err != nil {
		return err
	}

	updateSql := "update " + TABLE_WEIGH_INFO + " set weight = ? where symbol = ? and exchange = ?"
	_, err = r.DB.Exec(updateSql, weight, symbol, exchangeName)
	if err != nil {
		return err
	}

	return nil
}
