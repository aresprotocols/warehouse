package exchange

import (
	"encoding/json"
	"strconv"
)

//{"bids":[["45848.63","0.20914554",1]],"asks":[["45866.2","0.052",1]],"sequence":822411170}
type CoinbasePriceInfo struct {
	Bids [][]interface{} `json:"bids"`
}

func parseCoinbasePrice(priceJson string) (float64, error) {
	var coinbasePriceInfo CoinbasePriceInfo

	err := json.Unmarshal([]byte(priceJson), &coinbasePriceInfo)
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(coinbasePriceInfo.Bids[0][0].(string), 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
