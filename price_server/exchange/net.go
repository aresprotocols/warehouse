package exchange

import (
	"errors"
	"io/ioutil"
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
		}
		httpClient = &http.Client{
			Timeout:   time.Second * 10,
			Transport: netTransport,
		}
	}

	res, err := httpClient.Get(exchangeUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("status code :" + strconv.Itoa(res.StatusCode))
	}
	c, _ := ioutil.ReadAll(res.Body)
	return string(c), nil
}
