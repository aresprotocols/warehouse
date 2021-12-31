package vo

type SetIntervalReq struct {
	Interval int    `json:"interval"`
	Symbol   string `json:"symbol"`
}
