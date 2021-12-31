package conf

type PriceInfo struct {
	Symbol      string  `json:"symbol" db:"symbol" `
	Price       float64 `json:"price" db:"price"`
	PriceOrigin string  `json:"priceOrigin" db:"price_origin"`
	Weight      int64   `json:"weight" db:"weight"`
	TimeStamp   int64   `json:"timestamp" db:"timestamp"`
}

type PriceInfos struct {
	PriceInfos []PriceInfo
}

type PriceInfosCache struct {
	PriceInfosCache map[string][]PriceInfos
}
