package vo

import conf "price_api/price_server/config"

type RESPONSE_PRICE_CONF struct {
	Price  float64
	Conf   conf.ExchangeConfig
	Symbol string
}
