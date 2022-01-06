# get_price
This project is used for get price

## Getting the source
Your can get the source from github, run
``` bash
 git clone https://github.com/aresprotocols/get_price.git
```

## Building the source
Building project requires a go.

### Install with ubuntu
```bash
## install
wget https://studygolang.com/dl/golang/go1.13.4.linux-amd64.tar.gz
tar xfz go1.13.4.linux-amd64.tar.gz -C /usr/local
## config
vim ~/.bashrc
export GOPATH=/usr/local/go
export PATH=$GOPATH/bin:$PATH
source ～/.bashrc
```

### Install with mac
```bash
brew install go
```

Once you install, then
```bash
cd price_server
go build
```
That's all

## Configuration
Before run, you should config your project.
Using editor you like, such as
```bash
vim conf.toml
```
There is some config you need know.

> port = 5566 # server listen, default is 5566

> proxy = "http://127.0.0.1:7890/"   #  your proxy ip and port, or

> proxy = "" # Not use proxy


> symbols = ["btc-usdt", "eth-usdt", "dot-usdt", "xrp-usdt"]  #Transaction pair you interesting

> [mysql] # add your mysql config in here,db mean database name, just use a name you like

## Configure mysql
Install mysql server and start.

### Install with ubuntu
```bash
sudo apt update
sudo apt install mysql-server
```

### Install with mac
```bash
brew install mysql
```

If version >= 8.0, config with:
```bash
mysql -u root -p
use mysql;
GRANT ALL ON *.* TO 'root'@'%';
flush privileges;
ALTER USER 'root'@'localhost' IDENTIFIED BY '$yourpassword$' PASSWORD EXPIRE NEVER;
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY '$yourpassword$';
FLUSH PRIVILEGES;
```

## Start
Run
```bash
./start.sh
```

## Api
### Get exchange price
>
> http://127.0.0.1:5566/api/getPrice/$symbol/$market
>
> example: http://127.0.0.1:5566/api/getPrice/btcusdt/huobi

**Return**

``` javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"price": 43409.26,
		"timestamp": 1641435137
	}
}
```




### Get price after aggregation
>
> http://127.0.0.1:5566/api/getPartyPrice/$symbol
>
> example: http://127.0.0.1:5566/api/getPartyPrice/btcusdt
>
**Return**

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"price": 43340.53,
		"timestamp": 1641435560,
		"infos": [
			{
				"price": 43345.31,
				"weight": 2,
				"exchangeName": "huobi"
			},
			{
				"price": 43341,
				"weight": 1,
				"exchangeName": "kucoin"
			},
			{
				"price": 43340,
				"weight": 1,
				"exchangeName": "binance"
			},
			{
				"price": 43337.56,
				"weight": 1,
				"exchangeName": "bitstamp"
			},
			{
				"price": 43334,
				"weight": 1,
				"exchangeName": "bitfinex"
			}
		]
	}
}
```


### Get all price by symbol
>
> http://127.0.0.1:5566/api/getPriceAll/$symbol
>
> example: http://127.0.0.1:5566/api/getPriceAll/btcusdt
>
**Return**

``` javascript
{
	"code": 0,
	"message": "OK",
	"data": [
		{
			"name": "kucoin",
			"symbol": "btcusdt",
			"price": 43383,
			"timestamp": 1641435319,
			"weight": 1
		},
		{
			"name": "ok",
			"symbol": "btcusdt",
			"price": 43378.9,
			"timestamp": 1641435319,
			"weight": 1
		},
		{
			"name": "binance",
			"symbol": "btcusdt",
			"price": 43378.32,
			"timestamp": 1641435319,
			"weight": 1
		},
		{
			"name": "huobi",
			"symbol": "btcusdt",
			"price": 43377.56,
			"timestamp": 1641435319,
			"weight": 2
		},
		{
			"name": "coinbase",
			"symbol": "btcusdt",
			"price": 43376.92,
			"timestamp": 1641435319,
			"weight": 3
		},
		{
			"name": "bitfinex",
			"symbol": "btcusdt",
			"price": 43367,
			"timestamp": 1641435319,
			"weight": 1
		},
		{
			"name": "bitstamp",
			"symbol": "btcusdt",
			"price": 43337.56,
			"timestamp": 1641435319,
			"weight": 1
		}
	]
}
```


### Get coin history price (price must be stored in memory or db)

> http://127.0.0.1:5566/api/getHistoryPrice/$symbol?timestamp={}
>
> example:http://127.0.0.1:5566/api/getHistoryPrice/btcusdt?timestamp=1629341127
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"price": 43376.2625,
		"timestamp": 1641435319,
		"infos": [
			{
				"price": 43378.9,
				"weight": 1,
				"exchangeName": "ok"
			},
			{
				"price": 43378.32,
				"weight": 1,
				"exchangeName": "binance"
			},
			{
				"price": 43377.56,
				"weight": 2,
				"exchangeName": "huobi"
			},
			{
				"price": 43376.92,
				"weight": 3,
				"exchangeName": "coinbase"
			},
			{
				"price": 43367,
				"weight": 1,
				"exchangeName": "bitfinex"
			}
		]
	}
}
```



