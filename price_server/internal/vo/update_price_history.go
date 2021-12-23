package vo

type UpdatePirceHistory struct {
	Timestamp int64  `json:"timestamp" db:"timestamp"`
	Symbol    string `json:"symbol" db:"symbol"`
}
