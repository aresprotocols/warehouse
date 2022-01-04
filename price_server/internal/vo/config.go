package vo

type ReqConfigResp struct {
	Weight   []EXCHANGE_WEIGHT_INFO `json:"weight"`
	Interval int                    `json:"interval"`
}
