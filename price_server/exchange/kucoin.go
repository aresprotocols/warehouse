package exchange

import (
	"encoding/json"
	"errors"
	"strconv"
)

func parseKucoinPrice(priceJson string) (float64, error) {
	var kucoinPriceInfo map[string]interface{}

	err := json.Unmarshal([]byte(priceJson), &kucoinPriceInfo)
	if err != nil {
		return 0, err
	}
	if kucoinPriceInfo["code"] != "200000" {
		return 0, errors.New("some error")
	} else {
		data := kucoinPriceInfo["data"]
		if data == nil {
			return 0, nil
		}
		dataMap := data.(map[string]interface{})
		price, err := strconv.ParseFloat(dataMap["price"].(string), 64)
		if err != nil {
			return 0, err
		}

		return price, nil

	}
}
