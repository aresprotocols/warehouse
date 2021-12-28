package service

import (
	"github.com/golang/mock/gomock"
	"price_api/price_server/internal/repository"
	"price_api/price_server/internal/vo"
	"reflect"
	"testing"
)

func TestHttpErrorService_GetHttpErrorsByPage(t *testing.T) {
	type fields struct {
		httpErrorRepo repository.HttpErrorRepository
	}
	type args struct {
		idx    int
		size   int
		symbol string
	}

	args1 := args{
		idx:    0,
		size:   10,
		symbol: "btcusdt",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	httpErrorRepo := repository.NewMockHttpErrorRepository(ctrl)

	httpErrorRepo.EXPECT().GetTotalHttpErrorInfo(gomock.Eq(args1.symbol)).Return(1, nil)
	//([]vo.HTTP_ERROR_INFO, error)
	httpErrorRepo.EXPECT().GetHttpErrorInfo(gomock.Eq(args1.idx), gomock.Eq(args1.size), gomock.Eq(args1.symbol)).Return([]vo.HTTP_ERROR_INFO{}, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		want1   []vo.HTTP_ERROR_INFO
		wantErr bool
	}{
		{
			name:    "basic",
			fields:  fields{httpErrorRepo: httpErrorRepo},
			args:    args1,
			want:    1,
			want1:   []vo.HTTP_ERROR_INFO{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HttpErrorService{
				httpErrorRepo: tt.fields.httpErrorRepo,
			}
			got, got1, err := s.GetHttpErrorsByPage(tt.args.idx, tt.args.size, tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHttpErrorsByPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetHttpErrorsByPage() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetHttpErrorsByPage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestHttpErrorService_InsertHttpError(t *testing.T) {
	type fields struct {
		httpErrorRepo repository.HttpErrorRepository
	}
	type args struct {
		url       string
		symbol    string
		errorInfo string
	}

	args1 := args{
		url:       "https://api-pub.bitfinex.com/v2/tickers?symbols=t{$symbol}",
		symbol:    "knc-usdt",
		errorInfo: "status code :429 url:https://api-pub.bitfinex.com/v2/tickers?symbols=tKNCUSD",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	httpErrorRepo := repository.NewMockHttpErrorRepository(ctrl)

	httpErrorRepo.EXPECT().InsertHttpError(gomock.Eq(args1.url), gomock.Eq(args1.symbol), gomock.Eq(args1.errorInfo)).Return(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "basic",
			fields:  fields{httpErrorRepo: httpErrorRepo},
			args:    args1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HttpErrorService{
				httpErrorRepo: tt.fields.httpErrorRepo,
			}
			if err := s.InsertHttpError(tt.args.url, tt.args.symbol, tt.args.errorInfo); (err != nil) != tt.wantErr {
				t.Errorf("InsertHttpError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
