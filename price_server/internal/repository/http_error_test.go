package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"price_api/price_server/internal/vo"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"
)

var (
	httpErrorInfo = vo.HTTP_ERROR_INFO{
		Url:       "https://api-pub.bitfinex.com/v2/tickers?symbols=t{$symbol}",
		Symbol:    "knc-usdt",
		Error:     "status code :429 url:https://api-pub.bitfinex.com/v2/tickers?symbols=tKNCUSD",
		Timestamp: 1639641289,
	}
)

func TestHttpErrorRepository_GetHttpErrorInfo(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		idx      int
		symbol   string
		pageSize int
	}

	args1 := args{
		idx:      0,
		symbol:   httpErrorInfo.Symbol,
		pageSize: 20,
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	querySql := "select url,symbol,error,timestamp from " +
		TABLE_HTTP_ERROR + " where symbol = ? order by id desc limit ?,?;"

	rows := sqlmock.NewRows([]string{"url", "symbol", "error", "timestamp"}).
		AddRow(httpErrorInfo.Url, httpErrorInfo.Symbol, httpErrorInfo.Error, httpErrorInfo.Timestamp)

	mock.ExpectQuery(regexp.QuoteMeta(querySql)).WithArgs(args1.symbol, strconv.Itoa(args1.idx*args1.pageSize), strconv.Itoa(args1.pageSize)).WillReturnRows(rows)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []vo.HTTP_ERROR_INFO
		wantErr bool
	}{
		{
			name:    "basic",
			fields:  fields{DB: db},
			args:    args1,
			want:    []vo.HTTP_ERROR_INFO{httpErrorInfo},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &HttpErrorRepository{
				DB: tt.fields.DB,
			}
			got, err := r.GetHttpErrorInfo(tt.args.idx, tt.args.symbol, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHttpErrorInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHttpErrorInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpErrorRepository_GetTotalHttpErrorInfo(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		symbol string
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	args1 := args{symbol: httpErrorInfo.Symbol}
	querySql := "select count(1) from " + TABLE_HTTP_ERROR + " where symbol = ?;"
	rows := sqlmock.NewRows([]string{"count(1)"}).
		AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta(querySql)).WithArgs(httpErrorInfo.Symbol).WillReturnRows(rows)

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
			r := &HttpErrorRepository{
				DB: tt.fields.DB,
			}
			got, err := r.GetTotalHttpErrorInfo(tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTotalHttpErrorInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTotalHttpErrorInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpErrorRepository_InsertHttpError(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		url       string
		symbol    string
		errorInfo string
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	args1 := args{
		url:       httpErrorInfo.Url,
		symbol:    httpErrorInfo.Symbol,
		errorInfo: httpErrorInfo.Error,
	}

	insertSql := "insert into " + TABLE_HTTP_ERROR + " (url,symbol,error,timestamp)" +
		" values(?,?,?,?)"

	mock.ExpectExec(regexp.QuoteMeta(insertSql)).
		WithArgs(httpErrorInfo.Url, httpErrorInfo.Symbol, httpErrorInfo.Error, time.Now().Unix()).
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
			r := &HttpErrorRepository{
				DB: tt.fields.DB,
			}
			if err := r.InsertHttpError(tt.args.url, tt.args.symbol, tt.args.errorInfo); (err != nil) != tt.wantErr {
				t.Errorf("InsertHttpError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
