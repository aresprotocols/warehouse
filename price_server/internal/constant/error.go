package constant

const MSG_URL_NOT_FIND = "url not find"
const MSG_PRICE_NOT_READY = "price not ready"
const MSG_PARAM_NOT_TRUE = "param not true"
const MSG_GET_ARES_ERROR = "get ares info error"
const MSG_PARSE_PARAM_ERROR = "parse param error"
const MSG_GET_LOG_INFO_ERROR = "get log info error"
const MSG_CHECK_USER_ERROR = "user and password not match"

const (
	ERROR = iota - 1000
	NO_MATCH_FORMAT_ERROR
	PARAM_NOT_TRUE_ERROR
	GET_ARES_INFO_ERROR
	PARSE_PARAM_ERROR
	GET_LOG_INFO_ERROR
	GET_HTTP_ERROR_ERROR
	CHECK_USER_ERROR
	SET_WEIGHT_ERROR
	SET_UPDATE_INTERVAL_ERROR
)