### bulk get symbol price

> http://127.0.0.1:5566/api/getBulkPrices?symbol={}
>
> example:http://127.0.0.1:5566/api/getBulkPrices?symbol=btcusdt_ethusdt
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"btcusdt": {
			"price": 43341.02125,
			"timestamp": 1641435801
		},
		"ethusdt": {
			"price": 3512.214,
			"timestamp": 1641435833
		}
	}
}
```


### bulk get symbol price by currency

> http://127.0.0.1:5566/api/getBulkCurrencyPrices?symbol={}
>
> example:http://127.0.0.1:5566/api/getBulkCurrencyPrices?symbol=btc_eth_dot_link&currency=usdt
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"btcusdt": {
			"price": 43372.14375,
			"timestamp": 1641435862,
			"infos": [
				{
					"price": 43374,
					"weight": 1,
					"exchangeName": "bitfinex"
				},
				{
					"price": 43373.05,
					"weight": 3,
					"exchangeName": "coinbase"
				},
				{
					"price": 43371.6,
					"weight": 2,
					"exchangeName": "huobi"
				},
				{
					"price": 43370.7,
					"weight": 1,
					"exchangeName": "kucoin"
				},
				{
					"price": 43370.1,
					"weight": 1,
					"exchangeName": "ok"
				}
			]
		},
		"dotusdt": {
			"price": 26.262575,
			"timestamp": 1641435881,
			"infos": [
				{
					"price": 26.2722,
					"weight": 1,
					"exchangeName": "kucoin"
				},
				{
					"price": 26.2681,
					"weight": 1,
					"exchangeName": "huobi"
				},
				{
					"price": 26.26,
					"weight": 1,
					"exchangeName": "binance"
				},
				{
					"price": 26.25,
					"weight": 1,
					"exchangeName": "coinbase"
				}
			]
		},
		"ethusdt": {
			"price": 3512.214,
			"timestamp": 1641435833,
			"infos": [
				{
					"price": 3513.06,
					"weight": 1,
					"exchangeName": "huobi"
				},
				{
					"price": 3512.95,
					"weight": 1,
					"exchangeName": "binance"
				},
				{
					"price": 3512.5,
					"weight": 1,
					"exchangeName": "coinbase"
				},
				{
					"price": 3512.26,
					"weight": 1,
					"exchangeName": "ok"
				},
				{
					"price": 3510.3,
					"weight": 1,
					"exchangeName": "bitfinex"
				}
			]
		},
		"linkusdt": {
			"price": 23.68055,
			"timestamp": 1641435836,
			"infos": [
				{
					"price": 23.681,
					"weight": 1,
					"exchangeName": "ok"
				},
				{
					"price": 23.6801,
					"weight": 1,
					"exchangeName": "kucoin"
				}
			]
		}
	}
}
```
### Get getReqConfig

