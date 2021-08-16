package exchange

import (
	"errors"
	"log"
	conf "price_api/price_server/config"
	"strings"
)

func GetExchangePrice(cfg conf.Config) (conf.PriceInfos, error) {
	var retPriceInfos conf.PriceInfos
	for _, exchange := range cfg.Exchanges {
		for _, symbol := range cfg.Symbols {
			var priceInfo conf.PriceInfo
			lowName := strings.ToLower(exchange.Name)
			resJson, err := getPriceBySymbolExchange(exchange.Url, symbol, exchange.Name, cfg.Proxy)
			if err != nil {
				log.Println(err)
				continue
			}

			priceInfo.Weight = exchange.Weight
			priceInfo.Symbol = symbol
			priceInfo.PriceOrigin = exchange.Name

			//add price
			var price float64
			if lowName == "binance" {
				price, err = parseBinancePrice(resJson)
				if err != nil {
					log.Println(err)
					continue
				}
			} else if lowName == "huobi" {
				price, err = parseHuobiPrice(resJson)
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				log.Println("unknow exchange name:", exchange.Name)
				continue
			}

			priceInfo.Price = price
			//log.Println(priceInfo)
			retPriceInfos.PriceInfos = append(retPriceInfos.PriceInfos, priceInfo)
		}
	}

	return retPriceInfos, nil
}

func getPriceBySymbolExchange(url, symbol, exchangeName, proxy string) (string, error) {
	lowName := strings.ToLower(exchangeName)
	if lowName == "binance" {
		if strings.Contains(url, "{$symbol}") {
			return getPrice(strings.Replace(url, "{$symbol}", strings.ToUpper(symbol), -1), proxy)
		} else {
			return "", errors.New("symbol not find in binance url")
		}
	} else if lowName == "huobi" {
		if strings.Contains(url, "{$symbol}") {
			return getPrice(strings.Replace(url, "{$symbol}", strings.ToLower(symbol), -1), proxy)
		} else {
			return "", errors.New("symbol not find in binance url")
		}
	}
	return "", nil
}
