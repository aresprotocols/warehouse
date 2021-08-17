package exchange

import (
	"encoding/json"
	"strconv"
)

//{"best_ask":"0.06883","best_bid":"0.06881","instrument_id":"ETH-BTC","open_utc0":"0.06856","open_utc8":"0.0696","product_id":"ETH-BTC","last":"0.06884","last_qty":"0.27244","ask":"0.06883","best_ask_size":"17.43","bid":"0.06881","best_bid_size":"21.95","open_24h":"0.06887","high_24h":"0.06975","low_24h":"0.06825","base_volume_24h":"6591.948345","timestamp":"2021-08-17T03:52:26.903Z","quote_volume_24h":"455.12237"}
type OkPriceInfo struct {
	Ask string `json:"best_ask"`
}

func parseOkPrice(priceJson string) (float64, error) {
	var okPriceInfo OkPriceInfo

	err := json.Unmarshal([]byte(priceJson), &okPriceInfo)
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(okPriceInfo.Ask, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