> http://127.0.0.1:5566/api/getReqConfig

 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"aave-usdt": {
			"weight": [
				{
					"exchange": "binance",
					"weight": 2
				},
				{
					"exchange": "huobi",
					"weight": 4
				},
				{
					"exchange": "ok",
					"weight": 1
				},
				{
					"exchange": "kucoin",
					"weight": 1
				}
			],
			"interval": 60
		},
		"btc-usdt": {
			"weight": [
				{
					"exchange": "bitfinex",
					"weight": 1
				},
				{
					"exchange": "binance",
					"weight": 1
				},
				{
					"exchange": "huobi",
					"weight": 2
				},
				{
					"exchange": "ok",
					"weight": 1
				},
				{
					"exchange": "coinbase",
					"weight": 3
				},
				{
					"exchange": "kucoin",
					"weight": 1
				},
				{
					"exchange": "bitstamp",
					"weight": 1
				}
			],
			"interval": 60
		},
		"eth-usdt": {
			"weight": [
				{
					"exchange": "huobi",
					"weight": 1
				},
				{
					"exchange": "binance",
					"weight": 1
				},
				{
					"exchange": "bitfinex",
					"weight": 1
				},
				{
					"exchange": "ok",
					"weight": 1
				},
				{
					"exchange": "coinbase",
					"weight": 1
				},
				{
					"exchange": "kucoin",
					"weight": 1
				},
				{
					"exchange": "bitstamp",
					"weight": 1
				}
			],
			"interval": 60
		}
	}
}
```



### Get log info

> http://127.0.0.1:5566/api/getRequestInfo?index={}
>
> Header:
>
>     Authorization: Bearer TOKEN
>
> 
> example:http://127.0.0.1:5566/api/getRequestInfo?index=0
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"infos": [
			{
				"client_ip": "158.247.224.166",
				"request_time": "2022-01-06 02:35:42",
				"user_agent": "",
				"request_url": "/api/getBulkCurrencyPrices?currency=usdt&symbol=link",
				"response_time": "2022-01-06 02:35:42",
				"response": "{\"code\":0,\"message\":\"OK\",\"data\":{\"linkusdt\":{\"price\":23.7733,\"timestamp\":1641436498,\"infos\":[{\"price\":23.7761,\"weight\":1,\"exchangeName\":\"huobi\"},{\"price\":23.7705,\"weight\":1,\"exchangeName\":\"kucoin\"}]}}}"
			},
		]
	}
}
```

### get symbol request info

> http://127.0.0.1:5566/api/getRequestInfoBySymbol?index={}&symbol={}&ip={}
> 
> ip is optional
>
> example:http://127.0.0.1:5566/api/getRequestInfoBySymbol?index=0&symbol=btcusdt
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"curPage": 0,
		"totalNum": 207485,
		"items": [
			{
				"type": "getBulkCurrencyPrices",
				"client": {
					"ip": "45.77.30.9",
					"request_time": "2022-01-06 02:42:24",
					"request_timestamp": 1641436944
				},
				"price_info": {
					"price": 43310.613333,
					"timestamp": 1641436888
				},
				"price_infos": [
					{
						"price": 43320.1,
						"timestamp": 1641436888,
						"exchange": "ok",
						"weight": 1
					},
					{
						"price": 43309.9,
						"timestamp": 1641436888,
						"exchange": "kucoin",
						"weight": 1
					},
					{
						"price": 43309.17,
						"timestamp": 1641436888,
						"exchange": "binance",
						"weight": 1
					},
					{
						"price": 43309.07,
						"timestamp": 1641436888,
						"exchange": "huobi",
						"weight": 2
					},
					{
						"price": 43306.37,
						"timestamp": 1641436888,
						"exchange": "bitstamp",
						"weight": 1
					}
				]
			},
			
		]
	}
}
```


### Get http error info

> http://127.0.0.1:5566/api/getHttpErrorInfo?index={}
>
> example:http://127.0.0.1:5566/api/getHttpErrorInfo?index=0
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"curPage": 0,
		"totalNum": 674,
		"items": [
			{
				"url": "https://api-pub.bitfinex.com/v2/tickers?symbols=t{$symbol}",
				"symbol": "btc-usdt",
				"error": "status code :429 url:https://api-pub.bitfinex.com/v2/tickers?symbols=tBTCUSD",
				"Timestamp": 1641283069
			},
			{
				"url": "https://api-pub.bitfinex.com/v2/tickers?symbols=t{$symbol}",
				"symbol": "btc-usdt",
				"error": "Get \"https://api-pub.bitfinex.com/v2/tickers?symbols=tBTCUSD\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)",
				"Timestamp": 1640707125
			},
			{
				"url": "https://api.huobi.pro/market/detail/merged?symbol={$symbol}",
				"symbol": "btc-usdt",
				"error": "Get \"https://api.huobi.pro/market/detail/merged?symbol=btcusdt\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)",
				"Timestamp": 1640707125
			},
			{
				"url": "https://api.binance.com/api/v3/ticker/price?symbol={$symbol}",
				"symbol": "btc-usdt",
				"error": "Get \"https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT\": read tcp 141.164.58.241:58974->52.84.252.143:443: read: connection reset by peer",
				"Timestamp": 1640707124
			},
			
		]
	}
}
```

