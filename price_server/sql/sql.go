package sql

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	conf "price_api/price_server/config"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

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

	db = mysqlDb

	return nil
}

const TABLE_COIN_PRICE = "t_coin_history_info"
const TABLE_LOG_INFO = "t_log_info"
const TABLE_HTTP_ERROR = "t_http_error"
const TABLE_WEIGH_INFO = "t_weight_info"
const TABLE_UPDATE_PRICE_HISTORY = "t_update_price_history"

func InsertPriceInfo(cfg conf.PriceInfos) error {
	insertSql := "insert into " + TABLE_COIN_PRICE + " (symbol,timestamp,price,price_origin,weight)" +
		" values(?,?,?,?,?)"

	insertUpdateHistorySql := "insert into " + TABLE_UPDATE_PRICE_HISTORY + "(timestamp,symbol) value (?,?)"

	historyMap := make(map[int64]map[string]struct{})

	for _, info := range cfg.PriceInfos {
		if _, timestampOk := historyMap[info.TimeStamp]; timestampOk {
			symbolMap := historyMap[info.TimeStamp]
			if _, symbolOk := symbolMap[info.Symbol]; !symbolOk {
				symbolMap[info.Symbol] = struct{}{}
			}
		} else {
			symbolMap := make(map[string]struct{})
			symbolMap[info.Symbol] = struct{}{}
			historyMap[info.TimeStamp] = symbolMap
		}
	}

	for kTimestamp, _ := range historyMap {
		for kSymbol, _ := range historyMap[kTimestamp] {
			_, err := db.Exec(insertUpdateHistorySql, kTimestamp, kSymbol)
			if err != nil {
				return err
			}
		}
	}

	for _, info := range cfg.PriceInfos {
		//TODO battle
		_, err := db.Exec(insertSql, info.Symbol, info.TimeStamp, info.Price, info.PriceOrigin, info.Weight)
		if err != nil {
			return err
		}
	}

	return nil
}

func InsertLogInfo(mapInfo map[string]interface{}, t int) error {
	insertSql := "insert into " + TABLE_LOG_INFO + " (client_ip,request_time,user_agent,request_url," +
		"response_time,request_response, use_symbol,request_timestamp,response_timestamp)" +
		" values(?,?,?,?," +
		"?,?,?,?,?)"
	_, err := db.Exec(insertSql, mapInfo["request_client_ip"], mapInfo["request_time"], mapInfo["request_ua"], mapInfo["request_uri"],
		mapInfo["response_time"], mapInfo["response"], t, mapInfo["request_timestamp"], mapInfo["response_timestamp"])
	if err != nil {
		return err
	}

	return nil
}

func InsertHttpError(url string, symbol string, errorInfo string) error {
	insertSql := "insert into " + TABLE_HTTP_ERROR + " (url,symbol,error,timestamp)" +
		" values(?,?,?,?)"

	_, err := db.Exec(insertSql, url, symbol, errorInfo, time.Now().Unix())
	if err != nil {
		return err
	}

	return nil
}

