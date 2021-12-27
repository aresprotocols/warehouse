package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"regexp"
	"testing"
)

func TestWeightInfoRepository_CheckUpdateWeight(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		symbol       string
		exchangeName string
		weight       int64
	}

	args1 := args{
		symbol:       "btcusdt",
		exchangeName: "huobi",
		weight:       2,
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()
	querySql := "select weight from " + TABLE_WEIGH_INFO + " where symbol = ? and exchange = ?"
	rows := sqlmock.NewRows([]string{"weight"}).
		AddRow(args1.weight)
	mock.ExpectQuery(regexp.QuoteMeta(querySql)).WithArgs(args1.symbol, args1.exchangeName).WillReturnRows(rows)

	db2, mock2 := NewMock()
	defer func() {
		db2.Close()
	}()
	rows2 := sqlmock.NewRows([]string{"weight"})
	mock2.ExpectQuery(regexp.QuoteMeta(querySql)).WithArgs(args1.symbol, args1.exchangeName).WillReturnRows(rows2)
	insertSql := "insert into " + TABLE_WEIGH_INFO + " (symbol,exchange,weight)" +
		" values(?,?,?)"
	mock2.ExpectExec(regexp.QuoteMeta(insertSql)).
		WithArgs(args1.symbol, args1.exchangeName, args1.weight).
		WillReturnResult(sqlmock.NewResult(0, 1))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "read record",
			fields:  fields{DB: db},
			args:    args1,
			want:    args1.weight,
			wantErr: false,
		},
		{
			name:    "insert record",
			fields:  fields{DB: db2},
			args:    args1,
			want:    args1.weight,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &WeightInfoRepository{
				DB: tt.fields.DB,
			}
			got, err := r.CheckUpdateWeight(tt.args.symbol, tt.args.exchangeName, tt.args.weight)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUpdateWeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckUpdateWeight() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeightInfoRepository_SetWeight(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		symbol       string
		exchangeName string
		weight       int
	}

	args1 := args{
		symbol:       "btcusdt",
		exchangeName: "huobi",
		weight:       2,
	}

	db, mock := NewMock()
	defer func() {
		db.Close()
	}()
	querySql := "select weight from " + TABLE_WEIGH_INFO + " where symbol = ? and exchange = ?"
	rows := sqlmock.NewRows([]string{"weight"}).
		AddRow(args1.weight)
	mock.ExpectQuery(regexp.QuoteMeta(querySql)).WithArgs(args1.symbol, args1.exchangeName).WillReturnRows(rows)

	updateSql := "update " + TABLE_WEIGH_INFO + " set weight = ? where symbol = ? and exchange = ?"

	mock.ExpectExec(regexp.QuoteMeta(updateSql)).
		WithArgs(args1.weight, args1.symbol, args1.exchangeName).
		WillReturnResult(sqlmock.NewResult(0, 1))

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "basic",
			fields:  fields{db},
			args:    args1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &WeightInfoRepository{
				DB: tt.fields.DB,
			}
			if err := r.SetWeight(tt.args.symbol, tt.args.exchangeName, tt.args.weight); (err != nil) != tt.wantErr {
				t.Errorf("SetWeight() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
