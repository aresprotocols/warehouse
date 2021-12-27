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
	updatePriceHistory = vo.UpdatePirceHistory{
		Timestamp: 1639640386,
		Symbol:    "btcusdt",
	}
)

func TestUpdatePriceRepository_GetTotalUpdatePriceHistoryBySymbol(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		symbol string
	}

	args1 := args{symbol: updatePriceHistory.Symbol}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	querySql := "select count(1) from `" + TABLE_UPDATE_PRICE_HISTORY + "` where symbol = ?;"

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
			fields:  fields{db},
			args:    args1,
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UpdatePriceRepository{
				DB: tt.fields.DB,
			}
			got, err := r.GetTotalUpdatePriceHistoryBySymbol(tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTotalUpdatePriceHistoryBySymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTotalUpdatePriceHistoryBySymbol() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdatePriceRepository_GetUpdatePriceHistoryBySymbol(t *testing.T) {
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
		symbol:   updatePriceHistory.Symbol,
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()

	querySql := "select symbol, timestamp  from `" + TABLE_UPDATE_PRICE_HISTORY + "` where symbol = ? order by timestamp desc limit ?,? ;"

	rows := sqlmock.NewRows([]string{"symbol", "timestamp"}).
		AddRow(updatePriceHistory.Symbol, updatePriceHistory.Timestamp)

	mock.ExpectQuery(regexp.QuoteMeta(querySql)).WithArgs(args1.symbol, strconv.Itoa(args1.idx*args1.pageSize), strconv.Itoa(args1.pageSize)).WillReturnRows(rows)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []vo.UpdatePirceHistory
		wantErr bool
	}{
		{
			name:    "basic",
			fields:  fields{db},
			args:    args1,
			want:    []vo.UpdatePirceHistory{updatePriceHistory},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UpdatePriceRepository{
				DB: tt.fields.DB,
			}
			got, err := r.GetUpdatePriceHistoryBySymbol(tt.args.idx, tt.args.pageSize, tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUpdatePriceHistoryBySymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUpdatePriceHistoryBySymbol() got = %v, want %v", got, tt.want)
			}
		})
	}
}
