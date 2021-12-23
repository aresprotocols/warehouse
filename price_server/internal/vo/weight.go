package vo

type SetWeightReq struct {
	Weight   int    `json:"weight"`
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
}
