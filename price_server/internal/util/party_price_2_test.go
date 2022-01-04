package util

import (
	conf "price_api/price_server/config"
	"price_api/price_server/internal/vo"
	"reflect"
	"testing"
)

func TestPartyPrice2(t *testing.T) {
	type args struct {
		infos    []conf.PriceInfo
		bAverage bool
	}

	huobiPriceInfo := conf.PriceInfo{
		Symbol:      "bttusdt",
		Price:       0.00274154,
		PriceOrigin: "huobi",
		Weight:      2,
		TimeStamp:   1640681206,
	}

	kucoinPriceInfo := conf.PriceInfo{
		Symbol:      "bttusdt",
		Price:       0.00273997,
		PriceOrigin: "kucoin",
		Weight:      3,
		TimeStamp:   1640681206,
	}

	okPriceInfo := conf.PriceInfo{
		Symbol:      "bttusdt",
		Price:       0.0027403,
		PriceOrigin: "ok",
		Weight:      1,
		TimeStamp:   1640681206,
	}

	binancePriceInfo := conf.PriceInfo{
		Symbol:      "bttusdt",
		Price:       0.00274,
		PriceOrigin: "binance",
		Weight:      1,
		TimeStamp:   1640681206,
	}

	bitfinexPriceInfo := conf.PriceInfo{
		Symbol:      "bttusdt",
		Price:       0.0027392,
		PriceOrigin: "bitfinex",
		Weight:      1,
		TimeStamp:   1640681206,
	}

	tests := []struct {
		name  string
		args  args
		want  bool
		want1 vo.PartyPriceInfo
	}{
		{
			name: "infos length = 1 and average = true",
			args: args{
				infos: []conf.PriceInfo{
					bitfinexPriceInfo,
				},
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
					huobiPriceInfo,
					okPriceInfo,
				},
				bAverage: true,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     0.00274113,
				Timestamp: huobiPriceInfo.TimeStamp,
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
				},
			},
		},
		{
			name: "infos length = 3 and average = true",
			args: args{
				infos: []conf.PriceInfo{
					huobiPriceInfo,
					kucoinPriceInfo,
					okPriceInfo,
				},
				bAverage: true,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     0.0027403,
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
				},
				bAverage: true,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     0.00274004,
				Timestamp: okPriceInfo.TimeStamp,
				Infos: []vo.WeightInfo{
					{
						Price:        okPriceInfo.Price,
						Weight:       okPriceInfo.Weight,
						ExchangeName: okPriceInfo.PriceOrigin,
					},
					{
						Price:        binancePriceInfo.Price,
						Weight:       binancePriceInfo.Weight,
						ExchangeName: binancePriceInfo.PriceOrigin,
					},
					{
						Price:        kucoinPriceInfo.Price,
						Weight:       kucoinPriceInfo.Weight,
						ExchangeName: kucoinPriceInfo.PriceOrigin,
					},
				},
			},
		},
		{
			name: "infos length = 1 and average = false",
			args: args{
				infos: []conf.PriceInfo{
					bitfinexPriceInfo,
				},
				bAverage: false,
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
			name: "infos length = 2 and average = false",
			args: args{
				infos: []conf.PriceInfo{
					huobiPriceInfo,
					binancePriceInfo,
				},
				bAverage: false,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     0.00274103,
				Timestamp: huobiPriceInfo.TimeStamp,
				Infos: []vo.WeightInfo{
					{
						Price:        huobiPriceInfo.Price,
						Weight:       huobiPriceInfo.Weight,
						ExchangeName: huobiPriceInfo.PriceOrigin,
					},
					{
						Price:        binancePriceInfo.Price,
						Weight:       binancePriceInfo.Weight,
						ExchangeName: binancePriceInfo.PriceOrigin,
					},
				},
			},
		},
		{
			name: "infos length = 3 and average = false",
			args: args{
				infos: []conf.PriceInfo{
					huobiPriceInfo,
					binancePriceInfo,
					okPriceInfo,
				},
				bAverage: false,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     0.00274085,
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
						Price:        binancePriceInfo.Price,
						Weight:       binancePriceInfo.Weight,
						ExchangeName: binancePriceInfo.PriceOrigin,
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
				},
				bAverage: false,
			},
			want: true,
			want1: vo.PartyPriceInfo{
				Price:     0.00274031,
				Timestamp: bitfinexPriceInfo.TimeStamp,
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
						Price:        binancePriceInfo.Price,
						Weight:       binancePriceInfo.Weight,
						ExchangeName: binancePriceInfo.PriceOrigin,
					},
					{
						Price:        kucoinPriceInfo.Price,
						Weight:       kucoinPriceInfo.Weight,
						ExchangeName: kucoinPriceInfo.PriceOrigin,
					},
					{
						Price:        bitfinexPriceInfo.Price,
						Weight:       bitfinexPriceInfo.Weight,
						ExchangeName: bitfinexPriceInfo.PriceOrigin,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := PartyPrice(tt.args.infos, tt.args.bAverage)
			if got != tt.want {
				t.Errorf("PartyPrice() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PartyPrice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