### Get local prices

> http://127.0.0.1:5566/api/getLocalPrices?index={}&symbol={}
>
> example:http://127.0.0.1:5566/api/getLocalPrices?index=0&symbol=btcusdt
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"PriceInfosCache": {
			"btcusdt": [
				{
					"PriceInfos": [
						{
							"symbol": "btcusdt",
							"price": 46512.1,
							"priceOrigin": "kucoin",
							"weight": 1,
							"timestamp": 1641283060
						},
						{
							"symbol": "btcusdt",
							"price": 46508.92,
							"priceOrigin": "huobi",
							"weight": 2,
							"timestamp": 1641283060
						},
						{
							"symbol": "btcusdt",
							"price": 46507.77,
							"priceOrigin": "binance",
							"weight": 1,
							"timestamp": 1641283060
						},
						{
							"symbol": "btcusdt",
							"price": 46507.08,
							"priceOrigin": "coinbase",
							"weight": 3,
							"timestamp": 1641283060
						},
						{
							"symbol": "btcusdt",
							"price": 46505.5,
							"priceOrigin": "ok",
							"weight": 1,
							"timestamp": 1641283060
						},
						{
							"symbol": "btcusdt",
							"price": 46286.21,
							"priceOrigin": "bitstamp",
							"weight": 1,
							"timestamp": 1641283060
						}
					]
				},
				{
					"PriceInfos": [
						{
							"symbol": "btcusdt",
							"price": 46513,
							"priceOrigin": "bitfinex",
							"weight": 1,
							"timestamp": 1641283129
						},
						{
							"symbol": "btcusdt",
							"price": 46507,
							"priceOrigin": "kucoin",
							"weight": 1,
							"timestamp": 1641283129
						},
						{
							"symbol": "btcusdt",
							"price": 46505.97,
							"priceOrigin": "coinbase",
							"weight": 3,
							"timestamp": 1641283129
						},
						{
							"symbol": "btcusdt",
							"price": 46503.8,
							"priceOrigin": "ok",
							"weight": 1,
							"timestamp": 1641283129
						},
						{
							"symbol": "btcusdt",
							"price": 46502.45,
							"priceOrigin": "binance",
							"weight": 1,
							"timestamp": 1641283129
						},
						{
							"symbol": "btcusdt",
							"price": 46502.09,
							"priceOrigin": "huobi",
							"weight": 2,
							"timestamp": 1641283129
						},
						{
							"symbol": "btcusdt",
							"price": 46286.21,
							"priceOrigin": "bitstamp",
							"weight": 1,
							"timestamp": 1641283129
						}
					]
				}				
			]
		}
	}
}
```



### Get update price history

> http://127.0.0.1:5566/api/getUpdatePriceHistory?index={}&symbol={}
>
> example:http://127.0.0.1:5566/api/getUpdatePriceHistory?index=0&symbol=btcusdt
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"curPage": 0,
		"totalNum": 35787,
		"items": [
			{
				"timestamp": 1641438277,
				"symbol": "btcusdt",
				"price": 43299.32,
				"infos": [
					{
						"symbol": "btcusdt",
						"price": 43306.37,
						"priceOrigin": "bitstamp",
						"weight": 1,
						"timestamp": 1641438277
					},
					{
						"symbol": "btcusdt",
						"price": 43301.9,
						"priceOrigin": "kucoin",
						"weight": 1,
						"timestamp": 1641438277
					},
					{
						"symbol": "btcusdt",
						"price": 43301.3,
						"priceOrigin": "ok",
						"weight": 1,
						"timestamp": 1641438277
					},
					{
						"symbol": "btcusdt",
						"price": 43300.43,
						"priceOrigin": "binance",
						"weight": 1,
						"timestamp": 1641438277
					},
					{
						"symbol": "btcusdt",
						"price": 43299.8,
						"priceOrigin": "huobi",
						"weight": 2,
						"timestamp": 1641438277
					},
					{
						"symbol": "btcusdt",
						"price": 43297.11,
						"priceOrigin": "coinbase",
						"weight": 3,
						"timestamp": 1641438277
					},
					{
						"symbol": "btcusdt",
						"price": 43280,
						"priceOrigin": "bitfinex",
						"weight": 1,
						"timestamp": 1641438277
					}
				]
			},
			{
				"timestamp": 1641437976,
				"symbol": "btcusdt",
				"price": 43233.92875,
				"infos": [
					{
						"symbol": "btcusdt",
						"price": 43306.37,
						"priceOrigin": "bitstamp",
						"weight": 1,
						"timestamp": 1641437976
					},
					{
						"symbol": "btcusdt",
						"price": 43248.6,
						"priceOrigin": "ok",
						"weight": 1,
						"timestamp": 1641437976
					},
					{
						"symbol": "btcusdt",
						"price": 43242.5,
						"priceOrigin": "kucoin",
						"weight": 1,
						"timestamp": 1641437976
					},
					{
						"symbol": "btcusdt",
						"price": 43234.09,
						"priceOrigin": "binance",
						"weight": 1,
						"timestamp": 1641437976
					},
					{
						"symbol": "btcusdt",
						"price": 43229.32,
						"priceOrigin": "huobi",
						"weight": 2,
						"timestamp": 1641437976
					},
					{
						"symbol": "btcusdt",
						"price": 43229.2,
						"priceOrigin": "coinbase",
						"weight": 3,
						"timestamp": 1641437976
					},
					{
						"symbol": "btcusdt",
						"price": 43220,
						"priceOrigin": "bitfinex",
						"weight": 1,
						"timestamp": 1641437976
					}
				]
			}
		]
	}
}
```


