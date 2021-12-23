package vo

type HTTP_ERROR_INFO struct {
	Url       string `db:"url" json:"url"`
	Symbol    string `db:"symbol" json:"symbol"`
	Error     string `db:"error" json:"error"`
	Timestamp int64  `db:"timestamp" jsoon:"timestamp"`
}
