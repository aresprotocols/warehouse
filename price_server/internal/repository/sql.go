package repository

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	conf "price_api/price_server/config"
)

var DB *sqlx.DB

func createTable(cfg conf.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8", cfg.Mysql.Name, cfg.Mysql.Password, cfg.Mysql.Server, cfg.Mysql.Port)
	mysqlDb, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer mysqlDb.Close()

	err = mysqlDb.Ping()
	if err != nil {
		return err
	}

	createOrderTables(mysqlDb, cfg.Mysql.Db)

	return nil
}

//operation about mysql
func InitMysqlDB(cfg conf.Config) error {
	err := createTable(cfg)
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", cfg.Mysql.Name, cfg.Mysql.Password, cfg.Mysql.Server, cfg.Mysql.Port, cfg.Mysql.Db)
	mysqlDb, err := sqlx.Open("mysql", dsn)
	//db.SetMaxIdleConns(conf.Mysql.Conn.MaxIdle)
	//db.SetMaxOpenConns(conf.Mysql.Conn.Maxopen)
	//db.SetConnMaxLifetime(5 * time.Minute)
	if err != nil {
		return err
	}

	DB = mysqlDb

	return nil
}

const TABLE_COIN_PRICE = "t_coin_history_info"
const TABLE_LOG_INFO = "t_log_info"
const TABLE_HTTP_ERROR = "t_http_error"
const TABLE_WEIGH_INFO = "t_weight_info"
const TABLE_UPDATE_PRICE_HISTORY = "t_update_price_history"
const TABLE_UPDATE_INTERVAL = "t_update_interval"
