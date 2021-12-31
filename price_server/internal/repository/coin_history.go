package repository

import (
	"github.com/jmoiron/sqlx"
	logger "github.com/sirupsen/logrus"
	conf "price_api/price_server/config"
	"strconv"
)

//go:generate mockgen -destination mock/coin_history_mock.go price_api/price_server/internal/repository CoinHistoryRepository

type CoinHistoryRepository interface {
	InsertPriceInfo(cfg conf.PriceInfos) error
	GetTotalHistoryBySymbol(symbol string) (int, error)
	GetHistoryBySymbol(idx int, pageSize int, symbol string) ([]conf.PriceInfo, error)
	GetHistoryBySymbolAndTimestamp(symbol string, timestamp int64) ([]conf.PriceInfo, error)
	GetHistoryByTimestamp(timestamp int64) ([]conf.PriceInfo, error)
}

func NewCoinHistoryRepository(db *sqlx.DB) CoinHistoryRepository {
	return &coinHistoryRepository{db}
}

type coinHistoryRepository struct {
	db *sqlx.DB
}

func (r *coinHistoryRepository) InsertPriceInfo(cfg conf.PriceInfos) error {
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

	for kTimestamp := range historyMap {
		for kSymbol := range historyMap[kTimestamp] {
			_, err := r.db.Exec(insertUpdateHistorySql, kTimestamp, kSymbol)
			if err != nil {
				return err
			}
		}
	}

	for _, info := range cfg.PriceInfos {
		//TODO battle
		_, err := r.db.Exec(insertSql, info.Symbol, info.TimeStamp, info.Price, info.PriceOrigin, info.Weight)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *coinHistoryRepository) GetTotalHistoryBySymbol(symbol string) (int, error) {
	var total int
	querySql := "select count(1) from `" + TABLE_COIN_PRICE + "` where symbol = ?;"
	logger.Infoln("sql:", querySql, "symbol", symbol)
	err := r.db.QueryRow(querySql, symbol).Scan(&total)
	if err != nil {
		return total, err
	}
	return total, nil
}

func (r *coinHistoryRepository) GetHistoryBySymbol(idx int, pageSize int, symbol string) ([]conf.PriceInfo, error) {
	var infos []conf.PriceInfo
	querySql := "select symbol, timestamp, price, weight, price_origin from `" + TABLE_COIN_PRICE + "` where symbol = ? order by id desc limit ?,? ;"
	logger.Infoln("sql:", querySql, "symbol", symbol, " limit:", strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))

	err := r.db.Select(&infos, querySql, symbol, strconv.Itoa(idx*pageSize), strconv.Itoa(pageSize))
	if err != nil {
		return infos, err
	}
	return infos, nil
}

func (r *coinHistoryRepository) GetHistoryBySymbolAndTimestamp(symbol string, timestamp int64) ([]conf.PriceInfo, error) {
	var infos []conf.PriceInfo
	querySql := "select symbol, timestamp, price, weight, price_origin from `" + TABLE_COIN_PRICE + "` where symbol = ? and timestamp = ? order by id desc ;"
	logger.Infoln("sql:", querySql, "symbol", symbol, "timestamp", timestamp)

	err := r.db.Select(&infos, querySql, symbol, timestamp)
	if err != nil {
		return infos, err
	}
	return infos, nil
}

func (r *coinHistoryRepository) GetHistoryByTimestamp(timestamp int64) ([]conf.PriceInfo, error) {

	dbPriceInfos := make([]conf.PriceInfo, 0)
	querySql := "select symbol, timestamp, price, weight, price_origin from `" + TABLE_COIN_PRICE + "` where timestamp = ?;"
	logger.Infoln("sql:", querySql, "timestamp", timestamp)
	err := r.db.Select(&dbPriceInfos, querySql, timestamp)
	if err != nil {
		return []conf.PriceInfo{}, err
	}

	return dbPriceInfos, nil
}
