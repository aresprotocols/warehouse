package exchange

import (
	"encoding/json"
	"strconv"
)

//{"high": "47596.81", "last": "45736.67", "timestamp": "1629184431", "bid": "45783.47", "vwap": "46471.38", "volume": "20.28131577", "low": "45326.39", "ask": "45823.07", "open": "45887.86"}
type BitStampPriceInfo struct {
	Last string `json:"last"`
}

func parseBitStampPrice(priceJson string) (float64, error) {
	var bitstampPriceInfo BitStampPriceInfo

	err := json.Unmarshal([]byte(priceJson), &bitstampPriceInfo)
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(bitstampPriceInfo.Last, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
