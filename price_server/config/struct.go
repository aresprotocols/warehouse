package conf

type PriceInfo struct {
	Symbol      string  `db:"symbol"`
	Price       float64 `db:"price"`
	PriceOrigin string  `db:"price_origin"`
	Weight      int64   `db:"weight"`
	TimeStamp   int64   `db:"timestamp"`
}

type PriceInfos struct {
	PriceInfos []PriceInfo
}

type PriceInfosCache struct {
	PriceInfosCache []PriceInfos
}