type LOG_INFO struct {
	ClientIP     string `db:"client_ip" json:"client_ip"`
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
	querySql := "select client_ip," +
		"request_time,user_agent,request_url,response_time,request_response from " +
		TABLE_LOG_INFO + " order by id desc limit ?,?;"
	logger.Infoln("sql:", querySql, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	err := db.Select(&logInfos.Infos, querySql, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return LOG_INFOS{}, err
	}

	return logInfos, nil
}

type REQ_RSP_LOG_INFO struct {
	ReqUrl           string `json:"reqUrl" db:"request_url"`
	Response         string `json:"response" db:"request_response"`
	Ip               string `json:"ip" db:"client_ip"`
	RequestTime      string `json:"request_time" db:"request_time"`
	RequestTimestamp int64  `json:"request_timestamp" db:"request_timestamp"`
}

func GetTotalLogInfoBySymbol(symbol string) (int, error) {
	var total int
	querySql := "select count(1) from " +
		TABLE_LOG_INFO + " where ( request_response like '%" + symbol + "%'" +
		" or request_url like '%" + symbol + "%'" + " ) and use_symbol = 1 ;"
	logger.Infoln("sql:", querySql)
	err := db.QueryRow(querySql).Scan(&total)
	if err != nil {
		return total, err
	}

	return total, nil
}

func GetLogInfoBySymbol(idx int, pageSize int, symbol string) ([]REQ_RSP_LOG_INFO, error) {
	var logInfos []REQ_RSP_LOG_INFO
	querySql := "select client_ip,request_url,request_time,request_response,request_timestamp from " +
		TABLE_LOG_INFO + " where ( request_response like '%" + symbol + "%'" +
		" or request_url like '%" + symbol + "%'" + " ) and use_symbol = 1 order by id desc limit ?,?;"
	logger.Infoln("sql:", querySql, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	err := db.Select(&logInfos, querySql, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return logInfos, err
	}

	return logInfos, nil
}

func GetTotalHistoryBySymbol(symbol string) (int, error) {
	var total int
	querySql := "select count(1) from `" + TABLE_COIN_PRICE + "` where symbol = ?;"
	logger.Infoln("sql:", querySql, "symbol", symbol)
	err := db.QueryRow(querySql, symbol).Scan(&total)
	if err != nil {
		return total, err
	}
	return total, nil
}

func GetHistoryBySymbol(idx int, pageSize int, symbol string) ([]conf.PriceInfo, error) {
	var infos []conf.PriceInfo
	querySql := "select symbol, timestamp, price, weight, price_origin from `" + TABLE_COIN_PRICE + "` where symbol = ? order by id desc limit ?,? ;"
	logger.Infoln("sql:", querySql, "symbol", symbol, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))

	err := db.Select(&infos, querySql, symbol, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return infos, err
	}
	return infos, nil
}

type UpdatePirceHistory struct {
	Timestamp int64  `json:"timestamp" db:"timestamp"`
	Symbol    string `json:"symbol" db:"symbol"`
}

func GetTotalUpdatePriceHistoryBySymbol(symbol string) (int, error) {
	var total int
	querySql := "select count(1) from `" + TABLE_UPDATE_PRICE_HISTORY + "` where symbol = ?;"
	logger.Infoln("sql:", querySql, "symbol", symbol)
	err := db.QueryRow(querySql, symbol).Scan(&total)
	if err != nil {
		return total, err
	}
	return total, nil
}

func GetUpdatePriceHistoryBySymbol(idx int, pageSize int, symbol string) ([]UpdatePirceHistory, error) {
	var histories []UpdatePirceHistory
	querySql := "select symbol, timestamp  from `" + TABLE_UPDATE_PRICE_HISTORY + "` where symbol = ? order by timestamp desc limit ?,? ;"
	logger.Infoln("sql:", querySql, "symbol", symbol, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))

	err := db.Select(&histories, querySql, symbol, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return histories, err
	}
	return histories, nil
}

func GetHistoryBySymbolAndTimestamp(symbol string, timestamp int64) ([]conf.PriceInfo, error) {
	var infos []conf.PriceInfo
	querySql := "select symbol, timestamp, price, weight, price_origin from `" + TABLE_COIN_PRICE + "` where symbol = ? and timestamp = ? order by id desc ;"
	logger.Infoln("sql:", querySql, "symbol", symbol, "timestamp", timestamp)

	err := db.Select(&infos, querySql, symbol, timestamp)
	if err != nil {
		return infos, err
	}
	return infos, nil
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

type HTTP_ERROR_INFO struct {
	Url       string `db:"url" json:"url"`
	Symbol    string `db:"symbol" json:"symbol"`
	Error     string `db:"error" json:"error"`
	Timestamp int64  `db:"timestamp" jsoon:"timestamp"`
}

func GetHttpErrorInfo(idx int, symbol string, pageSize int) ([]HTTP_ERROR_INFO, error) {
	var infos = make([]HTTP_ERROR_INFO, 0)
	querySql := "select url,symbol,error,timestamp from " +
		TABLE_HTTP_ERROR + " where symbol = ? order by id desc limit ?,?;"
	logger.Infoln("sql:", querySql, "symbol: ", symbol, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	err := db.Select(&infos, querySql, symbol, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return nil, err
	}

	return infos, nil
}
func GetTotalHttpErrorInfo(symbol string) (int, error) {
	var total int
	querySql := "select count(1) from " + TABLE_HTTP_ERROR + " where symbol = ?;"
	logger.Infoln("sql:", querySql, "symbol: ", symbol)
	err := db.QueryRow(querySql, symbol).Scan(&total)
	if err != nil {
		return total, err
	}
	return total, nil
}

//check symbo exchangeName in db, if not, update. if in. get weight return
func CheckUpdateWeight(symbol, exchangeName string, weight int64) (int64, error) {
	var weightDb int64
	querySql := "select weight from " + TABLE_WEIGH_INFO + " where symbol = ? and exchange = ?"
	err := db.Get(&weightDb, querySql, symbol, exchangeName)
	if err != nil {
		//no result, insert weight to db
		if strings.Contains(err.Error(), "no rows in result set") {
			insertSql := "insert into " + TABLE_WEIGH_INFO + " (symbol,exchange,weight)" +
				" values(?,?,?)"
			_, err := db.Exec(insertSql, symbol, exchangeName, weight)
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

func SetWeight(symbol, exchangeName string, weight int) error {
	var weightDb int64
	querySql := "select weight from " + TABLE_WEIGH_INFO + " where symbol = ? and exchange = ?"
	err := db.Get(&weightDb, querySql, symbol, exchangeName)
	if err != nil {
		return err
	}

	updateSql := "update " + TABLE_WEIGH_INFO + " set weight = ? where symbol = ? and exchange = ?"
	_, err = db.Exec(updateSql, weight, symbol, exchangeName)
	if err != nil {
		return err
	}

	return nil
}
