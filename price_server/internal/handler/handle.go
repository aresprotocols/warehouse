package handler

import (
	conf "price_api/price_server/config"
	"price_api/price_server/exchange"
)

var (
	handle *Handle
)

type Handle struct {
	fetcher *exchange.Fetcher
}

func InitHandle(cfg conf.Config) *Handle {
	handle = &Handle{
		fetcher: exchange.InitFetcher(cfg),
	}

	handle.fetcher.Start()

	return handle
}

func (h *Handle) Stop() {
	h.fetcher.Stop()
}
