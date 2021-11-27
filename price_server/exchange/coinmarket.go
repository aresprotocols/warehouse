package exchange

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
)

//{"status":{"timestamp":"2021-08-25T03:09:41.060Z","error_code":0,"error_message":null,"elapsed":23,"credit_count":1,"notice":null},"data":{"8702":{"id":8702,"name":"Ares Protocol","symbol":"ARES","slug":"ares-protocol","num_market_pairs":7,"date_added":"2021-03-06T00:00:00.000Z","tags":["oracles","smart-contracts","substrate","polkadot-ecosystem","duckstarter"],"max_supply":1000000000,"circulating_supply":153866975.90409997,"total_supply":1000000000,"platform":{"id":1027,"name":"Ethereum","symbol":"ETH","slug":"ethereum","token_address":"0x358AA737e033F34df7c54306960a38d09AaBd523"},"is_active":1,"cmc_rank":1100,"is_fiat":0,"last_updated":"2021-08-25T03:07:15.000Z","quote":{"USD":{"price":0.04314443190181,"volume_24h":757997.31058595,"percent_change_1h":0.27304841,"percent_change_24h":-4.82047334,"percent_change_7d":-6.93738111,"percent_change_30d":-11.27392811,"percent_change_60d":-28.08850153,"percent_change_90d":-70.30226632,"market_cap":6638503.263831881,"market_cap_dominance":0.0003,"fully_diluted_market_cap":43144431.9,"last_updated":"2021-08-25T03:07:15.000Z"}}}}}

const ARES_URL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?slug=ares-protocol&CMC_PRO_API_KEY=64a35d97-aca1-4c5c-8c17-43864a23aa97"
const ARES_Gate_URL = "https://data.gateapi.io/api2/1/ticker/ares_usdt"

type InfoDetail struct {
	CmcRank int `json:"cmc_rank"`
	Quote   struct {
		Usd struct {
			Price         float64 `json:"price"`
			Volume        float64 `json:"volume_24h"`
			PercentChange float64 `json:"percent_change_24h"`
			MarketCap     float64 `json:"market_cap"`
		} `json:"USD`
	} `json:"quote"`
}

type AresInfo struct {
	Status struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_message"` //maybe not string
	} `json:"status"`
	Data struct {
		ID InfoDetail `json:"8702"`
	} `json:"data"`
}

type AresGateShowInfo struct {
	Price         string `json:"last"`
	PercentChange string `json:"percentChange"`
	Rank          int    `json:"rank"`
	MarketCap     string `json:"quoteVolume"`
	Volume        string `json:"baseVolume"`
}

type AresShowInfo struct {
	Price         float64 `json:"price"`
	PercentChange float64 `json:"percent_change"`
	Rank          int     `json:"rank"`
	MarketCap     float64 `json:"market_cap"`
	Volume        float64 `json:"volume"`
}

func GetAresInfo(proxy string) (AresShowInfo, error) {
	var aresInfo AresInfo

	resJson, err := getPrice(ARES_URL, proxy)
	if err != nil {
		log.Println(err)
		return AresShowInfo{}, err
	}

	err = json.Unmarshal([]byte(resJson), &aresInfo)
	if err != nil {
		return AresShowInfo{}, err
	}

	if aresInfo.Status.ErrorCode != 0 {
		return AresShowInfo{}, errors.New(aresInfo.Status.ErrorMsg)
	}

	aresShowInfo := AresShowInfo{Price: aresInfo.Data.ID.Quote.Usd.Price,
		PercentChange: aresInfo.Data.ID.Quote.Usd.PercentChange,
		Rank:          aresInfo.Data.ID.CmcRank,
		MarketCap:     aresInfo.Data.ID.Quote.Usd.MarketCap,
		Volume:        aresInfo.Data.ID.Quote.Usd.Volume}

	return aresShowInfo, nil
}

func GetGateAresInfo(proxy string) (AresShowInfo, error) {
	var aresInfo AresGateShowInfo

	resJson, err := getPrice(ARES_Gate_URL, proxy)
	if err != nil {
		log.Println(err)
		return AresShowInfo{}, err
	}

	err = json.Unmarshal([]byte(resJson), &aresInfo)
	if err != nil {
		return AresShowInfo{}, err
	}
	price, _ := strconv.ParseFloat(aresInfo.Price, 64)
	percentChange, _ := strconv.ParseFloat(aresInfo.PercentChange, 64)
	marketCap, _ := strconv.ParseFloat(aresInfo.MarketCap, 64)
	volume, _ := strconv.ParseFloat(aresInfo.Volume, 64)

	return AresShowInfo{
		Price:         price,
		PercentChange: percentChange,
		Rank:          1427,
		MarketCap:     marketCap,
		Volume:        volume,
	}, nil
}
