package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"price_api/price_server/internal/vo"
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

var (
	logInfo = vo.LOG_INFO{
		ClientIP:     "127.0.0.1",
		RequestTime:  "2021-12-16 15:39:46",
		UserAgent:    "insomnia/2021.7.1",
		RequestUrl:   "/api/getPrice/btcusdt/huobi",
		ResponseTime: "2021-12-16 15:39:46",
		Response:     "{\"code\":0,\"message\":\"OK\",\"data\":{\"price\":48903.4,\"timestamp\":1639640325}}",
	}

	reqRespLogInfo = vo.REQ_RSP_LOG_INFO{
		ReqUrl:           "/api/getPrice/btcusdt/huobi",
		Response:         "{\"code\":0,\"message\":\"OK\",\"data\":{\"price\":48903.4,\"timestamp\":1639640325}}",
		Ip:               "127.0.0.1",
		RequestTime:      "2021-12-16 15:39:46",
		RequestTimestamp: 1639640386,
	}
)

func TestLogInfoRepository_GetLogInfo(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}

	type args struct {
		idx      int
		pageSize int
	}

	args1 := args{
		idx:      0,
		pageSize: 20,
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	querySql := "select client_ip," +
		"request_time,user_agent,request_url,response_time,request_response from " +
		TABLE_LOG_INFO + " order by id desc limit ?,?;"

	rows := sqlmock.NewRows([]string{"client_ip", "request_time", "user_agent", "request_url", "response_time", "request_response"}).
		AddRow(logInfo.ClientIP, logInfo.RequestTime, logInfo.UserAgent, logInfo.RequestUrl, logInfo.ResponseTime, logInfo.Response)

	mock.ExpectQuery(regexp.QuoteMeta(querySql)).WithArgs(strconv.Itoa(args1.idx*args1.pageSize), strconv.Itoa(args1.pageSize)).WillReturnRows(rows)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    vo.LOG_INFOS
		wantErr bool
	}{
		{
			name:    "basic",
			fields:  fields{DB: db},
			args:    args1,
			want:    vo.LOG_INFOS{Infos: []vo.LOG_INFO{logInfo}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &logInfoRepository{
				DB: tt.fields.DB,
			}
			got, err := r.GetLogInfo(tt.args.idx, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLogInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogInfoRepository_GetLogInfoBySymbol(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		idx      int
		pageSize int
		symbol   string
		ip       string
	}

	symbol := "btcusdt"

	args1 := args{
		idx:      0,
		pageSize: 20,
		symbol:   symbol,
		ip:       "",
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	querySql := "select client_ip,request_url,request_time,request_response,request_timestamp from " +
		TABLE_LOG_INFO + " where ( request_response like '%" + symbol + "%'" +
		" or request_url like '%" + symbol + "%'" + " ) and use_symbol = 1 order by id desc limit ?,?;"

	rows := sqlmock.NewRows([]string{"client_ip", "request_url", "request_time", "request_response", "request_timestamp"}).
		AddRow(reqRespLogInfo.Ip, reqRespLogInfo.ReqUrl, reqRespLogInfo.RequestTime, reqRespLogInfo.Response, reqRespLogInfo.RequestTimestamp)

	mock.ExpectQuery(regexp.QuoteMeta(querySql)).WithArgs(strconv.Itoa(args1.idx*args1.pageSize), strconv.Itoa(args1.pageSize)).WillReturnRows(rows)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []vo.REQ_RSP_LOG_INFO
		wantErr bool
	}{
		{
			name:   "basic",
			fields: fields{DB: db},
			args:   args1,
			want: []vo.REQ_RSP_LOG_INFO{
				{
					ReqUrl:           reqRespLogInfo.ReqUrl,
					Response:         reqRespLogInfo.Response,
					Ip:               reqRespLogInfo.Ip,
					RequestTime:      reqRespLogInfo.RequestTime,
					RequestTimestamp: reqRespLogInfo.RequestTimestamp,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &logInfoRepository{
				DB: tt.fields.DB,
			}
			got, err := r.GetLogInfoBySymbol(tt.args.idx, tt.args.pageSize, tt.args.symbol, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogInfoBySymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLogInfoBySymbol() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogInfoRepository_GetTotalLogInfoBySymbol(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		symbol string
		ip     string
	}

	symbol := "btcusdt"

	args1 := args{symbol: symbol, ip: ""}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	querySql := "select count(1) from " +
		TABLE_LOG_INFO + " where ( request_response like '%" + symbol + "%'" +
		" or request_url like '%" + symbol + "%'" + " ) and use_symbol = 1"

	rows := sqlmock.NewRows([]string{"count(1)"}).AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta(querySql)).WillReturnRows(rows)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "basic",
			fields:  fields{DB: db},
			args:    args1,
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &logInfoRepository{
				DB: tt.fields.DB,
			}
			got, err := r.GetTotalLogInfoBySymbol(tt.args.symbol, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTotalLogInfoBySymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTotalLogInfoBySymbol() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogInfoRepository_InsertLogInfo(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}

	type args struct {
		mapInfo map[string]interface{}
		t       int
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	args1 := args{
		mapInfo: map[string]interface{}{
			"request_client_ip":  logInfo.ClientIP,
			"request_time":       logInfo.RequestTime,
			"request_ua":         logInfo.UserAgent,
			"request_uri":        logInfo.RequestUrl,
			"response_time":      logInfo.RequestTime,
			"response":           logInfo.Response,
			"request_timestamp":  reqRespLogInfo.RequestTimestamp,
			"response_timestamp": reqRespLogInfo.RequestTimestamp,
		},
		t: 1,
	}

	insertSql := "insert into " + TABLE_LOG_INFO + " (client_ip,request_time,user_agent,request_url," +
		"response_time,request_response, use_symbol,request_timestamp,response_timestamp)" +
		" values(?,?,?,?," +
		"?,?,?,?,?)"

	mock.ExpectExec(regexp.QuoteMeta(insertSql)).
		WithArgs(logInfo.ClientIP, logInfo.RequestTime, logInfo.UserAgent, logInfo.RequestUrl,
			logInfo.RequestTime, logInfo.Response, 1, reqRespLogInfo.RequestTimestamp, reqRespLogInfo.RequestTimestamp).
		WillReturnResult(sqlmock.NewResult(0, 1))

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "basic",
			fields:  fields{DB: db},
			args:    args1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &logInfoRepository{
				DB: tt.fields.DB,
			}
			if err := r.InsertLogInfo(tt.args.mapInfo, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("InsertLogInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
