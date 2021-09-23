package sql

import (
	"fmt"
	"log"
	conf "price_api/price_server/config"
	"strconv"

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

	err = mysqlDb.Ping()
	if err != nil {
		return err
	}

	db = mysqlDb
	db.SetMaxOpenConns(1)
	createOrderTables(mysqlDb, cfg.Mysql.Db)

	return nil
}

const TABLE_COIN_PRICE = "t_coin_history_info"
const TABLE_LOG_INFO = "t_log_info"

func InsertPriceInfo(cfg conf.PriceInfos) error {
	insertSql := "insert into " + TABLE_COIN_PRICE + " (symbol,timestamp,price,price_origin,weight)" +
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

func InsertLogInfo(mapInfo map[string]string) error {
	insertSql := "insert into " + TABLE_LOG_INFO + " (client_ip,method,post_data,request_proto," +
		"request_time,user_agent,request_url,response_time,request_response)" +
		" values(?,?,?,?," +
		"?,?,?,?,?)"
	_, err := db.Exec(insertSql, mapInfo["request_client_ip"], mapInfo["request_method"], mapInfo["request_post_data"], mapInfo["request_proto"],
		mapInfo["request_time"], mapInfo["request_ua"], mapInfo["request_uri"], mapInfo["response_time"], mapInfo["response"])
	if err != nil {
		return err
	}

	return nil
}

type LOG_INFO struct {
	ClientIP     string `db:"client_ip" json:"client_ip"`
	Method       string `db:"method" json:"method"`
	PostData     string `db:"post_data" json:"post_data"`
	Proto        string `db:"request_proto" json:"proto"`
	RequestTime  string `db:"request_time" json:"request_time"`
	UserAgent    string `db:"user_agent" json:"user_agent"`
	RequestUrl   string `db:"request_url" json:"request_url"`
	ResponseTime string `db:"response_time" json:"response_time"`
	Response     string `db:"request_response" json:"response"`
}

type LOG_INFOS struct {
	Infos []LOG_INFO `json:"infos"`
}

func GetLogInfo(idx int, pageSize int) (LOG_INFOS, error) {
	var logInfos LOG_INFOS
	querySql := "select client_ip,method,post_data,request_proto," +
		"request_time,user_agent,request_url,response_time,request_response from " +
		TABLE_LOG_INFO + " order by id desc limit ?,?;"
	log.Println("sql:", querySql, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	err := db.Select(&logInfos.Infos, querySql, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return LOG_INFOS{}, err
	}

	return logInfos, nil
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
