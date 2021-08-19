# get_price
This project is used for get price

## Getting the source
Your can get the source from github, run

>  git clone https://github.com/aresprotocols/get_price.git


## Building the source
Building getrue requires a go.
You can install them using your favourite package manager.
Once you install, then

> 1.cd price_server

> 2.go build

That's all

## Configuration
Before run, you should config your project.
Using editor you like, such as

> vim conf.toml

There is some config you need know.

> port = 5566 # server listen, default is 5566

> proxy = "http://127.0.0.1:7890/"   #  your proxy ip and port, or

> proxy = "" # Not use proxy


> symbols = ["btc-usdt", "eth-usdt", "dot-usdt", "xrp-usdt"]  #Transaction pair you interesting

> [mysql] # add your mysql config in here,db mean database name, just use a name you like

## Configure mysql
Install mysql server and start.

If version >= 8.0, config with this:

`https://www.cnblogs.com/liran123/p/10164564.html`

## Start
Run
> ./start.sh

## Api
>  get weight by config
>
> http://127.0.0.1:5566/api/getConfigWeight
>
> @return
>
> {"code":0,"message":"OK","data":{"weightInfos":[{"exchangeName":"binance","weight":1},{"exchangeName":"huobi","weight":1},{"exchangeName":"bitfinex","weight":1},{"exchangeName":"ok","weight":1},{"exchangeName":"cryptocompare","weight":1},{"exchangeName":"coinbase","weight":1},{"exchangeName":"bitstamp","weight":1}]}}




> get exchange price
>
> http://127.0.0.1:5566/api/getprice/$symbol/$market
>
> example: http://127.0.0.1:5566/api/getprice/btcusdt/huobi
>
> @return
>
> {"code":0,"message":"OK","data":{"timestamp":1629340675,"price":44721.54}}





> get price after aggregation
>
> http://127.0.0.1:5566/api/getPartyPrice/$symbol
>
> example: http://127.0.0.1:5566/api/getPartyPrice/btcusdt
>
> @return
>
> {"code":0,"message":"OK","data":{"price":44727.4,"timestamp":1629340811,"infos":[{"price":44731.7,"weight":1,"exchangeName":"ok"},{"price":44726.48,"weight":1,"exchangeName":"huobi"},{"price":44720,"weight":1,"exchangeName":"bitfinex"},{"price":44732.52,"weight":1,"exchangeName":"bitstamp"},{"price":44726.3,"weight":1,"exchangeName":"binance"}]}}



> get all price by symbol
>
> http://127.0.0.1:5566/api/getPriceAll/$symbol
>
> example: http://127.0.0.1:5566/api/getPriceAll/btcusdt
>
> @return
>
> {"code":0,"message":"OK","data":[{"name":"binance","symbol":"btcusdt","price":44673.34,"timestamp":1629340944},{"name":"huobi","symbol":"btcusdt","price":44671.41,"timestamp":1629340944},{"name":"bitfinex","symbol":"btcusdt","price":44694,"timestamp":1629340944},{"name":"ok","symbol":"btcusdt","price":44674.4,"timestamp":1629340944},{"name":"cryptocompare","symbol":"btcusdt","price":44688.36,"timestamp":1629340944},{"name":"coinbase","symbol":"btcusdt","price":44667.16,"timestamp":1629340944},{"name":"bitstamp","symbol":"btcusdt","price":44663.78,"timestamp":1629340944}]}



> get coin history price (price must be stored in memory or db)
>
> http://127.0.0.1:5566/api/getHistoryPrice/$symbol?timestamp={}
>
> example:http://127.0.0.1:5566/api/getHistoryPrice/btcusdt?timestamp=1629341127
>
> @return 
>
> {"code":0,"message":"OK","data":{"price":44655.439999999995,"timestamp":1629341547,"infos":[{"price":44655.27,"weight":1,"exchangeName":"cryptocompare"},{"price":44652.4,"weight":1,"exchangeName":"ok"},{"price":44666,"weight":1,"exchangeName":"huobi"},{"price":44665.62,"weight":1,"exchangeName":"binance"},{"price":44637.91,"weight":1,"exchangeName":"bitstamp"}]}}


