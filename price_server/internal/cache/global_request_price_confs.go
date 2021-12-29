package cache

import (
	conf "price_api/price_server/config"
	"sync"
)

//go:generate mockgen -destination mock/global_request_price_confs_mock.go price_api/price_server/internal/cache GlobalRequestPriceConfs

type GlobalRequestPriceConfs interface {
	SetConfs(conf map[string][]conf.ExchangeConfig)
	UpdateSymbolWeight(symbol, exchange string, weight int)
	GetConfs() map[string][]conf.ExchangeConfig
	GetConfsBySymbol(symbol string) []conf.ExchangeConfig
}

func NewGlobalRequestPriceConfs() GlobalRequestPriceConfs {
	return &globalRequestPriceConfs{
		confs: make(map[string][]conf.ExchangeConfig),
		m:     new(sync.RWMutex)}
}

type globalRequestPriceConfs struct {
	confs map[string][]conf.ExchangeConfig
	m     *sync.RWMutex
}

func (c *globalRequestPriceConfs) SetConfs(conf map[string][]conf.ExchangeConfig) {
	c.m.Lock()
	c.confs = conf
	c.m.Unlock()
}

func (c *globalRequestPriceConfs) GetConfs() map[string][]conf.ExchangeConfig {
	c.m.Lock()
	defer c.m.Unlock()
	return c.confs
}
func (c *globalRequestPriceConfs) GetConfsBySymbol(symbol string) []conf.ExchangeConfig {
	c.m.Lock()
	defer c.m.Unlock()
	return c.confs[symbol]
}

func (c *globalRequestPriceConfs) UpdateSymbolWeight(symbol, exchange string, weight int) {
	c.m.Lock()
	for i, confTemp := range c.confs[symbol] {
		if confTemp.Name == exchange {
			c.confs[symbol][i].Weight = int64(weight)
			break
		}
	}
	c.m.Unlock()
}
