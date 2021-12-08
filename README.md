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
mysql
use mysql;
GRANT ALL ON *.* TO 'root'@'%';
flush privileges;
ALTER USER 'root'@'localhost' IDENTIFIED BY 'yourpassword' PASSWORD EXPIRE NEVER;
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'yourpassword';
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
> http://127.0.0.1:5566/api/getprice/$symbol/$market
>
> example: http://127.0.0.1:5566/api/getprice/btcusdt/huobi

**Return**

``` javascript
{"code":0,"message":"OK","data":{"timestamp":1629340675,"price":44721.54}}
```




### Get price after aggregation
>
> http://127.0.0.1:5566/api/getPartyPrice/$symbol
>
> example: http://127.0.0.1:5566/api/getPartyPrice/btcusdt
>
**Return**

```javascript
{"code":0,"message":"OK","data":{"price":44727.4,"timestamp":1629340811,"infos":[{"price":44731.7,"weight":1,"exchangeName":"ok"},
{"price":44726.48,"weight":1,"exchangeName":"huobi"},{"price":44720,"weight":1,"exchangeName":"bitfinex"},{"price":44732.52,"weight":1,"exchangeName":"bitstamp"},
{"price":44726.3,"weight":1,"exchangeName":"binance"}]}}
```


### Get all price by symbol
>
> http://127.0.0.1:5566/api/getPriceAll/$symbol
>
> example: http://127.0.0.1:5566/api/getPriceAll/btcusdt
>
**Return**

``` javascript
{"code":0,"message":"OK","data":[{"name":"ok","symbol":"btcusdt","price":54398.9,"timestamp":1633678924,"weight":1},{"name":"bitfinex","symbol":"btcusdt","price":54431,"timestamp":1633678924,"weight":1},{"name":"binance","symbol":"btcusdt","price":54393.67,"timestamp":1633678924,"weight":1},{"name":"huobi","symbol":"btcusdt","price":54396,"timestamp":1633678924,"weight":2},{"name":"bitstamp","symbol":"btcusdt","price":54368.34,"timestamp":1633678924,"weight":1}]}
```


### Get coin history price (price must be stored in memory or db)

> http://127.0.0.1:5566/api/getHistoryPrice/$symbol?timestamp={}
>
> example:http://127.0.0.1:5566/api/getHistoryPrice/btcusdt?timestamp=1629341127
>
 **Return** 

```javascript
{"code":0,"message":"OK","data":{"price":44655.439999999995,"timestamp":1629341547,"infos":[{"price":44655.27,"weight":1,"exchangeName":"cryptocompare"},
{"price":44652.4,"weight":1,"exchangeName":"ok"},{"price":44666,"weight":1,"exchangeName":"huobi"},{"price":44665.62,"weight":1,"exchangeName":"binance"},
{"price":44637.91,"weight":1,"exchangeName":"bitstamp"}]}}
```


### Get ares info

> http://127.0.0.1:5566/api/getAresAll
>
> example:http://127.0.0.1:5566/api/getAresAll
>
 **Return** 

```javascript
{"code":0,"message":"OK","data":{"price":0.04235333740536,"percent_change":-5.38960837,"rank":1108,"market_cap":6516779.946008743,"volume":749528.82939821}}
```

### Get symbol price

> http://127.0.0.1:5566/api/getBulkPrices?symbol={}
>
> example:http://127.0.0.1:5566/api/getBulkPrices?symbol=btcusdt_ethusdt
>
 **Return** 

```javascript
{"code":0,"message":"OK","data":{"btcusdt":{"price":42174.990000000005,"timestamp":1632279887},"ethusdt":{"price":2874.3959999999997,"timestamp":1632279887}}}

{"code":0,"message":"OK","data":{"arrusdt":{"price":0,"timestamp":0}}}
```


### Get symbol price

> http://127.0.0.1:5566/api/getBulkCurrencyPrices?symbol={}
>
> example:https://api.aresprotocol.io/api/getBulkCurrencyPrices?symbol=btc_eth_dot_link&currency=usdt
>
 **Return** 

```javascript
{"code":0,"message":"OK","data":{"btcusdt":{"price":57465.51,"timestamp":1638167488},"dotusdt":{"price":36.158,"timestamp":1638167488},"ethusdt":{"price":4335.455,"timestamp":1638167488},"linkusdt":{"price":25.005,"timestamp":1638167488}}}
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
{"code":0,"message":"OK","data":{"infos":[{"client_ip":"127.0.0.1","method":"GET","post_data":"","proto":"HTTP/1.1","request_time":"2021-09-23 16:37:32","user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36","request_url":"/api/getRequestInfo?index=1","response_time":"2021-09-23 16:37:32","response":"{\"code\":0,\"message\":\"OK\",\"data\":{\"infos\":null}}"}]}}
```

### Get local prices

> http://127.0.0.1:5566/api/getLocalPrices?index={}&symbol={}
>
> example:http://127.0.0.1:5566/api/getLocalPrices?index=0&symbol=ethusdt
>
 **Return** 

