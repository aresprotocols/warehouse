package exchange

import (
	"errors"
	"strconv"
	"strings"
)

//{"USD":45883.67}
func parseCryptoComparePrice(priceJson string) (float64, error) {
	idx := strings.Index(priceJson, ":")
	if idx == -1 {
		return 0, errors.New("unknow rsp format:" + priceJson)
	}

	price, err := strconv.ParseFloat(priceJson[idx+1:len(priceJson)-1], 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}
