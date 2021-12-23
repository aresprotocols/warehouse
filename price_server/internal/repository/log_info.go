package repository

import (
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	"price_api/price_server/internal/vo"
	"strconv"
)

type LogInfoRepository struct {
	DB *sqlx.DB
}

func (r *LogInfoRepository) InsertLogInfo(mapInfo map[string]interface{}, t int) error {
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

func (r *LogInfoRepository) GetLogInfo(idx int, pageSize int) (vo.LOG_INFOS, error) {
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

func (r *LogInfoRepository) GetTotalLogInfoBySymbol(symbol string) (int, error) {
	var total int
	querySql := "select count(1) from " +
		TABLE_LOG_INFO + " where ( request_response like '%" + symbol + "%'" +
		" or request_url like '%" + symbol + "%'" + " ) and use_symbol = 1 ;"
	logger.Infoln("sql:", querySql)
	err := r.DB.QueryRow(querySql).Scan(&total)
	if err != nil {
		return total, err
	}

	return total, nil
}

func (r *LogInfoRepository) GetLogInfoBySymbol(idx int, pageSize int, symbol string) ([]vo.REQ_RSP_LOG_INFO, error) {
	var logInfos []vo.REQ_RSP_LOG_INFO
	querySql := "select client_ip,request_url,request_time,request_response,request_timestamp from " +
		TABLE_LOG_INFO + " where ( request_response like '%" + symbol + "%'" +
		" or request_url like '%" + symbol + "%'" + " ) and use_symbol = 1 order by id desc limit ?,?;"
	logger.Infoln("sql:", querySql, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	err := r.DB.Select(&logInfos, querySql, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return logInfos, err
	}

	return logInfos, nil
}
