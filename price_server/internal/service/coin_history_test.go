package service

import (
	"github.com/golang/mock/gomock"
	conf "price_api/price_server/config"
	"price_api/price_server/internal/cache"
	mock_cache "price_api/price_server/internal/cache/mock"
	"price_api/price_server/internal/repository"
	mock_repository "price_api/price_server/internal/repository/mock"
	"price_api/price_server/internal/vo"
	"reflect"
	"testing"
)

var (
	priceInfos1 = []conf.PriceInfo{
		{
			Symbol:      "btcusdt",
			Price:       47851.30000000,
			PriceOrigin: "ok",
			Weight:      1,
			TimeStamp:   1640745642,
		},
		{
			Symbol:      "btcusdt",
			Price:       47873.19000000,
			PriceOrigin: "bitstamp",
			Weight:      1,
			TimeStamp:   1640745642,
		},
		{
			Symbol:      "btcusdt",
			Price:       47851.04000000,
			PriceOrigin: "binance",
			Weight:      1,
			TimeStamp:   1640745642,
		},
		{
			Symbol:      "btcusdt",
			Price:       47852.37000000,
			PriceOrigin: "coinbase",
			Weight:      1,
			TimeStamp:   1640745642,
		},
		{
			Symbol:      "btcusdt",
			Price:       47846.10000000,
			PriceOrigin: "kucoin",
			Weight:      1,
			TimeStamp:   1640745642,
		},
		{
			Symbol:      "btcusdt",
			Price:       47853.17000000,
			PriceOrigin: "huobi",
			Weight:      2,
			TimeStamp:   1640745642,
		},
		{
			Symbol:      "btcusdt",
			Price:       47867.00000000,
			PriceOrigin: "bitfinex",
			Weight:      1,
			TimeStamp:   1640745642,
		},
	}
	exchangeConfigs = []conf.ExchangeConfig{
		{
			Name:   "binance",
			Weight: 1,
			Url:    "https://api.binance.com/api/v3/ticker/price?symbol={$symbol}",
		},
		{
			Name:   "huobi",
			Weight: 1,
			Url:    "https://api.huobi.pro/market/detail/merged?symbol={$symbol}",
		},
		{
			Name:   "bitfinex",
			Weight: 1,
			Url:    "https://api-pub.bitfinex.com/v2/tickers?symbols=t{$symbol}",
		},
		{
			Name:   "ok",
			Weight: 1,
			Url:    "https://www.okex.com/api/spot/v3/instruments/{$symbol1}-{$symbol2}/ticker",
		},
		{
			Name:   "coinbase",
			Weight: 1,
			Url:    "https://api.pro.coinbase.com/products/{$symbol}/book",
		},
		{
			Name:   "bitstamp",
			Weight: 1,
			Url:    "https://www.bitstamp.net/api/v2/ticker/{$symbol}",
		},
		{
			Name:   "kucoin",
			Weight: 1,
			Url:    "https://api.kucoin.com/api/v1/market/orderbook/level1?symbol={$symbol}",
		},
	}
	updatePriceHistory = vo.UpdatePirceHistory{
		Timestamp: 1639640386,
		Symbol:    "btcusdt",
	}
	priceInfo = conf.PriceInfo{
		Symbol:      "btcusdt",
		Price:       58609,
		PriceOrigin: "huobi",
		Weight:      2,
		TimeStamp:   1639640386,
	}
)

