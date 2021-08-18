package sql

import (
	"fmt"
	conf "price_api/price_server/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

//operation about mysql
func InitMysqlDB(cfg conf.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8", cfg.Mysql.Name, cfg.Mysql.Password, cfg.Mysql.Server, cfg.Mysql.Port)
	mysqlDb, err := sqlx.Open("mysql", dsn)
	//db.SetMaxIdleConns(conf.Mysql.Conn.MaxIdle)
	//db.SetMaxOpenConns(conf.Mysql.Conn.Maxopen)
	//db.SetConnMaxLifetime(5 * time.Minute)
	if err != nil {
		return err
	}

	db = mysqlDb
	createOrderTables(mysqlDb, cfg.Mysql.Db)

	return nil
}

const TABLE_COIN_PRICE = "t_coin_history_info"

func InsertPriceInfo(cfg conf.PriceInfos) error {
	insertSql := "insert into " + TABLE_COIN_PRICE + " (symbols,timestamp,price,price_origin,weight)" +
		" values(?,?,?,?,?)"
	for _, info := range cfg.PriceInfos {
		//TODO battle
		_, err := db.Exec(insertSql, info.Symbol, info.TimeStamp, info.Price, info.PriceOrigin, info.Weight)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetHistoryBySymbolTimestamp(symbol string, timestamp int64) ([]conf.PriceInfo, error) {
	var dbTimestamp int64
	querySql := "select timestamp from " + TABLE_COIN_PRICE + " where timestamp <= ? order by timestamp desc;"
	err := db.Get(&dbTimestamp, querySql, timestamp)
	if err != nil {
		return []conf.PriceInfo{}, err
	}

	var dbPriceInfos []conf.PriceInfo
	querySql = "select symbol, timestamp, price, weight, price_origin from `t_coin_history_info` where timestamp = ?;"
	err = db.Select(&dbPriceInfos, querySql, dbTimestamp)
	if err != nil {
		return []conf.PriceInfo{}, err
	}

	return dbPriceInfos, nil
}
