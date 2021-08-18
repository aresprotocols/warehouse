# get_price

## Step
1. cd price_server
2. go build
3. ./start.sh

## config
config file: conf.toml

proxy = "" or "your proxy ip:port"
symbols = ["btc-usdt", "eth-usdt", "dot-usdt", "xrp-usdt", "add token here"]

## api
base url:http://127.0.0.1:5566/api/
1. get coin history price (the price must store in db)
url : getHistoryPrice/$symbol
query param:timestamp
example:http://127.0.0.1:5566/api/getHistoryPrice/btcusdt?timestamp=1629197717

2. get exchange price
url : getprice/$symbol/$exchange
example: http://127.0.0.1:5566/api/getprice/btcusdt/huobi

3. get price after  aggregation
url : getPartyPrice/$symbol
example: http://127.0.0.1:5566/api/getPartyPrice/btcusdt

4. get all price by symbol
url : getPriceAll/$symbol
example: http://127.0.0.1:5566/api/getPriceAll/btcusdt

5. get weight by config
url : getConfigWeight
example: http://127.0.0.1:5566/api/getConfigWeight