func TestCoinHistoryService_GetUpdatePriceHeartbeat(t *testing.T) {
	type fields struct {
		gPriceInfosCache   cache.GlobalPriceInfoCache
		gReqeustPriceConfs cache.GlobalRequestPriceConfs
		updatePriceRepo    repository.UpdatePriceRepository
		coinHistoryRepo    repository.CoinHistoryRepository
	}
	type args struct {
		symbol   string
		interval int64
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gPriceInfosCache := mock_cache.NewMockGlobalPriceInfoCache(ctrl)
	gReqeustPriceConfs := mock_cache.NewMockGlobalRequestPriceConfs(ctrl)
	updatePriceRepo := mock_repository.NewMockUpdatePriceRepository(ctrl)
	coinHistoryRepo := mock_repository.NewMockCoinHistoryRepository(ctrl)

	gPriceInfosCache.EXPECT().GetLatestPriceInfos(gomock.Eq("btcusdt")).Return(conf.PriceInfos{PriceInfos: priceInfos1})
	gReqeustPriceConfs.EXPECT().GetConfsBySymbol(gomock.Eq("btc-usdt")).Return(exchangeConfigs)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    vo.HEARTBEAT_INFO
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				gPriceInfosCache:   gPriceInfosCache,
				gReqeustPriceConfs: gReqeustPriceConfs,
				updatePriceRepo:    updatePriceRepo,
				coinHistoryRepo:    coinHistoryRepo,
			},
			args: args{symbol: "btcusdt", interval: 60},
			want: vo.HEARTBEAT_INFO{
				ExpectResources: 7,
				ActualResources: 7,
				LatestTimestamp: 1640745642,
				Interval:        60,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CoinHistoryService{
				gPriceInfosCache:   tt.fields.gPriceInfosCache,
				gReqeustPriceConfs: tt.fields.gReqeustPriceConfs,
				updatePriceRepo:    tt.fields.updatePriceRepo,
				coinHistoryRepo:    tt.fields.coinHistoryRepo,
			}
			got, err := s.GetUpdatePriceHeartbeat(tt.args.symbol, tt.args.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUpdatePriceHeartbeat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUpdatePriceHeartbeat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinHistoryService_GetUpdatePriceHistory(t *testing.T) {
	type fields struct {
		gPriceInfosCache   cache.GlobalPriceInfoCache
		gReqeustPriceConfs cache.GlobalRequestPriceConfs
		updatePriceRepo    repository.UpdatePriceRepository
		coinHistoryRepo    repository.CoinHistoryRepository
	}
	type args struct {
		idx      int
		pageSize int
		symbol   string
	}

	args1 := args{
		idx:      0,
		pageSize: 20,
		symbol:   "btcusdt",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gPriceInfosCache := mock_cache.NewMockGlobalPriceInfoCache(ctrl)
	gReqeustPriceConfs := mock_cache.NewMockGlobalRequestPriceConfs(ctrl)
	updatePriceRepo := mock_repository.NewMockUpdatePriceRepository(ctrl)
	coinHistoryRepo := mock_repository.NewMockCoinHistoryRepository(ctrl)

	updatePriceRepo.EXPECT().GetUpdatePriceHistoryBySymbol(gomock.Eq(args1.idx), gomock.Eq(args1.pageSize), gomock.Eq(args1.symbol)).Return([]vo.UpdatePirceHistory{updatePriceHistory}, nil)
	updatePriceRepo.EXPECT().GetTotalUpdatePriceHistoryBySymbol(gomock.Eq(args1.symbol)).Return(1, nil)
	coinHistoryRepo.EXPECT().GetHistoryBySymbolAndTimestamp(gomock.Eq(updatePriceHistory.Symbol), gomock.Eq(updatePriceHistory.Timestamp)).Return([]conf.PriceInfo{priceInfo}, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		want1   []vo.UpdatePriceHistoryResp
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				gPriceInfosCache:   gPriceInfosCache,
				gReqeustPriceConfs: gReqeustPriceConfs,
				updatePriceRepo:    updatePriceRepo,
				coinHistoryRepo:    coinHistoryRepo,
			},
			args: args1,
			want: 1,
			want1: []vo.UpdatePriceHistoryResp{{
				Timestamp: updatePriceHistory.Timestamp,
				Symbol:    updatePriceHistory.Symbol,
				Price:     priceInfo.Price,
				Infos:     []conf.PriceInfo{priceInfo},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CoinHistoryService{
				gPriceInfosCache:   tt.fields.gPriceInfosCache,
				gReqeustPriceConfs: tt.fields.gReqeustPriceConfs,
				updatePriceRepo:    tt.fields.updatePriceRepo,
				coinHistoryRepo:    tt.fields.coinHistoryRepo,
			}
			got, got1, err := s.GetUpdatePriceHistory(tt.args.idx, tt.args.pageSize, tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUpdatePriceHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUpdatePriceHistory() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetUpdatePriceHistory() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
