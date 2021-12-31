package service

import (
	"github.com/golang/mock/gomock"
	mock_repository "price_api/price_server/internal/repository/mock"
	"price_api/price_server/internal/vo"
	"reflect"
	"testing"
)

var (
	logInfos = vo.LOG_INFOS{
		Infos: []vo.LOG_INFO{
			{
				ClientIP:     "127.0.0.1",
				RequestTime:  "2021-12-31 09:42:27",
				UserAgent:    "insomnia/2021.7.2",
				RequestUrl:   "/api/getPriceAll/btcusdt",
				ResponseTime: "2021-12-31 09:42:27",
				Response:     "{\"code\":0,\"message\":\"OK\",\"data\":[{\"name\":\"ok\",\"symbol\":\"btcusdt\",\"price\":47145.7,\"timestamp\":1640914944,\"weight\":1},{\"name\":\"bitfinex\",\"symbol\":\"btcusdt\",\"price\":47161,\"timestamp\":1640914944,\"weight\":1},{\"name\":\"binance\",\"symbol\":\"btcusdt\",\"price\":47150.9,\"timestamp\":1640914944,\"weight\":1},{\"name\":\"huobi\",\"symbol\":\"btcusdt\",\"price\":47151.82,\"timestamp\":1640914944,\"weight\":2},{\"name\":\"coinbase\",\"symbol\":\"btcusdt\",\"price\":47143.64,\"timestamp\":1640914944,\"weight\":1},{\"name\":\"bitstamp\",\"symbol\":\"btcusdt\",\"price\":47198.07,\"timestamp\":1640914944,\"weight\":1},{\"name\":\"kucoin\",\"symbol\":\"btcusdt\",\"price\":47154.9,\"timestamp\":1640914944,\"weight\":1}]}",
			},
			{
				ClientIP:     "127.0.0.1",
				RequestTime:  "2021-12-31 09:42:27",
				UserAgent:    "insomnia/2021.7.2",
				RequestUrl:   "/api/getPrice/btcusdt/huobi",
				ResponseTime: "2021-12-31 09:42:27",
				Response:     "{\"code\":0,\"message\":\"OK\",\"data\":{\"price\":47265.09,\"timestamp\":1640875437}}",
			},
			{
				ClientIP:     "127.0.0.1",
				RequestTime:  "2021-12-30 21:39:41",
				UserAgent:    "insomnia/2021.7.2",
				RequestUrl:   "/api/getPartyPrice/btcusdt",
				ResponseTime: "2021-12-30 21:39:41",
				Response:     "{\"code\":0,\"message\":\"OK\",\"data\":{\"price\":47497.29,\"timestamp\":1640871542,\"infos\":[{\"price\":47502.48,\"weight\":1,\"exchangeName\":\"coinbase\"},{\"price\":47501,\"weight\":1,\"exchangeName\":\"bitfinex\"},{\"price\":47495.43,\"weight\":2,\"exchangeName\":\"huobi\"},{\"price\":47494.7,\"weight\":1,\"exchangeName\":\"ok\"},{\"price\":47494.7,\"weight\":1,\"exchangeName\":\"kucoin\"}]}}",
			},
			{
				ClientIP:     "127.0.0.1",
				RequestTime:  "2021-12-31 09:43:00",
				UserAgent:    "insomnia/2021.7.2",
				RequestUrl:   "/api/getHistoryPrice/btcusdt?timestamp=1639640575",
				ResponseTime: "2021-12-31 09:43:00",
				Response:     "{\"code\":0,\"message\":\"OK\",\"data\":{\"price\":48898.17,\"timestamp\":1639640575,\"infos\":[{\"price\":48940,\"weight\":1,\"exchangeName\":\"bitfinex\"},{\"price\":48890.1,\"weight\":1,\"exchangeName\":\"ok\"},{\"price\":48881.78,\"weight\":1,\"exchangeName\":\"binance\"},{\"price\":48880.8,\"weight\":1,\"exchangeName\":\"huobi\"}]}}",
			},
			{
				ClientIP:     "127.0.0.1",
				RequestTime:  "2021-12-23 11:41:35",
				UserAgent:    "insomnia/2021.7.2",
				RequestUrl:   "/api/getBulkPrices?symbol=btcusdt",
				ResponseTime: "2021-12-23 11:41:35",
				Response:     "{\"code\":0,\"message\":\"OK\",\"data\":{\"btcusdt\":{\"price\":48415.346000000005,\"timestamp\":1640230878}}}",
			},
			{
				ClientIP:     "127.0.0.1",
				RequestTime:  "2021-12-30 22:44:53",
				UserAgent:    "insomnia/2021.7.2",
				RequestUrl:   "/api/getBulkCurrencyPrices?symbol=btc_eth_avax_ltc_bch_fil_etc_eos_dash_comp_matic&currency=usdt",
				ResponseTime: "2021-12-30 22:44:53",
				Response:     "{\"code\":0,\"message\":\"OK\",\"data\":{\"avaxusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"bchusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"btcusdt\":{\"price\":47253.101667,\"timestamp\":1640875484,\"infos\":[{\"price\":47258.08,\"weight\":2,\"exchangeName\":\"huobi\"},{\"price\":47257.7,\"weight\":1,\"exchangeName\":\"ok\"},{\"price\":47254.97,\"weight\":1,\"exchangeName\":\"binance\"},{\"price\":47252.78,\"weight\":1,\"exchangeName\":\"coinbase\"},{\"price\":47237,\"weight\":1,\"exchangeName\":\"bitfinex\"}]},\"compusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"dashusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"eosusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"etcusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"ethusdt\":{\"price\":3706.6,\"timestamp\":1640875483,\"infos\":[{\"price\":3707.29,\"weight\":1,\"exchangeName\":\"huobi\"},{\"price\":3707.15,\"weight\":1,\"exchangeName\":\"kucoin\"},{\"price\":3706.59,\"weight\":1,\"exchangeName\":\"binance\"},{\"price\":3706.07,\"weight\":1,\"exchangeName\":\"coinbase\"},{\"price\":3705.9,\"weight\":1,\"exchangeName\":\"ok\"}]},\"filusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"ltcusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"maticusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null}}}",
			},
		},
	}

	getPriceAllReqRspLogInfo = vo.REQ_RSP_LOG_INFO{
		Ip:               "127.0.0.1",
		RequestTime:      "2021-12-31 09:42:27",
		ReqUrl:           "/api/getPriceAll/btcusdt",
		Response:         "{\"code\":0,\"message\":\"OK\",\"data\":[{\"name\":\"ok\",\"symbol\":\"btcusdt\",\"price\":47145.7,\"timestamp\":1640914944,\"weight\":1},{\"name\":\"bitfinex\",\"symbol\":\"btcusdt\",\"price\":47161,\"timestamp\":1640914944,\"weight\":1},{\"name\":\"binance\",\"symbol\":\"btcusdt\",\"price\":47150.9,\"timestamp\":1640914944,\"weight\":1},{\"name\":\"huobi\",\"symbol\":\"btcusdt\",\"price\":47151.82,\"timestamp\":1640914944,\"weight\":2},{\"name\":\"coinbase\",\"symbol\":\"btcusdt\",\"price\":47143.64,\"timestamp\":1640914944,\"weight\":1},{\"name\":\"bitstamp\",\"symbol\":\"btcusdt\",\"price\":47198.07,\"timestamp\":1640914944,\"weight\":1},{\"name\":\"kucoin\",\"symbol\":\"btcusdt\",\"price\":47154.9,\"timestamp\":1640914944,\"weight\":1}]}",
		RequestTimestamp: 1640914947,
	}

	getPriceReqRspLogInfo = vo.REQ_RSP_LOG_INFO{
		Ip:               "127.0.0.1",
		RequestTime:      "2021-12-31 09:42:27",
		ReqUrl:           "/api/getPrice/btcusdt/huobi",
		RequestTimestamp: 1640914947,
		Response:         "{\"code\":0,\"message\":\"OK\",\"data\":{\"price\":47265.09,\"timestamp\":1640875437}}",
	}

	getPartyPriceReqRspLogInfo = vo.REQ_RSP_LOG_INFO{
		Ip:               "127.0.0.1",
		RequestTime:      "2021-12-30 21:39:41",
		ReqUrl:           "/api/getPartyPrice/btcusdt",
		RequestTimestamp: 1640914947,
		Response:         "{\"code\":0,\"message\":\"OK\",\"data\":{\"price\":47497.29,\"timestamp\":1640871542,\"infos\":[{\"price\":47502.48,\"weight\":1,\"exchangeName\":\"coinbase\"},{\"price\":47501,\"weight\":1,\"exchangeName\":\"bitfinex\"},{\"price\":47495.43,\"weight\":2,\"exchangeName\":\"huobi\"},{\"price\":47494.7,\"weight\":1,\"exchangeName\":\"ok\"},{\"price\":47494.7,\"weight\":1,\"exchangeName\":\"kucoin\"}]}}",
	}

	getHistoryPriceReqRspLogInfo = vo.REQ_RSP_LOG_INFO{
		Ip:               "127.0.0.1",
		RequestTime:      "2021-12-31 09:43:00",
		ReqUrl:           "/api/getHistoryPrice/btcusdt?timestamp=1639640575",
		RequestTimestamp: 1640914947,
		Response:         "{\"code\":0,\"message\":\"OK\",\"data\":{\"price\":48898.17,\"timestamp\":1639640575,\"infos\":[{\"price\":48940,\"weight\":1,\"exchangeName\":\"bitfinex\"},{\"price\":48890.1,\"weight\":1,\"exchangeName\":\"ok\"},{\"price\":48881.78,\"weight\":1,\"exchangeName\":\"binance\"},{\"price\":48880.8,\"weight\":1,\"exchangeName\":\"huobi\"}]}}",
	}

	getBulkPricesReqRspLogInfo = vo.REQ_RSP_LOG_INFO{
		Ip:               "127.0.0.1",
		RequestTime:      "2021-12-23 11:41:35",
		ReqUrl:           "/api/getBulkPrices?symbol=btcusdt",
		RequestTimestamp: 1640914947,
		Response:         "{\"code\":0,\"message\":\"OK\",\"data\":{\"btcusdt\":{\"price\":48415.346000000005,\"timestamp\":1640230878}}}",
	}
	getBulkCurrencyPricesReqRspLogInfo = vo.REQ_RSP_LOG_INFO{
		Ip:               "127.0.0.1",
		RequestTime:      "2021-12-30 22:44:53",
		ReqUrl:           "/api/getBulkCurrencyPrices?symbol=btc_eth_avax_ltc_bch_fil_etc_eos_dash_comp_matic&currency=usdt",
		RequestTimestamp: 1640914947,
		Response:         "{\"code\":0,\"message\":\"OK\",\"data\":{\"avaxusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"bchusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"btcusdt\":{\"price\":47253.101667,\"timestamp\":1640875484,\"infos\":[{\"price\":47258.08,\"weight\":2,\"exchangeName\":\"huobi\"},{\"price\":47257.7,\"weight\":1,\"exchangeName\":\"ok\"},{\"price\":47254.97,\"weight\":1,\"exchangeName\":\"binance\"},{\"price\":47252.78,\"weight\":1,\"exchangeName\":\"coinbase\"},{\"price\":47237,\"weight\":1,\"exchangeName\":\"bitfinex\"}]},\"compusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"dashusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"eosusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"etcusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"ethusdt\":{\"price\":3706.6,\"timestamp\":1640875483,\"infos\":[{\"price\":3707.29,\"weight\":1,\"exchangeName\":\"huobi\"},{\"price\":3707.15,\"weight\":1,\"exchangeName\":\"kucoin\"},{\"price\":3706.59,\"weight\":1,\"exchangeName\":\"binance\"},{\"price\":3706.07,\"weight\":1,\"exchangeName\":\"coinbase\"},{\"price\":3705.9,\"weight\":1,\"exchangeName\":\"ok\"}]},\"filusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"ltcusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null},\"maticusdt\":{\"price\":0,\"timestamp\":0,\"infos\":null}}}",
	}
	reqLogInfos = []vo.REQ_RSP_LOG_INFO{
		getPriceAllReqRspLogInfo,
		getPriceReqRspLogInfo,
		getPartyPriceReqRspLogInfo,
		getHistoryPriceReqRspLogInfo,
		getBulkPricesReqRspLogInfo,
		getBulkCurrencyPricesReqRspLogInfo,
	}

	getPriceAllParseLog = vo.PARTY_PRICE_INFO{
		Type: "getPriceAll",
		Client: vo.CLIENT_INFO{
			Ip:               "127.0.0.1",
			RequestTime:      "2021-12-31 09:42:27",
			RequestTimestamp: 1640914947,
		},
		PriceInfo: vo.PRICE_INFO{
			Price:     47145.7,
			Timestamp: 1640914944,
		},
		PriceInfos: []vo.PRICE_EXCHANGE_WEIGHT_INFO{
			{
				Exchange:  "ok",
				Price:     47145.7,
				Timestamp: 1640914944,
				Weight:    1,
			}, {
				Exchange:  "bitfinex",
				Price:     47161,
				Timestamp: 1640914944,
				Weight:    1,
			}, {
				Exchange:  "binance",
				Price:     47150.9,
				Timestamp: 1640914944,
				Weight:    1,
			}, {
				Exchange:  "huobi",
				Price:     47151.82,
				Timestamp: 1640914944,
				Weight:    2,
			}, {
				Exchange:  "coinbase",
				Price:     47143.64,
				Timestamp: 1640914944,
				Weight:    1,
			}, {
				Exchange:  "bitstamp",
				Price:     47198.07,
				Timestamp: 1640914944,
				Weight:    1,
			}, {
				Exchange:  "kucoin",
				Price:     47154.9,
				Timestamp: 1640914944,
				Weight:    1,
			},
		},
	}
	getPriceParseLog = vo.PARTY_PRICE_INFO{
		Type: "getPrice",
		Client: vo.CLIENT_INFO{
			Ip:               "127.0.0.1",
			RequestTime:      "2021-12-31 09:42:27",
			RequestTimestamp: 1640914947,
		},
		PriceInfo: vo.PRICE_INFO{
			Price:     47265.09,
			Timestamp: 1640875437,
		},
		PriceInfos: []vo.PRICE_EXCHANGE_WEIGHT_INFO{},
	}
	getPartyPriceParseLog = vo.PARTY_PRICE_INFO{
		Type: "getPartyPrice",
		Client: vo.CLIENT_INFO{
			Ip:               "127.0.0.1",
			RequestTime:      "2021-12-30 21:39:41",
			RequestTimestamp: 1640914947,
		},
		PriceInfo: vo.PRICE_INFO{
			Price:     47497.29,
			Timestamp: 1640871542,
		},
		PriceInfos: []vo.PRICE_EXCHANGE_WEIGHT_INFO{
			{
				Price:     47502.48,
				Weight:    1,
				Exchange:  "coinbase",
				Timestamp: 1640871542,
			}, {
				Price:     47501,
				Weight:    1,
				Exchange:  "bitfinex",
				Timestamp: 1640871542,
			}, {
				Price:     47495.43,
				Weight:    2,
				Exchange:  "huobi",
				Timestamp: 1640871542,
			}, {
				Price:     47494.7,
				Weight:    1,
				Exchange:  "ok",
				Timestamp: 1640871542,
			}, {
				Price:     47494.7,
				Weight:    1,
				Exchange:  "kucoin",
				Timestamp: 1640871542,
			},
		},
	}
	getHistoryPriceParseLog = vo.PARTY_PRICE_INFO{
		Type: "getHistoryPrice",
		Client: vo.CLIENT_INFO{
			Ip:               "127.0.0.1",
			RequestTime:      "2021-12-31 09:43:00",
			RequestTimestamp: 1640914947,
		},
		PriceInfo: vo.PRICE_INFO{
			Price:     48898.17,
			Timestamp: 1639640575,
		},
		PriceInfos: []vo.PRICE_EXCHANGE_WEIGHT_INFO{
			{
				Price:     48940,
				Weight:    1,
				Exchange:  "bitfinex",
				Timestamp: 1639640575,
			}, {
				Price:     48890.1,
				Weight:    1,
				Exchange:  "ok",
				Timestamp: 1639640575,
			}, {
				Price:     48881.78,
				Weight:    1,
				Exchange:  "binance",
				Timestamp: 1639640575,
			}, {
				Price:     48880.8,
				Weight:    1,
				Exchange:  "huobi",
				Timestamp: 1639640575,
			},
		},
	}
	getBulkPricesParseLog = vo.PARTY_PRICE_INFO{
		Type: "getBulkPrices",
		Client: vo.CLIENT_INFO{
			Ip:               "127.0.0.1",
			RequestTime:      "2021-12-23 11:41:35",
			RequestTimestamp: 1640914947,
		},
		PriceInfo: vo.PRICE_INFO{
			Price:     48415.346000000005,
			Timestamp: 1640230878,
		},
		PriceInfos: []vo.PRICE_EXCHANGE_WEIGHT_INFO{},
	}
	getBulkCurrencyPricesParseLog = vo.PARTY_PRICE_INFO{
		Type: "getBulkCurrencyPrices",
		Client: vo.CLIENT_INFO{
			Ip:               "127.0.0.1",
			RequestTime:      "2021-12-30 22:44:53",
			RequestTimestamp: 1640914947,
		},
		PriceInfo: vo.PRICE_INFO{
			Price:     47253.101667,
			Timestamp: 1640875484,
		},
		PriceInfos: []vo.PRICE_EXCHANGE_WEIGHT_INFO{
			{
				Price:     47258.08,
				Weight:    2,
				Exchange:  "huobi",
				Timestamp: 1640875484,
			}, {
				Price:     47257.7,
				Weight:    1,
				Exchange:  "ok",
				Timestamp: 1640875484,
			}, {
				Price:     47254.97,
				Weight:    1,
				Exchange:  "binance",
				Timestamp: 1640875484,
			}, {
				Price:     47252.78,
				Weight:    1,
				Exchange:  "coinbase",
				Timestamp: 1640875484,
			}, {
				Price:     47237,
				Weight:    1,
				Exchange:  "bitfinex",
				Timestamp: 1640875484,
			},
		},
	}

	parseLogs = []vo.PARTY_PRICE_INFO{
		getPriceAllParseLog,
		getPriceParseLog,
		getPartyPriceParseLog,
		getHistoryPriceParseLog,
		getBulkPricesParseLog,
		getBulkCurrencyPricesParseLog,
	}

	logInfo = vo.LOG_INFO{
		ClientIP:     "127.0.0.1",
		RequestTime:  "2021-12-16 15:39:46",
		UserAgent:    "insomnia/2021.7.1",
		RequestUrl:   "/api/getPrice/btcusdt/huobi",
		ResponseTime: "2021-12-16 15:39:46",
		Response:     "{\"code\":0,\"message\":\"OK\",\"data\":{\"price\":48903.4,\"timestamp\":1639640325}}",
	}

	mapInfo = map[string]interface{}{
		"request_client_ip":  logInfo.ClientIP,
		"request_time":       logInfo.RequestTime,
		"request_ua":         logInfo.UserAgent,
		"request_uri":        logInfo.RequestUrl,
		"response_time":      logInfo.RequestTime,
		"response":           logInfo.Response,
		"request_timestamp":  1640914947,
		"response_timestamp": 1640914947,
	}
)

