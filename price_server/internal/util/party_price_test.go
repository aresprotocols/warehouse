package util

import (
	conf "price_api/price_server/config"
	"price_api/price_server/internal/vo"
	"reflect"
	"testing"
)

func TestPartyPrice(t *testing.T) {
	type args struct {
		infos    []conf.PriceInfo
		symbol   string
		bAverage bool
	}

	bitfinexPriceInfo := conf.PriceInfo{
		Symbol:      "btcusdt",
		Price:       50814,
		PriceOrigin: "bitfinex",
		Weight:      1,
		TimeStamp:   1640592554,
	}

	kucoinPriceInfo := conf.PriceInfo{
		Symbol:      "btcusdt",
		Price:       50789,
		PriceOrigin: "kucoin",
		Weight:      1,
		TimeStamp:   1640592554,
	}
	huobiPriceInfo := conf.PriceInfo{
		Symbol:      "btcusdt",
		Price:       50787.31,
		PriceOrigin: "huobi",
		Weight:      2,
		TimeStamp:   1640592554,
	}

	binancePriceInfo := conf.PriceInfo{
		Symbol:      "btcusdt",
		Price:       50787.44,
		PriceOrigin: "binance",
		Weight:      1,
		TimeStamp:   1640592554,
	}

	okPriceInfo := conf.PriceInfo{
		Symbol:      "btcusdt",
		Price:       50785,
		PriceOrigin: "ok",
		Weight:      1,
		TimeStamp:   1640592554,
	}
	bitstampPriceInfo := conf.PriceInfo{
		Symbol:      "btcusdt",
		Price:       50734.34,
		PriceOrigin: "bitstamp",
		Weight:      1,
		TimeStamp:   1640592554,
	}
	coinbasePriceInfo := conf.PriceInfo{
		Symbol:      "btcusdt",
		Price:       50781.23,
		PriceOrigin: "coinbase",
		Weight:      3,
		TimeStamp:   1640592554,
	}

	tests := []struct {
		name  string
		args  args
		want  bool
		want1 vo.PartyPriceInfo
	}{
		{
			name: "infos length is empty",
			args: args{
				infos:    []conf.PriceInfo{},
				symbol:   "btcusdt",
				bAverage: false,
			},
			want:  false,
			want1: vo.PartyPriceInfo{},
		},
		{
			name: "infos not contain symbol",
			args: args{
				infos: []conf.PriceInfo{
					bitfinexPriceInfo,
					kucoinPriceInfo,
					huobiPriceInfo,
					binancePriceInfo,
					okPriceInfo,
					bitstampPriceInfo,
					coinbasePriceInfo,
				},
				symbol:   "ethusdt",
				bAverage: true,
			},
			want:  false,
			want1: vo.PartyPriceInfo{},
		},
		{
			name: "infos length = 1 and average = true",
			args: args{
				infos: []conf.PriceInfo{
					bitfinexPriceInfo,
				},
				symbol:   "btcusdt",
				bAverage: true,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     bitfinexPriceInfo.Price,
				Timestamp: bitfinexPriceInfo.TimeStamp,
				Infos: []vo.WeightInfo{
					{
						Price:        bitfinexPriceInfo.Price,
						Weight:       bitfinexPriceInfo.Weight,
						ExchangeName: bitfinexPriceInfo.PriceOrigin,
					},
				},
			},
		},
		{
			name: "infos length = 2 and average = true",
			args: args{
				infos: []conf.PriceInfo{
					coinbasePriceInfo,
					okPriceInfo,
				},
				symbol:   "btcusdt",
				bAverage: true,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     50782.1725,
				Timestamp: coinbasePriceInfo.TimeStamp,
				Infos: []vo.WeightInfo{
					{
						Price:        okPriceInfo.Price,
						Weight:       okPriceInfo.Weight,
						ExchangeName: okPriceInfo.PriceOrigin,
					},
					{
						Price:        coinbasePriceInfo.Price,
						Weight:       coinbasePriceInfo.Weight,
						ExchangeName: coinbasePriceInfo.PriceOrigin,
					},
				},
			},
		},
		{
			name: "infos length = 3 and average = true",
			args: args{
				infos: []conf.PriceInfo{
					huobiPriceInfo,
					okPriceInfo,
					coinbasePriceInfo,
				},
				symbol:   "btcusdt",
				bAverage: true,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     50785,
				Timestamp: okPriceInfo.TimeStamp,
				Infos: []vo.WeightInfo{
					{
						Price:        okPriceInfo.Price,
						Weight:       okPriceInfo.Weight,
						ExchangeName: okPriceInfo.PriceOrigin,
					},
				},
			},
		},
		{
			name: "infos length over 3 and average = true",
			args: args{
				infos: []conf.PriceInfo{
					bitfinexPriceInfo,
					kucoinPriceInfo,
					huobiPriceInfo,
					binancePriceInfo,
					okPriceInfo,
					bitstampPriceInfo,
					coinbasePriceInfo,
				},
				symbol:   "btcusdt",
				bAverage: true,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     50784.96875,
				Timestamp: 1640592554,
				Infos: []vo.WeightInfo{
					{
						Price:        kucoinPriceInfo.Price,
						Weight:       kucoinPriceInfo.Weight,
						ExchangeName: kucoinPriceInfo.PriceOrigin,
					},
					{
						Price:        binancePriceInfo.Price,
						Weight:       binancePriceInfo.Weight,
						ExchangeName: binancePriceInfo.PriceOrigin,
					},
					{
						Price:        huobiPriceInfo.Price,
						Weight:       huobiPriceInfo.Weight,
						ExchangeName: huobiPriceInfo.PriceOrigin,
					},
					{
						Price:        okPriceInfo.Price,
						Weight:       okPriceInfo.Weight,
						ExchangeName: okPriceInfo.PriceOrigin,
					},
					{
						Price:        coinbasePriceInfo.Price,
						Weight:       coinbasePriceInfo.Weight,
						ExchangeName: coinbasePriceInfo.PriceOrigin,
					},
				},
			},
		},
		{
			name: "infos length = 1 and average = false",
			args: args{
				infos: []conf.PriceInfo{
					bitstampPriceInfo,
				},
				symbol:   "btcusdt",
				bAverage: false,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     bitstampPriceInfo.Price,
				Timestamp: bitstampPriceInfo.TimeStamp,
				Infos: []vo.WeightInfo{
					{
						Price:        bitstampPriceInfo.Price,
						Weight:       bitstampPriceInfo.Weight,
						ExchangeName: bitstampPriceInfo.PriceOrigin,
					},
				},
			},
		},
		{
			name: "infos length = 2 and average = false",
			args: args{
				infos: []conf.PriceInfo{
					huobiPriceInfo,
					binancePriceInfo,
				},
				symbol:   "btcusdt",
				bAverage: false,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     50787.35333333333,
				Timestamp: huobiPriceInfo.TimeStamp,
				Infos: []vo.WeightInfo{
					{
						Price:        binancePriceInfo.Price,
						Weight:       binancePriceInfo.Weight,
						ExchangeName: binancePriceInfo.PriceOrigin,
					},
					{
						Price:        huobiPriceInfo.Price,
						Weight:       huobiPriceInfo.Weight,
						ExchangeName: huobiPriceInfo.PriceOrigin,
					},
				},
			},
		},
		{
			name: "infos length = 3 and average = false",
			args: args{
				infos: []conf.PriceInfo{
					huobiPriceInfo,
					coinbasePriceInfo,
					okPriceInfo,
				},
				symbol:   "btcusdt",
				bAverage: false,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     50783.885,
				Timestamp: okPriceInfo.TimeStamp,
				Infos: []vo.WeightInfo{
					{
						Price:        huobiPriceInfo.Price,
						Weight:       huobiPriceInfo.Weight,
						ExchangeName: huobiPriceInfo.PriceOrigin,
					},
					{
						Price:        okPriceInfo.Price,
						Weight:       okPriceInfo.Weight,
						ExchangeName: okPriceInfo.PriceOrigin,
					},

					{
						Price:        coinbasePriceInfo.Price,
						Weight:       coinbasePriceInfo.Weight,
						ExchangeName: coinbasePriceInfo.PriceOrigin,
					},
				},
			},
		},
		{
			name: "infos length > 3 and average = false",
			args: args{
				infos: []conf.PriceInfo{
					bitfinexPriceInfo,
					kucoinPriceInfo,
					huobiPriceInfo,
					binancePriceInfo,
					okPriceInfo,
					bitstampPriceInfo,
					coinbasePriceInfo,
				},
				symbol:   "btcusdt",
				bAverage: false,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     50782.809,
				Timestamp: bitfinexPriceInfo.TimeStamp,
				Infos: []vo.WeightInfo{
					{
						Price:        bitfinexPriceInfo.Price,
						Weight:       bitfinexPriceInfo.Weight,
						ExchangeName: bitfinexPriceInfo.PriceOrigin,
					},
					{
						Price:        kucoinPriceInfo.Price,
						Weight:       kucoinPriceInfo.Weight,
						ExchangeName: kucoinPriceInfo.PriceOrigin,
					},
					{
						Price:        binancePriceInfo.Price,
						Weight:       binancePriceInfo.Weight,
						ExchangeName: binancePriceInfo.PriceOrigin,
					},
					{
						Price:        huobiPriceInfo.Price,
						Weight:       huobiPriceInfo.Weight,
						ExchangeName: huobiPriceInfo.PriceOrigin,
					},
					{
						Price:        okPriceInfo.Price,
						Weight:       okPriceInfo.Weight,
						ExchangeName: okPriceInfo.PriceOrigin,
					},
					{
						Price:        coinbasePriceInfo.Price,
						Weight:       coinbasePriceInfo.Weight,
						ExchangeName: coinbasePriceInfo.PriceOrigin,
					},
					{
						Price:        bitstampPriceInfo.Price,
						Weight:       bitstampPriceInfo.Weight,
						ExchangeName: bitstampPriceInfo.PriceOrigin,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := PartyPrice(tt.args.infos, tt.args.symbol, tt.args.bAverage)
			if got != tt.want {
				t.Errorf("PartyPrice() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PartyPrice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
