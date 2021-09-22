package exchange

import (
	"errors"
	"log"
	conf "price_api/price_server/config"
	"strings"
	"time"
)

func GetExchangePrice(cfg conf.Config) (conf.PriceInfos, error) {
	var retPriceInfos conf.PriceInfos

	timestamp := time.Now().Unix()
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
			priceInfo.Symbol = strings.Replace(symbol, "-", "", -1)
			priceInfo.PriceOrigin = exchange.Name
			priceInfo.TimeStamp = timestamp

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
			} else if lowName == "bitfinex" {
				price, err = parseBitfinexPrice(resJson)
				if err != nil {
					log.Println(err)
					continue
				}
			} else if lowName == "ok" {
				price, err = parseOkPrice(resJson)
				if err != nil {
					log.Println(err)
					continue
				}
			} else if lowName == "cryptocompare" {
				price, err = parseCryptoComparePrice(resJson)
				if err != nil {
					log.Println(resJson)
					log.Println(err)
					continue
				}
			} else if lowName == "coinbase" {
				price, err = parseCoinbasePrice(resJson)
				if err != nil {
					log.Println(err)
					continue
				}
			} else if lowName == "bitstamp" {
				price, err = parseBitStampPrice(resJson)
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
			symbol = strings.ReplaceAll(symbol, "-", "")
			return getPrice(strings.Replace(url, "{$symbol}", strings.ToUpper(symbol), -1), proxy)
		} else {
			return "", errors.New("symbol not find in binance url")
		}
	} else if lowName == "huobi" {
		if strings.Contains(url, "{$symbol}") {
			symbol = strings.ReplaceAll(symbol, "-", "")
			return getPrice(strings.Replace(url, "{$symbol}", strings.ToLower(symbol), -1), proxy)
		} else {
			return "", errors.New("symbol not find in huobi url")
		}
	} else if lowName == "bitfinex" {
		if strings.Contains(url, "{$symbol}") {
			symbol = strings.ReplaceAll(symbol, "-", "")
			if strings.Contains(symbol, "usdt") {
				symbol = strings.Replace(symbol, "usdt", "usd", -1)
			}
			return getPrice(strings.Replace(url, "{$symbol}", strings.ToUpper(symbol), -1), proxy)
		} else {
			return "", errors.New("symbol not find in bitfinex url")
		}
	} else if lowName == "ok" {
		if strings.Contains(url, "{$symbol1}") && strings.Contains(url, "{$symbol2}") {
			idx := strings.Index(symbol, "-")
			symbol1 := symbol[0:idx]
			symbol2 := symbol[idx+1:]
			url = strings.Replace(url, "{$symbol1}", symbol1, -1)
			url = strings.Replace(url, "{$symbol2}", symbol2, -1)

			return getPrice(url, proxy)
		} else {
			return "", errors.New("symbol not find in ok url")
		}
	} else if lowName == "cryptocompare" {
		if strings.Contains(url, "{$symbol1}") && strings.Contains(url, "{$symbol2}") {
			idx := strings.Index(symbol, "-")
			symbol1 := symbol[0:idx]
			symbol2 := symbol[idx+1:]
			url = strings.Replace(url, "{$symbol1}", symbol1, -1)
			url = strings.Replace(url, "{$symbol2}", symbol2, -1)

			return getPrice(url, proxy)
		} else {
			return "", errors.New("symbol not find in cryptocompare url")
		}
	} else if lowName == "coinbase" {
		if strings.Contains(url, "{$symbol}") {
			return getPrice(strings.Replace(url, "{$symbol}", strings.ToLower(symbol), -1), proxy)
		} else {
			return "", errors.New("symbol not find in coinbase url")
		}
	} else if lowName == "bitstamp" {
		if strings.Contains(url, "{$symbol}") {
			symbol = strings.ReplaceAll(symbol, "-", "")
			return getPrice(strings.Replace(url, "{$symbol}", strings.ToLower(symbol), -1), proxy)
		} else {
			return "", errors.New("symbol not find in bitstamp url")
		}
	} else {
		return "", errors.New("unknow exchangeName:" + exchangeName)
	}
}
