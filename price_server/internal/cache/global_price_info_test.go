package cache

import (
	conf "price_api/price_server/config"
	"reflect"
	"sync"
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
	priceInfos2 = []conf.PriceInfo{
		{
			Symbol:      "btcusdt",
			Price:       47929.80000000,
			PriceOrigin: "ok",
			Weight:      1,
			TimeStamp:   1640745981,
		},
		{
			Symbol:      "btcusdt",
			Price:       47930.75000000,
			PriceOrigin: "coinbase",
			Weight:      1,
			TimeStamp:   1640745981,
		},
		{
			Symbol:      "btcusdt",
			Price:       47873.19000000,
			PriceOrigin: "bitstamp",
			Weight:      1,
			TimeStamp:   1640745981,
		},
		{
			Symbol:      "btcusdt",
			Price:       47937.79000000,
			PriceOrigin: "huobi",
			Weight:      2,
			TimeStamp:   1640745981,
		},
		{
			Symbol:      "btcusdt",
			Price:       47927.56000000,
			PriceOrigin: "binance",
			Weight:      1,
			TimeStamp:   1640745981,
		},
		{
			Symbol:      "btcusdt",
			Price:       47934.20000000,
			PriceOrigin: "kucoin",
			Weight:      1,
			TimeStamp:   1640745981,
		},
		{
			Symbol:      "btcusdt",
			Price:       47950.00000000,
			PriceOrigin: "bitfinex",
			Weight:      1,
			TimeStamp:   1640745981,
		},
	}
	priceInfos3 = []conf.PriceInfo{
		{
			Symbol:      "btcusdt",
			Price:       47942.90000000,
			PriceOrigin: "kucoin",
			Weight:      1,
			TimeStamp:   1640746301,
		},
		{
			Symbol:      "btcusdt",
			Price:       47940.80000000,
			PriceOrigin: "ok",
			Weight:      1,
			TimeStamp:   1640746301,
		},
		{
			Symbol:      "btcusdt",
			Price:       47873.19000000,
			PriceOrigin: "bitstamp",
			Weight:      1,
			TimeStamp:   1640746301,
		},
		{
			Symbol:      "btcusdt",
			Price:       47948.04000000,
			PriceOrigin: "coinbase",
			Weight:      1,
			TimeStamp:   1640746301,
		},
		{
			Symbol:      "btcusdt",
			Price:       47948.19000000,
			PriceOrigin: "huobi",
			Weight:      2,
			TimeStamp:   1640746301,
		},
		{
			Symbol:      "btcusdt",
			Price:       47940.37000000,
			PriceOrigin: "binance",
			Weight:      1,
			TimeStamp:   1640746301,
		},
		{
			Symbol:      "btcusdt",
			Price:       47963.00000000,
			PriceOrigin: "bitfinex",
			Weight:      1,
			TimeStamp:   1640746301,
		},
	}
)

func generateTestPriceInfosCache() conf.PriceInfosCache {
	return conf.PriceInfosCache{
		PriceInfosCache: map[string][]conf.PriceInfos{
			"btcusdt": {
				{
					PriceInfos: priceInfos1,
				},
				{
					PriceInfos: priceInfos2,
				},
			},
		},
	}
}

