package handler

import (
	"price_api/price_server/internal/config"
	"price_api/price_server/internal/exchange"
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
