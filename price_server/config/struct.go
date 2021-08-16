package conf

type PriceInfo struct {
	Symbol      string
	Price       float64
	PriceOrigin string
	Weight      int64
}

type PriceInfos struct {
	PriceInfos []PriceInfo
}

type PriceInfosCache struct {
	PriceInfosCache []PriceInfos
}
