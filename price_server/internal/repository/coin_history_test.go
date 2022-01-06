package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"price_api/price_server/internal/config"
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

var (
	priceInfo = conf.PriceInfo{
		Symbol:      "btcusd",
		Price:       58609,
		PriceOrigin: "huobi",
		Weight:      2,
		TimeStamp:   1640330341,
	}
)

func TestCoinHistoryRepository_GetHistoryBySymbol(t *testing.T) {

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		idx      int
		pageSize int
		symbol   string
	}

	args1 := args{
		idx:      0,
		pageSize: 20,
		symbol:   "btcusd",
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	querySql := "select symbol, timestamp, price, weight, price_origin from `" + TABLE_COIN_PRICE + "` where symbol = ? order by id desc limit ?,? ;"

	rows := sqlmock.NewRows([]string{"symbol", "timestamp", "price", "weight", "price_origin"}).
		AddRow(priceInfo.Symbol, priceInfo.TimeStamp, priceInfo.Price, priceInfo.Weight, priceInfo.PriceOrigin)

	mock.ExpectQuery(regexp.QuoteMeta(querySql)).WithArgs(args1.symbol, strconv.Itoa(args1.idx*args1.pageSize), strconv.Itoa(args1.pageSize)).WillReturnRows(rows)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []conf.PriceInfo
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{DB: func() *sqlx.DB {
				return db
			}()},
			args:    args1,
			want:    []conf.PriceInfo{priceInfo},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &coinHistoryRepository{
				db: tt.fields.DB,
			}
			got, err := r.GetHistoryBySymbol(tt.args.idx, tt.args.pageSize, tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHistoryBySymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHistoryBySymbol() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinHistoryRepository_GetHistoryBySymbolAndTimestamp(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		symbol    string
		timestamp int64
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	args1 := args{
		symbol:    "btcusd",
		timestamp: 1640330341,
	}

	querySql := "select symbol, timestamp, price, weight, price_origin from `" + TABLE_COIN_PRICE + "` where symbol = ? and timestamp = ? order by id desc ;"

	rows := sqlmock.NewRows([]string{"symbol", "timestamp", "price", "weight", "price_origin"}).
		AddRow(priceInfo.Symbol, priceInfo.TimeStamp, priceInfo.Price, priceInfo.Weight, priceInfo.PriceOrigin)

	mock.ExpectQuery(regexp.QuoteMeta(querySql)).WithArgs(args1.symbol, args1.timestamp).WillReturnRows(rows)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []conf.PriceInfo
		wantErr bool
	}{
		{
			name:   "basic",
			fields: fields{DB: db},
			args:   args1,
			want: []conf.PriceInfo{
				priceInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &coinHistoryRepository{
				db: tt.fields.DB,
			}
			got, err := r.GetHistoryBySymbolAndTimestamp(tt.args.symbol, tt.args.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHistoryBySymbolAndTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHistoryBySymbolAndTimestamp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinHistoryRepository_GetTotalHistoryBySymbol(t *testing.T) {
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

	args1 := args{symbol: priceInfo.Symbol}
	querySql := "select count(1) from `" + TABLE_COIN_PRICE + "` where symbol = ?;"
	rows := sqlmock.NewRows([]string{"count(1)"}).
		AddRow(10)
	mock.ExpectQuery(regexp.QuoteMeta(querySql)).WithArgs(priceInfo.Symbol).WillReturnRows(rows)

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
			r := &coinHistoryRepository{
				db: tt.fields.DB,
			}
			got, err := r.GetTotalHistoryBySymbol(tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTotalHistoryBySymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTotalHistoryBySymbol() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinHistoryRepository_InsertPriceInfo(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		cfg conf.PriceInfos
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	args1 := args{cfg: conf.PriceInfos{PriceInfos: []conf.PriceInfo{priceInfo}}}

	insertSql := "insert into " + TABLE_COIN_PRICE + " (symbol,timestamp,price,price_origin,weight)" +
		" values(?,?,?,?,?)"

	insertUpdateHistorySql := "insert into " + TABLE_UPDATE_PRICE_HISTORY + "(timestamp,symbol) value (?,?)"

	mock.ExpectExec(regexp.QuoteMeta(insertUpdateHistorySql)).
		WithArgs(priceInfo.TimeStamp, priceInfo.Symbol).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(insertSql)).
		WithArgs(priceInfo.Symbol, priceInfo.TimeStamp, priceInfo.Price, priceInfo.PriceOrigin, priceInfo.Weight).
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
			r := &coinHistoryRepository{
				db: tt.fields.DB,
			}
			if err := r.InsertPriceInfo(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("InsertPriceInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
