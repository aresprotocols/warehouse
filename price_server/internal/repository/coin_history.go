package repository

import (
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	conf "price_api/price_server/config"
	"strconv"
)

type CoinHistoryRepository struct {
	DB *sqlx.DB
}

func (r *CoinHistoryRepository) InsertPriceInfo(cfg conf.PriceInfos) error {
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
			_, err := r.DB.Exec(insertUpdateHistorySql, kTimestamp, kSymbol)
			if err != nil {
				return err
			}
		}
	}

	for _, info := range cfg.PriceInfos {
		//TODO battle
		_, err := r.DB.Exec(insertSql, info.Symbol, info.TimeStamp, info.Price, info.PriceOrigin, info.Weight)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *CoinHistoryRepository) GetTotalHistoryBySymbol(symbol string) (int, error) {
	var total int
	querySql := "select count(1) from `" + TABLE_COIN_PRICE + "` where symbol = ?;"
	logger.Infoln("sql:", querySql, "symbol", symbol)
	err := r.DB.QueryRow(querySql, symbol).Scan(&total)
	if err != nil {
		return total, err
	}
	return total, nil
}

func (r *CoinHistoryRepository) GetHistoryBySymbol(idx int, pageSize int, symbol string) ([]conf.PriceInfo, error) {
	var infos []conf.PriceInfo
	querySql := "select symbol, timestamp, price, weight, price_origin from `" + TABLE_COIN_PRICE + "` where symbol = ? order by id desc limit ?,? ;"
	logger.Infoln("sql:", querySql, "symbol", symbol, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))

	err := r.DB.Select(&infos, querySql, symbol, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return infos, err
	}
	return infos, nil
}

func (r *CoinHistoryRepository) GetHistoryBySymbolAndTimestamp(symbol string, timestamp int64) ([]conf.PriceInfo, error) {
	var infos []conf.PriceInfo
	querySql := "select symbol, timestamp, price, weight, price_origin from `" + TABLE_COIN_PRICE + "` where symbol = ? and timestamp = ? order by id desc ;"
	logger.Infoln("sql:", querySql, "symbol", symbol, "timestamp", timestamp)

	err := r.DB.Select(&infos, querySql, symbol, timestamp)
	if err != nil {
		return infos, err
	}
	return infos, nil
}

func (r *CoinHistoryRepository) GetHistoryByTimestamp(timestamp int64) ([]conf.PriceInfo, error) {

	dbPriceInfos := make([]conf.PriceInfo, 0)
	querySql := "select symbol, timestamp, price, weight, price_origin from `t_coin_history_info` where timestamp = ?;"
	err := r.DB.Select(&dbPriceInfos, querySql, timestamp)
	if err != nil {
		return []conf.PriceInfo{}, err
	}

	return dbPriceInfos, nil
}
