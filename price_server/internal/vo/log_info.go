package vo

type LOG_INFO struct {
	ClientIP     string `db:"client_ip" json:"client_ip"`
	RequestTime  string `db:"request_time" json:"request_time"`
	UserAgent    string `db:"user_agent" json:"user_agent"`
	RequestUrl   string `db:"request_url" json:"request_url"`
	ResponseTime string `db:"response_time" json:"response_time"`
	Response     string `db:"request_response" json:"response"`
}

type LOG_INFOS struct {
	Infos []LOG_INFO `json:"infos"`
}

type REQ_RSP_LOG_INFO struct {
	ReqUrl           string `json:"reqUrl" db:"request_url"`
	Response         string `json:"response" db:"request_response"`
	Ip               string `json:"ip" db:"client_ip"`
	RequestTime      string `json:"request_time" db:"request_time"`
	RequestTimestamp int64  `json:"request_timestamp" db:"request_timestamp"`
}
