package exchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	urlStr := "https://data.gateapi.io/api2/1/ticker/ares_usdt"
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)

	var aresInfo AresGateShowInfo

	err = json.Unmarshal(body, &aresInfo)
	fmt.Println("aresInfo ", aresInfo.Price)
	fmt.Println("aresInfo ", aresInfo)

	price, _ := strconv.ParseFloat(aresInfo.Price, 64)
	percentChange, _ := strconv.ParseFloat(aresInfo.PercentChange, 64)
	marketCap, _ := strconv.ParseFloat(aresInfo.MarketCap, 64)
	volume, _ := strconv.ParseFloat(aresInfo.Volume, 64)

	v := AresShowInfo{
		Price:         price,
		PercentChange: percentChange,
		MarketCap:     marketCap,
		Volume:        volume,
	}
	fmt.Println("v ", v)
}

func TestGetAres(t *testing.T) {
	urlStr := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?slug=ares-protocol&CMC_PRO_API_KEY=64a35d97-aca1-4c5c-8c17-43864a23aa97"
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)
}

func TestTime(t *testing.T) {
	fmt.Println("", time.Now())
	fmt.Println("", time.Now().Minute())
	var aresShowInfo AresShowInfo
	fmt.Println("", aresShowInfo)
	info, _ := GetGateAresInfo("")
	fmt.Println("info ", info)
	info, _ = GetGateAresInfo("")
	fmt.Println("info ", info)
	time.Sleep(time.Minute)
	info, _ = GetGateAresInfo("")
	fmt.Println("info Sleep ", info)
}
