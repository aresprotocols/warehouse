package exchange

import (
	"encoding/json"
	"strconv"
)

//{"symbol":"BTCUSDT","price":"47653.01000000"}
type BinancePriceInfo struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func parseBinancePrice(priceJson string) (float64, error) {
	var binancePriceInfo BinancePriceInfo

	err := json.Unmarshal([]byte(priceJson), &binancePriceInfo)
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(binancePriceInfo.Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
