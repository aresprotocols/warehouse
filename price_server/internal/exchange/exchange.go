package exchange

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	conf2 "price_api/price_server/internal/config"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
	"strings"
	"time"
)

var gResPriceConfCH = make(chan vo.RESPONSE_PRICE_CONF)

func GetSymbolExchangePrice(symbol string, reqConf map[string][]conf2.ExchangeConfig, cfg conf2.Config) (conf2.PriceInfos, error) {
	var retPriceInfos conf2.PriceInfos
	timestamp := time.Now().Unix()
	var gCh = make(chan conf2.PriceInfo)
	reqCount := 0
	for _, exchangeConf := range reqConf[symbol] {
		reqCount++
		go getPriceInfo(gCh, exchangeConf, symbol, cfg, reqConf)
	}
	for i := 0; i < reqCount; i++ {
		priceInfo := <-gCh
		priceInfo.TimeStamp = timestamp
		if priceInfo.Price != 0 {
			retPriceInfos.PriceInfos = append(retPriceInfos.PriceInfos, priceInfo)
		}
	}
	// end := time.Now().Unix()
	// log.Println("cost time:", end-timestamp)
	return retPriceInfos, nil
}

func getPriceInfo(gCh chan conf2.PriceInfo, exchange conf2.ExchangeConfig, symbol string, cfg conf2.Config, reqConf map[string][]conf2.ExchangeConfig) {
	var priceInfo conf2.PriceInfo
	defer func() {
		gCh <- priceInfo
	}()

	infos := reqConf[symbol]
	for _, info := range infos {
		if info.Name == exchange.Name {
			priceInfo.Weight = info.Weight
			break
		}
	}

	priceInfo.Symbol = strings.Replace(symbol, "-", "", -1)
	priceInfo.PriceOrigin = exchange.Name

	priceInfo.Price = getPriceByConf(exchange, symbol, cfg, true)
}

func getPriceBySymbolExchange(url, symbol, exchangeName, proxy string, replaces map[string]map[string]string) (string, error) {
	lowName := strings.ToLower(exchangeName)
	if len(replaces[symbol]) != 0 {
		if _, ok := replaces[symbol][lowName]; ok {
			symbol = replaces[symbol][lowName]
		}
	}
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
	} else if lowName == "kucoin" {
		if strings.Contains(url, "{$symbol}") {
			return getPrice(strings.Replace(url, "{$symbol}", strings.ToUpper(symbol), -1), proxy)
		} else {
			return "", errors.New("symbol not find in bitstamp url")
		}
	} else {
		return "", errors.New("unknow exchangeName:" + exchangeName)
	}
}

func InitSymbolsUpdateInterval(cfg conf2.Config) error {
	updateIntervalService := service.Svc.UpdateInterval()
	for _, symbol := range cfg.Symbols {
		interval, err := updateIntervalService.CheckUpdateInterval(symbol, int(cfg.Interval))
		if err != nil {
			return err
		}
		updateIntervalService.SetUpdateInterval(symbol, interval)
	}
	return nil
}

func InitRequestPriceConf(cfg conf2.Config) (map[string][]conf2.ExchangeConfig, error) {
	retRequestPriceConf := make(map[string][]conf2.ExchangeConfig)

	for _, exchange := range cfg.Exchanges {
		for _, symbol := range cfg.Symbols {
			go initRequestPrice(exchange, symbol, cfg)
			if exchange.Name == "coinbase" {
				time.Sleep(110 * time.Millisecond)
			}
			if exchange.Name == "bitfinex" {
				time.Sleep(110 * time.Millisecond)
			}
		}
	}

	for i := 0; i < len(cfg.Exchanges)*len(cfg.Symbols); i++ {
		resPriceConf := <-gResPriceConfCH
		if resPriceConf.Price != 0 {
			retRequestPriceConf[resPriceConf.Symbol] = append(retRequestPriceConf[resPriceConf.Symbol], resPriceConf.Conf)
		}
	}

	weightInfoService := service.Svc.WeightInfo()

	for symbol, configs := range retRequestPriceConf {
		for i, config := range configs {
			weight, err := weightInfoService.CheckUpdateWeight(symbol, config.Name, config.Weight)
			if err != nil {
				return retRequestPriceConf, err
			}
			retRequestPriceConf[symbol][i].Weight = weight
		}
	}

	return retRequestPriceConf, nil
}

func initRequestPrice(exchange conf2.ExchangeConfig, symbol string, cfg conf2.Config) {
	resPriceConf := vo.RESPONSE_PRICE_CONF{Conf: exchange, Symbol: symbol}
	defer func() {
		gResPriceConfCH <- resPriceConf
	}()

	resPriceConf.Price = getPriceByConf(exchange, symbol, cfg, false)
}

func getPriceByConf(exchange conf2.ExchangeConfig, symbol string, cfg conf2.Config, bRemberDb bool) float64 {
	var resJson string
	var err error

	lowName := strings.ToLower(exchange.Name)
	for i := 0; i < int(cfg.RetryCount); i++ {
		resJson, err = getPriceBySymbolExchange(exchange.Url, symbol, exchange.Name, cfg.Proxy, cfg.SymbolReplaces)
		if err == nil {
			break
		}
		if err != nil && strings.Contains(err.Error(), "404") { // skip retry when catch 404 error
			break
		}
		time.Sleep(time.Second * 3)
	}

	httpErrorService := service.Svc.HttpError()

	if err != nil {
		logger.WithError(err).Errorf("get price by symbol exchange error,symbol:%s,exchange:%s", symbol, exchange.Name)
		if bRemberDb {
			err = httpErrorService.InsertHttpError(exchange.Url, symbol, err.Error())
			if err != nil {
				logger.WithError(err).Errorf("insert http error to db occur error")
			}
		}
		return 0
	}

	//add price
	var price float64
	if lowName == "binance" {
		price, err = parseBinancePrice(resJson)
		if err != nil {
			logger.Errorln("response:", resJson, " err:", err)
			return 0
		}
	} else if lowName == "huobi" {
		price, err = parseHuobiPrice(resJson)
		if err != nil {
			logger.Errorln("response:", resJson, " err:", err)
			return 0
		}
	} else if lowName == "bitfinex" {
		price, err = parseBitfinexPrice(resJson)
		if err != nil {
			logger.Errorln("response:", resJson, " err:", err)
			return 0
		}
	} else if lowName == "ok" {
		price, err = parseOkPrice(resJson)
		if err != nil {
			logger.Errorln("response:", resJson, " err:", err)
			return 0
		}
	} else if lowName == "cryptocompare" {
		price, err = parseCryptoComparePrice(resJson)
		if err != nil {
			logger.Errorln("response:", resJson, " err:", err)
			return 0
		}
	} else if lowName == "coinbase" {
		price, err = parseCoinbasePrice(resJson)
		if err != nil {
			logger.Errorln("response:", resJson, " err:", err)
			return 0
		}
	} else if lowName == "bitstamp" {
		price, err = parseBitStampPrice(resJson)
		if err != nil {
			logger.Errorln("response:", resJson, " err:", err)
			return 0
		}
	} else if lowName == "kucoin" {
		price, err = parseKucoinPrice(resJson)
		if err != nil {
			logger.Errorln("response:", resJson, " err:", err)
			return 0
		}
	} else {
		logger.Errorf("unknown exchange,symbol:%s,exchange:%s, response:%s", symbol, exchange.Name, resJson)
		return 0
	}
	return price
}