```javascript
{"code":0,"message":"OK","data":{"PriceInfosCache":[{"PriceInfos":[{"Symbol":"ethusdt","Price":3126.18,"PriceOrigin":"bitstamp","Weight":1,"TimeStamp":1632722382},{"Symbol":"ethusdt","Price":3130.29,"PriceOrigin":"coinbase","Weight":1,"TimeStamp":1632722382},{"Symbol":"ethusdt","Price":3134.2,"PriceOrigin":"bitfinex","Weight":1,"TimeStamp":1632722382},{"Symbol":"ethusdt","Price":3130.65,"PriceOrigin":"huobi","Weight":1,"TimeStamp":1632722382},{"Symbol":"ethusdt","Price":3130.4,"PriceOrigin":"binance","Weight":1,"TimeStamp":1632722382},{"Symbol":"ethusdt","Price":3131.31,"PriceOrigin":"ok","Weight":1,"TimeStamp":1632722382}]}]}}
```

### Get getReqConfig

> http://127.0.0.1:5566/api/getReqConfig

 **Return** 

```javascript
{"code":0,"message":"OK","data":{"1INCH-usdt":["huobi","binance","cryptocompare"],"aave-usdt":["binance","ok","huobi"],"ada-usdt":["ok","huobi","bitfinex","binance"],"algo-usdt":["ok","huobi"],"atom-usdt":["huobi","binance","ok","cryptocompare"],"avax-usdt":["cryptocompare","huobi","ok"],"axs-usdt":["bitfinex","coinbase","huobi","binance","ok"],"bat-usdt":["cryptocompare","huobi"],"bch-usdt":["huobi","cryptocompare","binance","ok"],"bnt-usdt":["binance","huobi","ok","bitfinex"],"btc-usdt":["ok","bitstamp","coinbase","huobi","cryptocompare","bitfinex"],"btt-usdt":["binance","ok"],"celo-usdt":["binance","ok","cryptocompare"],"chz-usdt":["cryptocompare","binance","ok","coinbase"],"comp-usdt":["cryptocompare","ok"],"crv-usdt":["ok","binance","huobi"],"dash-usdt":["huobi","cryptocompare","ok"],"dcr-usdt":["huobi","ok","cryptocompare","bitfinex"],"doge-usdt":["huobi","binance","ok","coinbase"],"dot-usdt":["huobi","bitfinex","coinbase","binance","ok"],"egld-usdt":["ok"],"enj-usdt":["huobi","ok"],"eos-usdt":["binance","huobi","ok","bitfinex"],"etc-usdt":["binance","huobi","bitfinex","cryptocompare"],"eth-usdt":["bitstamp","huobi","coinbase","binance"],"fet-usdt":["binance","bitfinex","coinbase"],"fil-usdt":["binance","huobi","ok"],"ftm-usdt":["ok","binance","bitfinex"],"ftt-usdt":["binance","huobi","cryptocompare"],"grt-usdt":["huobi","ok","cryptocompare","binance","bitfinex"],"hbar-usdt":["binance","huobi","ok"],"icp-usdt":["huobi","bitfinex","cryptocompare"],"icx-usdt":["binance","ok","huobi"],"iost-usdt":["binance","huobi","ok"],"iota-usdt":["huobi","binance"],"iotx-usdt":["huobi","coinbase","cryptocompare"],"kava-usdt":["huobi","binance"],"ksm-usdt":["binance","cryptocompare","bitfinex","ok"],"link-usdt":["ok","binance"],"lrc-usdt":["cryptocompare","binance","bitfinex","ok","huobi"],"ltc-usdt":["huobi","ok","binance"],"luna-usdt":["binance","huobi","ok","cryptocompare"],"mana-usdt":["binance","ok","huobi","cryptocompare"],"matic-usdt":["bitstamp","ok"],"mkr-usdt":["bitfinex","ok","binance","cryptocompare"],"nano-usdt":["cryptocompare","ok","huobi"],"near-usdt":["binance","ok","huobi"],"neo-usdt":["bitfinex","cryptocompare","ok","huobi","binance"],"omg-usdt":["bitfinex","cryptocompare","ok","huobi","binance"],"ont-usdt":["binance","ok"],"qtum-usdt":["binance","cryptocompare","huobi","ok"],"ren-usdt":["binance","huobi","cryptocompare","ok"],"sand-usdt":["binance","huobi","ok","cryptocompare"],"sc-usdt":["huobi","binance","ok"],"snx-usdt":["huobi","binance","ok","bitfinex"],"sol-usdt":["coinbase","ok","bitfinex","huobi","binance"],"srm-usdt":["cryptocompare","binance"],"stx-usdt":["ok"],"sushi-usdt":["binance","ok","huobi"],"theta-usdt":["ok","binance"],"trx-usdt":["bitfinex","huobi","binance"],"uma-usdt":["binance","huobi","ok"],"uni-usdt":["huobi","bitfinex","binance","ok","cryptocompare"],"vet-usdt":["binance"],"waves-usdt":["binance","cryptocompare","huobi","ok"],"xem-usdt":["binance","cryptocompare","huobi"],"xlm-usdt":["binance","bitfinex","cryptocompare","huobi"],"xmr-usdt":["huobi","binance","bitfinex"],"xrp-usdt":["bitfinex","bitstamp","binance","ok"],"xtz-usdt":["bitfinex","huobi","binance","ok"],"yfi-usdt":["bitfinex","huobi","binance","ok"],"zec-usdt":["huobi","ok","binance","bitfinex"],"zen-usdt":["huobi","coinbase","cryptocompare","binance","ok"],"zil-usdt":["binance","bitfinex","ok"],"zrx-usdt":["ok","bitfinex","huobi","cryptocompare","binance"]}}
```