func TestGlobalPriceInfoCache_GetLatestPriceInfos(t *testing.T) {
	type fields struct {
		gPriceInfosCache conf.PriceInfosCache
		m                *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   conf.PriceInfos
	}{
		{
			name: "basic",
			fields: fields{
				gPriceInfosCache: generateTestPriceInfosCache(),
				m:                new(sync.RWMutex),
			},
			want: conf.PriceInfos{PriceInfos: priceInfos2},
		},
		{
			name: "empty",
			fields: fields{
				gPriceInfosCache: conf.PriceInfosCache{},
				m:                new(sync.RWMutex),
			},
			want: conf.PriceInfos{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &globalPriceInfoCache{
				gPriceInfosCache: tt.fields.gPriceInfosCache,
				m:                tt.fields.m,
			}
			if got := c.GetLatestPriceInfos("btcusdt"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLatestPriceInfos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGlobalPriceInfoCache_GetPriceInfosEqualTimestamp(t *testing.T) {
	type fields struct {
		gPriceInfosCache conf.PriceInfosCache
		m                *sync.RWMutex
	}
	type args struct {
		timestamp int64
		symbol    string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  conf.PriceInfos
	}{
		{
			name: "found at index 0",
			fields: fields{
				gPriceInfosCache: generateTestPriceInfosCache(),
				m:                new(sync.RWMutex),
			},
			args:  args{timestamp: 1640745642, symbol: "btcusdt"},
			want:  true,
			want1: conf.PriceInfos{PriceInfos: priceInfos1},
		},
		{
			name: "found at index 1",
			fields: fields{
				gPriceInfosCache: generateTestPriceInfosCache(),
				m:                new(sync.RWMutex),
			},
			args:  args{timestamp: 1640745981, symbol: "btcusdt"},
			want:  true,
			want1: conf.PriceInfos{PriceInfos: priceInfos2},
		},
		{
			name: "skip when priceInfos length is 0",
			fields: fields{
				gPriceInfosCache: conf.PriceInfosCache{
					PriceInfosCache: map[string][]conf.PriceInfos{
						"btcusdt": {
							{
								PriceInfos: priceInfos1,
							},
							{
								PriceInfos: []conf.PriceInfo{},
							},
							{
								PriceInfos: priceInfos2,
							},
						},
					},
				},
				m: new(sync.RWMutex),
			},
			args:  args{timestamp: 1640745981, symbol: "btcusdt"},
			want:  true,
			want1: conf.PriceInfos{PriceInfos: priceInfos2},
		},

		{
			name: "not found",
			fields: fields{
				gPriceInfosCache: generateTestPriceInfosCache(),
				m:                new(sync.RWMutex),
			},
			args:  args{timestamp: 1640745641},
			want:  false,
			want1: conf.PriceInfos{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &globalPriceInfoCache{
				gPriceInfosCache: tt.fields.gPriceInfosCache,
				m:                tt.fields.m,
			}
			got, got1 := c.GetPriceInfosEqualTimestamp(tt.args.symbol, tt.args.timestamp)
			if got != tt.want {
				t.Errorf("GetPriceInfosEqualTimestamp() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetPriceInfosEqualTimestamp() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGlobalPriceInfoCache_GetPriceInfosByRange(t *testing.T) {
	type fields struct {
		gPriceInfosCache conf.PriceInfosCache
		m                *sync.RWMutex
	}
	type args struct {
		symbol string
		start  int
		end    int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   conf.PriceInfosCache
	}{
		{
			name: "find start 0 ,to 1",
			fields: fields{
				gPriceInfosCache: generateTestPriceInfosCache(),
				m:                new(sync.RWMutex),
			},
			args: args{
				symbol: "btcusdt",
				start:  0,
				end:    1,
			},
			want: conf.PriceInfosCache{PriceInfosCache: map[string][]conf.PriceInfos{
				"btcusdt": {
					{
						PriceInfos: priceInfos1,
					},
				},
			}},
		},
		{
			name: "find start 0 ,to 2",
			fields: fields{
				gPriceInfosCache: generateTestPriceInfosCache(),
				m:                new(sync.RWMutex),
			},
			args: args{
				symbol: "btcusdt",
				start:  0,
				end:    2,
			},
			want: conf.PriceInfosCache{PriceInfosCache: map[string][]conf.PriceInfos{
				"btcusdt": {
					{
						PriceInfos: priceInfos1,
					},
					{
						PriceInfos: priceInfos2,
					},
				},
			}},
		},
		{
			name: "find start 0 ,to 3",
			fields: fields{
				gPriceInfosCache: generateTestPriceInfosCache(),
				m:                new(sync.RWMutex),
			},
			args: args{
				symbol: "btcusdt",
				start:  0,
				end:    3,
			},
			want: conf.PriceInfosCache{PriceInfosCache: map[string][]conf.PriceInfos{
				"btcusdt": {
					{
						PriceInfos: priceInfos1,
					},
					{
						PriceInfos: priceInfos2,
					},
				},
			},
			},
		},
		{
			name: "find start 1 ,to 3",
			fields: fields{
				gPriceInfosCache: generateTestPriceInfosCache(),
				m:                new(sync.RWMutex),
			},
			args: args{
				symbol: "btcusdt",
				start:  1,
				end:    3,
			},
			want: conf.PriceInfosCache{PriceInfosCache: map[string][]conf.PriceInfos{
				"btcusdt": {
					{
						PriceInfos: priceInfos2,
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &globalPriceInfoCache{
				gPriceInfosCache: tt.fields.gPriceInfosCache,
				m:                tt.fields.m,
			}
			if got := c.GetPriceInfosByRange(tt.args.symbol, tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPriceInfosByRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGlobalPriceInfoCache_GetCacheLength(t *testing.T) {
	type fields struct {
		gPriceInfosCache conf.PriceInfosCache
		m                *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "basic",
			fields: fields{
				gPriceInfosCache: generateTestPriceInfosCache(),
				m:                new(sync.RWMutex),
			},
			want: 1,
		},
		{
			name: "empty",
			fields: fields{
				gPriceInfosCache: conf.PriceInfosCache{},
				m:                new(sync.RWMutex),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &globalPriceInfoCache{
				gPriceInfosCache: tt.fields.gPriceInfosCache,
				m:                tt.fields.m,
			}
			if got := c.GetCacheLength(); got != tt.want {
				t.Errorf("GetCacheLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGlobalPriceInfoCache_UpdateCachePrice(t *testing.T) {
	type fields struct {
		gPriceInfosCache conf.PriceInfosCache
		m                *sync.RWMutex
	}
	type args struct {
		symbol     string
		infos      conf.PriceInfos
		maxMemTime int
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantLength    int
		wantPriceInfo conf.PriceInfos
	}{
		{
			name: "basic",
			fields: fields{
				gPriceInfosCache: generateTestPriceInfosCache(),
				m:                new(sync.RWMutex),
			},
			args: args{
				symbol:     "btcusdt",
				infos:      conf.PriceInfos{priceInfos3},
				maxMemTime: 5,
			},
			wantLength:    3,
			wantPriceInfo: conf.PriceInfos{priceInfos3},
		},
		{
			name: "over maxMemTime",
			fields: fields{
				gPriceInfosCache: generateTestPriceInfosCache(),
				m:                new(sync.RWMutex),
			},
			args: args{
				symbol:     "btcusdt",
				infos:      conf.PriceInfos{priceInfos3},
				maxMemTime: 2,
			},
			wantLength:    2,
			wantPriceInfo: conf.PriceInfos{priceInfos3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &globalPriceInfoCache{
				gPriceInfosCache: tt.fields.gPriceInfosCache,
				m:                tt.fields.m,
			}
			c.UpdateCachePrice(tt.args.symbol, tt.args.infos, tt.args.maxMemTime)

			gotLatestPriceInfos := c.GetLatestPriceInfos(tt.args.symbol)
			gotLength := c.GetSymbolCacheLength(tt.args.symbol)

			if gotLength != tt.wantLength {
				t.Errorf("UpdateCachePrice() gotLength = %v, wantLength %v", gotLength, tt.wantLength)
			}
			if !reflect.DeepEqual(gotLatestPriceInfos, tt.wantPriceInfo) {
				t.Errorf("UpdateCachePrice() gotLatestPriceInfos = %v, wantPriceInfo %v", gotLatestPriceInfos, tt.wantPriceInfo)
			}

		})
	}
}