### set weight

> POST http://127.0.0.1:5566/api/setWeight
> 
> Header:  
> 
>       Authorization: Bearer TOKEN
> 
> Body: {"symbol":"btc-usdt", "exchange":"huobi", "weight":2 }
>
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": null
}
```

### Get ares info

> http://127.0.0.1:5566/api/getAresAll
>
> example:http://127.0.0.1:5566/api/getAresAll
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"price": 0.0129,
		"percent_change": -10.29,
		"rank": 1438,
		"market_cap": 3004786.5293474947,
		"volume": 1751021.98137502
	}
}
```

### Get Dex price info

> http://127.0.0.1:5566/api/getDexPrice
>
> example:http://127.0.0.1:5566/api/getDexPrice
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"uni": {
			"price": "0.01368310053",
			"timestamp": 1641438741
		},
		"pancake": {
			"price": "0.01426012838",
			"timestamp": 1641438738
		}
	}
}
```



### auth

> POST http://127.0.0.1:5566/api/auth
>
> body: {"user":"root", "password":"xxxxxx"}
> 
> password need md5
> 
**Return**

```javascript
{
	"code": 0,
	"message": "OK",
	"data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJvb3QiLCJleHAiOjE2NDE0NDM3MDUsImlzcyI6ImdldC1wcmljZSJ9.VKz70nqgizZbMtUwUt-z4_pIHjj2stbuWVQX8ULrG3c"
}
```




### Get update price heartbeat
>
> http://127.0.0.1:5566/api/getUpdatePriceHeartbeat/$symbol
>
> example: http://127.0.0.1:5566/api/getUpdatePriceHeartbeat/btcusdt
>
**Return**

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"expect_resources": 7,
		"actual_resources": 7,
		"latest_timestamp": 1641439787,
		"interval": 60
	}
}
```



### bulk get symbol state

> http://127.0.0.1:5566/api/getBulkSymbolsState?symbol={}
>
> example:http://127.0.0.1:5566/api/getBulkSymbolsState?symbol=btc_eth_dot_link&currency=usdt
>

**Return**

```javascript
{
	"code": 0,
	"message": "OK",
	"data": {
		"btcusdt": true,
		"dotusdt": true,
		"ethusdt": true,
		"linkusdt": true
	}
}
```


### set update price interval

> POST http://127.0.0.1:5566/api/setInterval
> 
> Header:  
> 
>       Authorization: Bearer TOKEN
> 
> Body: {"symbol":"btc-usdt", "interval":20 }
>
>
 **Return** 

```javascript
{
	"code": 0,
	"message": "OK",
	"data": null
}
```