func TestRequestInfoService_GetLogInfos(t *testing.T) {
	type fields struct {
		logInfoRepo *mock_repository.MockLogInfoRepository
	}
	type args struct {
		idx      int
		pageSize int
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    vo.LOG_INFOS
		wantErr bool
	}{
		{
			name: "basic",
			prepare: func(f *fields) {
				f.logInfoRepo.EXPECT().GetLogInfo(gomock.Eq(0), gomock.Eq(20)).Return(logInfos, nil)
			},
			args: args{
				idx:      0,
				pageSize: 20,
			},
			want:    logInfos,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				logInfoRepo: mock_repository.NewMockLogInfoRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &RequestInfoService{
				logInfoRepo: f.logInfoRepo,
			}
			got, err := s.GetLogInfos(tt.args.idx, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogInfos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLogInfos() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestInfoService_GetRequestInfoBySymbol(t *testing.T) {
	type fields struct {
		logInfoRepo *mock_repository.MockLogInfoRepository
	}
	type args struct {
		idx      int
		pageSize int
		symbol   string
		ip       string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    int
		want1   []vo.PARTY_PRICE_INFO
		wantErr bool
	}{
		{
			name: "basic",
			prepare: func(f *fields) {
				f.logInfoRepo.EXPECT().GetTotalLogInfoBySymbol(gomock.Eq("btcusdt"), gomock.Eq("")).Return(10, nil)
				f.logInfoRepo.EXPECT().GetLogInfoBySymbol(gomock.Eq(0), gomock.Eq(20), gomock.Eq("btcusdt"), gomock.Eq("")).Return(reqLogInfos, nil)
			},
			args: args{
				idx:      0,
				pageSize: 20,
				symbol:   "btcusdt",
				ip:       "",
			},
			want:    10,
			want1:   parseLogs,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				logInfoRepo: mock_repository.NewMockLogInfoRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &RequestInfoService{
				logInfoRepo: f.logInfoRepo,
			}
			got, got1, err := s.GetRequestInfoBySymbol(tt.args.idx, tt.args.pageSize, tt.args.symbol, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRequestInfoBySymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRequestInfoBySymbol() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetRequestInfoBySymbol() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRequestInfoService_InsertLogInfo(t *testing.T) {
	type fields struct {
		logInfoRepo *mock_repository.MockLogInfoRepository
	}
	type args struct {
		mapInfo map[string]interface{}
		t       int
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
				f.logInfoRepo.EXPECT().InsertLogInfo(gomock.Eq(mapInfo), gomock.Eq(1)).Return(nil)
			},
			args: args{
				mapInfo: mapInfo,
				t:       1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				logInfoRepo: mock_repository.NewMockLogInfoRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &RequestInfoService{
				logInfoRepo: f.logInfoRepo,
			}
			if err := s.InsertLogInfo(tt.args.mapInfo, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("InsertLogInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestInfoService_parseLogInfos(t *testing.T) {
	type fields struct {
		logInfoRepo *mock_repository.MockLogInfoRepository
	}
	type args struct {
		logInfos []vo.REQ_RSP_LOG_INFO
		symbol   string
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    []vo.PARTY_PRICE_INFO
	}{
		{
			name:    "basic",
			prepare: nil,
			args: args{
				logInfos: reqLogInfos,
				symbol:   "btcusdt",
			},
			want: parseLogs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				logInfoRepo: mock_repository.NewMockLogInfoRepository(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			s := &RequestInfoService{
				logInfoRepo: f.logInfoRepo,
			}
			if got := s.parseLogInfos(tt.args.logInfos, tt.args.symbol); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLogInfos() = %v, want %v", got, tt.want)
			}
		})
	}
}
