package vo

import conf "price_api/price_server/config"

type RESPONSE_PRICE_CONF struct {
	Price  float64
	Conf   conf.ExchangeConfig
	Symbol string
}

type EXCHANGE_WEIGHT_INFO struct {
	Exchange string `json:"exchange"`
	Weight   int64  `json:"weight"`
}
