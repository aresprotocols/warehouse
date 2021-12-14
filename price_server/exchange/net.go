package exchange

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func getPrice(exchangeUrl string, proxyAddr string) (string, error) {
	var httpClient *http.Client
	if proxyAddr == "" {
		httpClient = &http.Client{
			Timeout: time.Second * 10,
		}
	} else {
		proxy, err := url.Parse(proxyAddr)
		if err != nil {
			return "", err
		}
		netTransport := &http.Transport{
			Proxy:                 http.ProxyURL(proxy),
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * time.Duration(5),
			DisableKeepAlives:     true,
		}
		httpClient = &http.Client{
			Timeout:   time.Second * 10,
			Transport: netTransport,
		}
	}

	res, err := httpClient.Get(exchangeUrl)
	if err != nil {
		log.Println(fmt.Sprintf("getPrice error:%v", err))
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("status code :" + strconv.Itoa(res.StatusCode) + " url:" + exchangeUrl)
	}
	c, _ := ioutil.ReadAll(res.Body)
	return string(c), nil
}
