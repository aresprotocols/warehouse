package exchange

import (
	"errors"
	"strconv"
	"strings"
)

//[["tETHBTC",0.068768,60.457862539999994,0.068782,87.18505796000001,-0.000422,-0.0061,0.068768,4345.75880076,0.069731,0.0683]]
func parseBitfinexPrice(priceJson string) (float64, error) {
	firstIdx := strings.Index(priceJson, ",")
	if firstIdx == -1 {
		return 0, errors.New("unknow rsp format:" + priceJson)
	}

	lastIdx := strings.Index(priceJson[firstIdx+1:], ",")
	if firstIdx == -1 {
		return 0, errors.New("unknow rsp format:" + priceJson)
	}

	price, err := strconv.ParseFloat(priceJson[firstIdx+1:firstIdx+1+lastIdx], 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}
