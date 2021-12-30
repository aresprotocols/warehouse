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
	priceInfos2 = []conf.PriceInfo{
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
			Symbol:      "avaxusdt",
			Price:       108.51400000,
			PriceOrigin: "kucoin",
			Weight:      1,
			TimeStamp:   1640745642,
		},
		{
			Symbol:      "avaxusdt",
			Price:       108.58400000,
			PriceOrigin: "huobi",
			Weight:      1,
			TimeStamp:   1640745642,
		},
		{
			Symbol:      "avaxusdt",
			Price:       108.55000000,
			PriceOrigin: "coinbase",
			Weight:      1,
			TimeStamp:   1640745642,
		},
		{
			Symbol:      "avaxusdt",
			Price:       108.60000000,
			PriceOrigin: "binance",
			Weight:      1,
			TimeStamp:   1640745642,
		},
		{
			Symbol:      "avaxusdt",
			Price:       108.57900000,
			PriceOrigin: "ok",
			Weight:      1,
			TimeStamp:   1640745642,
		},
	}
)

func TestPriceService_GetBulkCurrencyPrices(t *testing.T) {
	type fields struct {
		gPriceInfosCache   cache.GlobalPriceInfoCache
		coinHistoryRepo    repository.CoinHistoryRepository
		gRequestPriceConfs cache.GlobalRequestPriceConfs
	}
	type args struct {
		symbol   string
		currency string
	}

	//args1 := args{
	//	symbol:   "btc_avax_ltc_bch_fil_etc_eos_dash_comp_matic",
	//	currency: "usdt",
	//}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gPriceInfosCache := mock_cache.NewMockGlobalPriceInfoCache(ctrl)
	gReqeustPriceConfs := mock_cache.NewMockGlobalRequestPriceConfs(ctrl)
	coinHistoryRepo := mock_repository.NewMockCoinHistoryRepository(ctrl)

	gPriceInfosCache.EXPECT().GetLatestPriceInfos().Return(conf.PriceInfos{PriceInfos: priceInfos2}).Times(2)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]vo.PartyPriceInfo
	}{
		{
			name: "only btc",
			fields: fields{
				gPriceInfosCache:   gPriceInfosCache,
				coinHistoryRepo:    coinHistoryRepo,
				gRequestPriceConfs: gReqeustPriceConfs,
			},
			args: args{
				symbol:   "btc",
				currency: "usdt",
			},
			want: map[string]vo.PartyPriceInfo{
				"btcusdt": {
					Price:     47851.835,
					Timestamp: 1640745642,
					Infos: []vo.WeightInfo{
						{
							Price:        47852.37,
							Weight:       1,
							ExchangeName: "coinbase",
						},
						{
							Price:        47851.30,
							Weight:       1,
							ExchangeName: "ok",
						},
					},
				},
			},
		},
		{
			name: " btc and avax",
			fields: fields{
				gPriceInfosCache:   gPriceInfosCache,
				coinHistoryRepo:    coinHistoryRepo,
				gRequestPriceConfs: gReqeustPriceConfs,
			},
			args: args{
				symbol:   "btc_avax",
				currency: "usdt",
			},
			want: map[string]vo.PartyPriceInfo{
				"btcusdt": {
					Price:     47851.835,
					Timestamp: 1640745642,
					Infos: []vo.WeightInfo{
						{
							Price:        47852.37,
							Weight:       1,
							ExchangeName: "coinbase",
						},
						{
							Price:        47851.30,
							Weight:       1,
							ExchangeName: "ok",
						},
					},
				},
				"avaxusdt": {
					Price:     108.571,
					Timestamp: 1640745642,
					Infos: []vo.WeightInfo{
						{
							Price:        108.584,
							Weight:       1,
							ExchangeName: "huobi",
						},
						{
							Price:        108.579,
							Weight:       1,
							ExchangeName: "ok",
						},
						{
							Price:        108.55,
							Weight:       1,
							ExchangeName: "coinbase",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PriceService{
				gPriceInfosCache:   tt.fields.gPriceInfosCache,
				coinHistoryRepo:    tt.fields.coinHistoryRepo,
				gRequestPriceConfs: tt.fields.gRequestPriceConfs,
			}
			if got := s.GetBulkCurrencyPrices(tt.args.symbol, tt.args.currency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBulkCurrencyPrices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceService_GetBulkPrices(t *testing.T) {
	type fields struct {
		gPriceInfosCache   *mock_cache.MockGlobalPriceInfoCache
		coinHistoryRepo    *mock_repository.MockCoinHistoryRepository
		gRequestPriceConfs *mock_cache.MockGlobalRequestPriceConfs
	}
	type args struct {
		symbol string
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    map[string]vo.PRICE_INFO
	}{
		{
			name: "only btc",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetLatestPriceInfos().Return(conf.PriceInfos{PriceInfos: priceInfos2})
			},
			args: args{symbol: "btcusdt"},
			want: map[string]vo.PRICE_INFO{
				"btcusdt": {
					Price:     47851.835,
					Timestamp: 1640745642,
				},
			},
		},
		{
			name: "btc and avax",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetLatestPriceInfos().Return(conf.PriceInfos{PriceInfos: priceInfos2})
			},
			args: args{symbol: "btcusdt_avaxusdt"},
			want: map[string]vo.PRICE_INFO{
				"btcusdt": {
					Price:     47851.835,
					Timestamp: 1640745642,
				},
				"avaxusdt": {
					Price:     108.571,
					Timestamp: 1640745642,
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				gPriceInfosCache:   mock_cache.NewMockGlobalPriceInfoCache(ctrl),
				coinHistoryRepo:    mock_repository.NewMockCoinHistoryRepository(ctrl),
				gRequestPriceConfs: mock_cache.NewMockGlobalRequestPriceConfs(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &PriceService{
				gPriceInfosCache:   f.gPriceInfosCache,
				coinHistoryRepo:    f.coinHistoryRepo,
				gRequestPriceConfs: f.gRequestPriceConfs,
			}
			if got := s.GetBulkPrices(tt.args.symbol); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBulkPrices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceService_GetBulkSymbolsState(t *testing.T) {
	type fields struct {
		gPriceInfosCache   *mock_cache.MockGlobalPriceInfoCache
		coinHistoryRepo    *mock_repository.MockCoinHistoryRepository
		gRequestPriceConfs *mock_cache.MockGlobalRequestPriceConfs
	}
	type args struct {
		symbolStr string
		currency  string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    map[string]bool
	}{
		{
			name: "btc",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetLatestPriceInfos().Return(conf.PriceInfos{PriceInfos: priceInfos2})
				f.gRequestPriceConfs.EXPECT().GetConfsBySymbol(gomock.Eq("btc-usdt")).Return(exchangeConfigs)
			},
			args: args{
				symbolStr: "btc",
				currency:  "usdt",
			},
			want: map[string]bool{
				"btcusdt": true,
			},
		},
		{
			name: "btc and avax",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetLatestPriceInfos().Return(conf.PriceInfos{PriceInfos: priceInfos2})
				f.gRequestPriceConfs.EXPECT().GetConfsBySymbol(gomock.Eq("btc-usdt")).Return(exchangeConfigs)
				f.gRequestPriceConfs.EXPECT().GetConfsBySymbol(gomock.Eq("avax-usdt")).Return(exchangeConfigs)
			},
			args: args{
				symbolStr: "btc_avax",
				currency:  "usdt",
			},
			want: map[string]bool{
				"btcusdt":  true,
				"avaxusdt": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				gPriceInfosCache:   mock_cache.NewMockGlobalPriceInfoCache(ctrl),
				coinHistoryRepo:    mock_repository.NewMockCoinHistoryRepository(ctrl),
				gRequestPriceConfs: mock_cache.NewMockGlobalRequestPriceConfs(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &PriceService{
				gPriceInfosCache:   f.gPriceInfosCache,
				coinHistoryRepo:    f.coinHistoryRepo,
				gRequestPriceConfs: f.gRequestPriceConfs,
			}
			if got := s.GetBulkSymbolsState(tt.args.symbolStr, tt.args.currency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBulkSymbolsState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceService_GetCacheLength(t *testing.T) {
	type fields struct {
		gPriceInfosCache   *mock_cache.MockGlobalPriceInfoCache
		coinHistoryRepo    *mock_repository.MockCoinHistoryRepository
		gRequestPriceConfs *mock_cache.MockGlobalRequestPriceConfs
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		want    int
	}{
		{
			name: "basic",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetCacheLength().Return(1)
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				gPriceInfosCache:   mock_cache.NewMockGlobalPriceInfoCache(ctrl),
				coinHistoryRepo:    mock_repository.NewMockCoinHistoryRepository(ctrl),
				gRequestPriceConfs: mock_cache.NewMockGlobalRequestPriceConfs(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &PriceService{
				gPriceInfosCache:   f.gPriceInfosCache,
				coinHistoryRepo:    f.coinHistoryRepo,
				gRequestPriceConfs: f.gRequestPriceConfs,
			}
			if got := s.GetCacheLength(); got != tt.want {
				t.Errorf("GetCacheLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceService_GetHistoryPrice(t *testing.T) {
	type fields struct {
		gPriceInfosCache   *mock_cache.MockGlobalPriceInfoCache
		coinHistoryRepo    *mock_repository.MockCoinHistoryRepository
		gRequestPriceConfs *mock_cache.MockGlobalRequestPriceConfs
	}
	type args struct {
		symbol    string
		timestamp int64
		bAverage  bool
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    bool
		want1   vo.PartyPriceInfo
	}{
		{
			name: "from cache",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetPriceInfosEqualTimestamp(gomock.Eq(int64(1640745642))).Return(true, conf.PriceInfos{PriceInfos: priceInfos1})
				//f.coinHistoryRepo.EXPECT().GetHistoryByTimestamp(gomock.Eq(1640745642)).Return(priceInfos1)
			},
			args: args{
				symbol:    "btcusdt",
				timestamp: 1640745642,
				bAverage:  true,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     47854.675,
				Timestamp: 1640745642,
				Infos: []vo.WeightInfo{
					{
						Price:        47867,
						Weight:       1,
						ExchangeName: "bitfinex",
					},
					{
						Price:        47853.17,
						Weight:       2,
						ExchangeName: "huobi",
					},
					{
						Price:        47852.37,
						Weight:       1,
						ExchangeName: "coinbase",
					},
					{
						Price:        47851.3,
						Weight:       1,
						ExchangeName: "ok",
					},
					{
						Price:        47851.04,
						Weight:       1,
						ExchangeName: "binance",
					},
				},
			},
		},
		{
			name: "from db",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetPriceInfosEqualTimestamp(gomock.Eq(int64(1640745642))).Return(false, conf.PriceInfos{})
				f.coinHistoryRepo.EXPECT().GetHistoryByTimestamp(gomock.Eq(int64(1640745642))).Return(priceInfos1, nil)
			},
			args: args{
				symbol:    "btcusdt",
				timestamp: 1640745642,
				bAverage:  true,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     47854.675,
				Timestamp: 1640745642,
				Infos: []vo.WeightInfo{
					{
						Price:        47867,
						Weight:       1,
						ExchangeName: "bitfinex",
					},
					{
						Price:        47853.17,
						Weight:       2,
						ExchangeName: "huobi",
					},
					{
						Price:        47852.37,
						Weight:       1,
						ExchangeName: "coinbase",
					},
					{
						Price:        47851.3,
						Weight:       1,
						ExchangeName: "ok",
					},
					{
						Price:        47851.04,
						Weight:       1,
						ExchangeName: "binance",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				gPriceInfosCache:   mock_cache.NewMockGlobalPriceInfoCache(ctrl),
				coinHistoryRepo:    mock_repository.NewMockCoinHistoryRepository(ctrl),
				gRequestPriceConfs: mock_cache.NewMockGlobalRequestPriceConfs(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &PriceService{
				gPriceInfosCache:   f.gPriceInfosCache,
				coinHistoryRepo:    f.coinHistoryRepo,
				gRequestPriceConfs: f.gRequestPriceConfs,
			}
			got, got1 := s.GetHistoryPrice(tt.args.symbol, tt.args.timestamp, tt.args.bAverage)
			if got != tt.want {
				t.Errorf("GetHistoryPrice() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetHistoryPrice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPriceService_GetLocalPrices(t *testing.T) {
	type fields struct {
		gPriceInfosCache   *mock_cache.MockGlobalPriceInfoCache
		coinHistoryRepo    *mock_repository.MockCoinHistoryRepository
		gRequestPriceConfs *mock_cache.MockGlobalRequestPriceConfs
	}
	type args struct {
		start  int
		end    int
		symbol string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    conf.PriceInfosCache
	}{
		{
			name: "btc",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetPriceInfosByRange(gomock.Eq(0), gomock.Eq(1)).Return(conf.PriceInfosCache{PriceInfosCache: []conf.PriceInfos{
					{PriceInfos: priceInfos2},
				}})
			},
			args: args{
				start:  0,
				end:    1,
				symbol: "btcusdt",
			},
			want: conf.PriceInfosCache{PriceInfosCache: []conf.PriceInfos{
				{PriceInfos: []conf.PriceInfo{
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
				}},
			}},
		},
		{
			name: "avax",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetPriceInfosByRange(gomock.Eq(0), gomock.Eq(1)).Return(conf.PriceInfosCache{PriceInfosCache: []conf.PriceInfos{
					{PriceInfos: priceInfos2},
				}})
			},
			args: args{
				start:  0,
				end:    1,
				symbol: "avaxusdt",
			},
			want: conf.PriceInfosCache{PriceInfosCache: []conf.PriceInfos{
				{PriceInfos: []conf.PriceInfo{
					{
						Symbol:      "avaxusdt",
						Price:       108.51400000,
						PriceOrigin: "kucoin",
						Weight:      1,
						TimeStamp:   1640745642,
					},
					{
						Symbol:      "avaxusdt",
						Price:       108.58400000,
						PriceOrigin: "huobi",
						Weight:      1,
						TimeStamp:   1640745642,
					},
					{
						Symbol:      "avaxusdt",
						Price:       108.55000000,
						PriceOrigin: "coinbase",
						Weight:      1,
						TimeStamp:   1640745642,
					},
					{
						Symbol:      "avaxusdt",
						Price:       108.60000000,
						PriceOrigin: "binance",
						Weight:      1,
						TimeStamp:   1640745642,
					},
					{
						Symbol:      "avaxusdt",
						Price:       108.57900000,
						PriceOrigin: "ok",
						Weight:      1,
						TimeStamp:   1640745642,
					},
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				gPriceInfosCache:   mock_cache.NewMockGlobalPriceInfoCache(ctrl),
				coinHistoryRepo:    mock_repository.NewMockCoinHistoryRepository(ctrl),
				gRequestPriceConfs: mock_cache.NewMockGlobalRequestPriceConfs(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &PriceService{
				gPriceInfosCache:   f.gPriceInfosCache,
				coinHistoryRepo:    f.coinHistoryRepo,
				gRequestPriceConfs: f.gRequestPriceConfs,
			}
			if got := s.GetLocalPrices(tt.args.start, tt.args.end, tt.args.symbol); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLocalPrices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceService_GetPartyPrice(t *testing.T) {
	type fields struct {
		gPriceInfosCache   *mock_cache.MockGlobalPriceInfoCache
		coinHistoryRepo    *mock_repository.MockCoinHistoryRepository
		gRequestPriceConfs *mock_cache.MockGlobalRequestPriceConfs
	}
	type args struct {
		symbol string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    bool
		want1   vo.PartyPriceInfo
	}{
		{
			name: "btc",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetLatestPriceInfos().Return(conf.PriceInfos{PriceInfos: priceInfos2})
			},
			args: args{symbol: "btcusdt"},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     47851.835,
				Timestamp: 1640745642,
				Infos: []vo.WeightInfo{
					{
						Price:        47852.37,
						Weight:       1,
						ExchangeName: "coinbase",
					},
					{
						Price:        47851.30,
						Weight:       1,
						ExchangeName: "ok",
					},
				},
			},
		},
		{
			name: "avax",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetLatestPriceInfos().Return(conf.PriceInfos{PriceInfos: priceInfos2})
			},
			args: args{symbol: "avaxusdt"},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     108.571,
				Timestamp: 1640745642,
				Infos: []vo.WeightInfo{
					{
						Price:        108.584,
						Weight:       1,
						ExchangeName: "huobi",
					},
					{
						Price:        108.579,
						Weight:       1,
						ExchangeName: "ok",
					},
					{
						Price:        108.55,
						Weight:       1,
						ExchangeName: "coinbase",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				gPriceInfosCache:   mock_cache.NewMockGlobalPriceInfoCache(ctrl),
				coinHistoryRepo:    mock_repository.NewMockCoinHistoryRepository(ctrl),
				gRequestPriceConfs: mock_cache.NewMockGlobalRequestPriceConfs(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &PriceService{
				gPriceInfosCache:   f.gPriceInfosCache,
				coinHistoryRepo:    f.coinHistoryRepo,
				gRequestPriceConfs: f.gRequestPriceConfs,
			}
			got, got1 := s.GetPartyPrice(tt.args.symbol)
			if got != tt.want {
				t.Errorf("GetPartyPrice() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetPartyPrice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPriceService_GetPrice(t *testing.T) {
	type fields struct {
		gPriceInfosCache   *mock_cache.MockGlobalPriceInfoCache
		coinHistoryRepo    *mock_repository.MockCoinHistoryRepository
		gRequestPriceConfs *mock_cache.MockGlobalRequestPriceConfs
	}
	type args struct {
		symbol   string
		exchange string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    bool
		want1   vo.PRICE_INFO
	}{
		{
			name: "btc ok",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetLatestPriceInfos().Return(conf.PriceInfos{PriceInfos: priceInfos2})
			},
			args: args{
				symbol:   "btcusdt",
				exchange: "ok",
			},
			want: true,
			want1: vo.PRICE_INFO{
				Price:     47851.3,
				Timestamp: 1640745642,
			},
		},
		{
			name: "btc huobi",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetLatestPriceInfos().Return(conf.PriceInfos{PriceInfos: priceInfos2})
			},
			args: args{
				symbol:   "btcusdt",
				exchange: "huobi",
			},
			want:  false,
			want1: vo.PRICE_INFO{},
		},
		{
			name: "avax ok",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetLatestPriceInfos().Return(conf.PriceInfos{PriceInfos: priceInfos2})
			},
			args: args{
				symbol:   "avaxusdt",
				exchange: "ok",
			},
			want: true,
			want1: vo.PRICE_INFO{
				Price:     108.57900000,
				Timestamp: 1640745642,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				gPriceInfosCache:   mock_cache.NewMockGlobalPriceInfoCache(ctrl),
				coinHistoryRepo:    mock_repository.NewMockCoinHistoryRepository(ctrl),
				gRequestPriceConfs: mock_cache.NewMockGlobalRequestPriceConfs(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &PriceService{
				gPriceInfosCache:   f.gPriceInfosCache,
				coinHistoryRepo:    f.coinHistoryRepo,
				gRequestPriceConfs: f.gRequestPriceConfs,
			}
			got, got1 := s.GetPrice(tt.args.symbol, tt.args.exchange)
			if got != tt.want {
				t.Errorf("GetPrice() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetPrice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPriceService_GetPriceAll(t *testing.T) {
	type fields struct {
		gPriceInfosCache   *mock_cache.MockGlobalPriceInfoCache
		coinHistoryRepo    *mock_repository.MockCoinHistoryRepository
		gRequestPriceConfs *mock_cache.MockGlobalRequestPriceConfs
	}
	type args struct {
		symbol string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    bool
		want1   []vo.PriceAllInfo
	}{
		{
			name: "btc",
			prepare: func(f *fields) {
				f.gPriceInfosCache.EXPECT().GetLatestPriceInfos().Return(conf.PriceInfos{PriceInfos: priceInfos2})
			},
			args: args{symbol: "btcusdt"},
			want: true,
			want1: []vo.PriceAllInfo{
				{
					Symbol:    "btcusdt",
					Price:     47851.30000000,
					Name:      "ok",
					Weight:    1,
					Timestamp: 1640745642,
				},
				{
					Symbol:    "btcusdt",
					Price:     47873.19000000,
					Name:      "bitstamp",
					Weight:    1,
					Timestamp: 1640745642,
				},
				{
					Symbol:    "btcusdt",
					Price:     47851.04000000,
					Name:      "binance",
					Weight:    1,
					Timestamp: 1640745642,
				},
				{
					Symbol:    "btcusdt",
					Price:     47852.37000000,
					Name:      "coinbase",
					Weight:    1,
					Timestamp: 1640745642,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				gPriceInfosCache:   mock_cache.NewMockGlobalPriceInfoCache(ctrl),
				coinHistoryRepo:    mock_repository.NewMockCoinHistoryRepository(ctrl),
				gRequestPriceConfs: mock_cache.NewMockGlobalRequestPriceConfs(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &PriceService{
				gPriceInfosCache:   f.gPriceInfosCache,
				coinHistoryRepo:    f.coinHistoryRepo,
				gRequestPriceConfs: f.gRequestPriceConfs,
			}
			got, got1 := s.GetPriceAll(tt.args.symbol)
			if got != tt.want {
				t.Errorf("GetPriceAll() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetPriceAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPriceService_InsertPriceInfo(t *testing.T) {
	type fields struct {
		gPriceInfosCache   *mock_cache.MockGlobalPriceInfoCache
		coinHistoryRepo    *mock_repository.MockCoinHistoryRepository
		gRequestPriceConfs *mock_cache.MockGlobalRequestPriceConfs
	}
	type args struct {
		cfg conf.PriceInfos
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr bool
	}{
		{
			name: "basic",
			prepare: func(f *fields) {
				f.coinHistoryRepo.EXPECT().InsertPriceInfo(gomock.Eq(conf.PriceInfos{PriceInfos: priceInfos2})).Return(nil)
			},
			args:    args{cfg: conf.PriceInfos{PriceInfos: priceInfos2}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				gPriceInfosCache:   mock_cache.NewMockGlobalPriceInfoCache(ctrl),
				coinHistoryRepo:    mock_repository.NewMockCoinHistoryRepository(ctrl),
				gRequestPriceConfs: mock_cache.NewMockGlobalRequestPriceConfs(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &PriceService{
				gPriceInfosCache:   f.gPriceInfosCache,
				coinHistoryRepo:    f.coinHistoryRepo,
				gRequestPriceConfs: f.gRequestPriceConfs,
			}
			if err := s.InsertPriceInfo(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("InsertPriceInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
