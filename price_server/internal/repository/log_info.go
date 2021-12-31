package repository

import (
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	"price_api/price_server/internal/vo"
	"strconv"
)

//go:generate mockgen -destination mock/log_info_mock.go price_api/price_server/internal/repository LogInfoRepository

type LogInfoRepository interface {
	InsertLogInfo(mapInfo map[string]interface{}, t int) error
	GetLogInfo(idx int, pageSize int) (vo.LOG_INFOS, error)
	GetTotalLogInfoBySymbol(symbol string, ip string) (int, error)
	GetLogInfoBySymbol(idx int, pageSize int, symbol string, ip string) ([]vo.REQ_RSP_LOG_INFO, error)
}

func NewLogInfoRepository(db *sqlx.DB) LogInfoRepository {
	return &logInfoRepository{db}
}

type logInfoRepository struct {
	DB *sqlx.DB
}

func (r *logInfoRepository) InsertLogInfo(mapInfo map[string]interface{}, t int) error {
	insertSql := "insert into " + TABLE_LOG_INFO + " (client_ip,request_time,user_agent,request_url," +
		"response_time,request_response, use_symbol,request_timestamp,response_timestamp)" +
		" values(?,?,?,?," +
		"?,?,?,?,?)"
	_, err := r.DB.Exec(insertSql, mapInfo["request_client_ip"], mapInfo["request_time"], mapInfo["request_ua"], mapInfo["request_uri"],
		mapInfo["response_time"], mapInfo["response"], t, mapInfo["request_timestamp"], mapInfo["response_timestamp"])
	if err != nil {
		return err
	}

	return nil
}

func (r *logInfoRepository) GetLogInfo(idx int, pageSize int) (vo.LOG_INFOS, error) {
	var logInfos vo.LOG_INFOS
	querySql := "select client_ip," +
		"request_time,user_agent,request_url,response_time,request_response from " +
		TABLE_LOG_INFO + " order by id desc limit ?,?;"
	logger.Infoln("sql:", querySql, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	err := r.DB.Select(&logInfos.Infos, querySql, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return vo.LOG_INFOS{}, err
	}

	return logInfos, nil
}

func (r *logInfoRepository) GetTotalLogInfoBySymbol(symbol string, ip string) (int, error) {
	var total int
	argsArr := make([]interface{}, 0)
	querySql := "select count(1) from " +
		TABLE_LOG_INFO + " where ( request_response like '%" + symbol + "%'" +
		" or request_url like '%" + symbol + "%'" + " ) and use_symbol = 1 "
	if ip != "" {
		querySql += " and client_ip = ? "
		argsArr = append(argsArr, ip)
	}
	logger.Infoln("sql:", querySql, "args", argsArr)
	err := r.DB.QueryRow(querySql, argsArr...).Scan(&total)
	if err != nil {
		return total, err
	}

	return total, nil
}

func (r *logInfoRepository) GetLogInfoBySymbol(idx int, pageSize int, symbol string, ip string) ([]vo.REQ_RSP_LOG_INFO, error) {
	var logInfos []vo.REQ_RSP_LOG_INFO
	argsArr := make([]interface{}, 0)

	querySql := "select client_ip,request_url,request_time,request_response,request_timestamp from " +
		TABLE_LOG_INFO + " where ( request_response like '%" + symbol + "%'" +
		" or request_url like '%" + symbol + "%'" + " ) and use_symbol = 1 "

	if ip != "" {
		querySql += " and client_ip = ? "
		argsArr = append(argsArr, ip)
	}
	querySql += "order by id desc limit ?,?;"
	argsArr = append(argsArr, strconv.Itoa(idx*pageSize))
	argsArr = append(argsArr, strconv.Itoa(pageSize))
	logger.Infoln("sql:", querySql, " args:", argsArr)
	err := r.DB.Select(&logInfos, querySql, argsArr...)
	if err != nil {
		return logInfos, err
	}

	return logInfos, nil
}