### Get http error info

> http://127.0.0.1:5566/api/getHttpErrorInfo?index={}
>
> example:http://127.0.0.1:5566/api/getHttpErrorInfo?index=0
>
 **Return** 

```javascript
{"code":0,"message":"OK","data":{"infos":[{"url":"https://www.okex.com/api/spot/v3/instruments/{$symbol1}-{$symbol2}/ticker","symbol":"btc-usdt","error":"test error","Timestamp":1632727340},{"url":"https://api.huobi.pro/market/detail/merged?symbol={$symbol}","symbol":"btc-usdt","error":"test error","Timestamp":1632727288}]}}
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
> example:http://127.0.0.1:5566/api/setWeight?exchange=coinbase&symbol=btc-usdt&weight=3
>
 **Return** 

```javascript
{"code":0,"message":"OK","data":null}
```

### get symbol request info

> http://127.0.0.1:5566/api/getRequestInfoBySymbol?index={}&symbol={}
>
> example:http://127.0.0.1:5566/api/getRequestInfoBySymbol?index=0&symbol=btcusdt
>
 **Return** 

```javascript
{"code":0,"message":"OK","data":{"getBulkPrices":[{"client":{"ip":"127.0.0.1","request_time":"2021-09-29 11:03:26"},"price_info":{"price":41700.786666666674,"timestamp":1632884557}}],"getHistoryPrice":[{"client":{"ip":"127.0.0.1","request_time":"2021-09-29 11:02:55"},"price_info":{"price":45287.56999999999,"timestamp":1629268984},"price_infos":[{"price":45286.2,"timestamp":1629268984,"exchange":"ok","weight":1},{"price":45295,"timestamp":1629268984,"exchange":"bitfinex","weight":1},{"price":45283.32,"timestamp":1629268984,"exchange":"huobi","weight":1},{"price":45290.2,"timestamp":1629268984,"exchange":"bitstamp","weight":1},{"price":45283.13,"timestamp":1629268984,"exchange":"binance","weight":1}]}],"getPartyPrice":[{"client":{"ip":"127.0.0.1","request_time":"2021-09-29 14:00:27"},"price_info":{"price":42288.427500000005,"timestamp":1632895181},"price_infos":[{"price":42291.66,"timestamp":1632895181,"exchange":"huobi","weight":2},{"price":42294.39,"timestamp":1632895181,"exchange":"binance","weight":1},{"price":42276,"timestamp":1632895181,"exchange":"ok","weight":1}]},{"client":{"ip":"127.0.0.1","request_time":"2021-09-28 16:54:31"},"price_info":{"price":41721.92571428571,"timestamp":1632819233},"price_infos":[{"price":41679.61,"timestamp":1632819233,"exchange":"huobi","weight":2},{"price":41939.72,"timestamp":1632819233,"exchange":"bitstamp","weight":1},{"price":41750,"timestamp":1632819233,"exchange":"bitfinex","weight":1},{"price":41668.18,"timestamp":1632819233,"exchange":"coinbase","weight":3}]}],"getPrice":[{"client":{"ip":"127.0.0.1","request_time":"2021-09-28 16:55:53"},"price_info":{"price":41724,"timestamp":1632819294}}],"getPriceAll":[{"client":{"ip":"127.0.0.1","request_time":"2021-10-08 15:42:38"},"price_infos":[{"price":54398.9,"timestamp":1633678924,"exchange":"ok","weight":1},{"price":54431,"timestamp":1633678924,"exchange":"bitfinex","weight":1},{"price":54393.67,"timestamp":1633678924,"exchange":"binance","weight":1},{"price":54396,"timestamp":1633678924,"exchange":"huobi","weight":2},{"price":54368.34,"timestamp":1633678924,"exchange":"bitstamp","weight":1}]}]}}
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
{"code": 0, "message": "OK", "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJvb3QiLCJleHAiOjE2Mzg5MzY3MTEsImlzcyI6ImdldC1wcmljZSJ9.LEgHx3ZfbCNIf6UMxYFNDmg8cd5lu2oqLUk2Eb1TwOw"}
```
