package vo

import conf "price_api/price_server/config"

type PriceAllInfo struct {
	Name      string  `json:"name"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
	Weight    int64   `json:"weight"`
}

type WeightInfo struct {
	Price        float64 `json:"price"`
	Weight       int64   `json:"weight"`
	ExchangeName string  `json:"exchangeName"`
}

type PartyPriceInfo struct {
	Price     float64      `json:"price"`
	Timestamp int64        `json:"timestamp"`
	Infos     []WeightInfo `json:"infos"`
}

type PRICE_INFO struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type PRICE_EXCHANGE_INFO struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
	Exchange  string  `json:"exchange"`
	Weight    int64   `json:"weight"`
}

type PRICE_EXCHANGE_WEIGHT_INFO struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
	Exchange  string  `json:"exchange"`
	Weight    int     `json:"weight"`
}

type CLIENT_INFO struct {
	Ip               string `json:"ip"`
	RequestTime      string `json:"request_time"`
	RequestTimestamp int64  `json:"request_timestamp"`
}

type CLIENT_PRICE_INFO struct {
	Client    CLIENT_INFO `json:"client"`
	PriceInfo PRICE_INFO  `json:"price_info"`
}

type CLIENT_PRICEALL_INFO struct {
	Client     CLIENT_INFO           `json:"client"`
	PriceInfos []PRICE_EXCHANGE_INFO `json:"price_infos"`
}

type PARTY_PRICE_INFO struct {
	Type       string                       `json:"type"`
	Client     CLIENT_INFO                  `json:"client"`
	PriceInfo  PRICE_INFO                   `json:"price_info"`
	PriceInfos []PRICE_EXCHANGE_WEIGHT_INFO `json:"price_infos"`
}

type UpdatePriceHistoryResp struct {
	Timestamp int64            `json:"timestamp"`
	Symbol    string           `json:"symbol"`
	Price     float64          `json:"price"`
	Infos     []conf.PriceInfo `json:"infos"`
}
