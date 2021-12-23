package vo

type RESPONSE struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Pagination struct {
	CurPage  int         `json:"curPage"`
	TotalNum int         `json:"totalNum"`
	Items    interface{} `json:"items"`
}
