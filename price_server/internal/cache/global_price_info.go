package cache

import (
	conf "price_api/price_server/config"
	"strings"
	"sync"
)

//go:generate mockgen -destination mock/global_price_info_mock.go price_api/price_server/internal/cache GlobalPriceInfoCache

type GlobalPriceInfoCache interface {
	GetLatestPriceInfos(symbol string) conf.PriceInfos
	GetPriceInfosEqualTimestamp(symbol string, timestamp int64) (bool, conf.PriceInfos)
	GetPriceInfosByRange(symbol string, start, end int) conf.PriceInfosCache
	GetCacheLength() int
	GetSymbolCacheLength(symbol string) int
	UpdateCachePrice(symbol string, infos conf.PriceInfos, maxMemTime int)
}

type globalPriceInfoCache struct {
	gPriceInfosCache conf.PriceInfosCache
	m                *sync.RWMutex
}

func NewGlobalPriceInfoCache() GlobalPriceInfoCache {
	return &globalPriceInfoCache{
		gPriceInfosCache: conf.PriceInfosCache{PriceInfosCache: map[string][]conf.PriceInfos{}},
		m:                new(sync.RWMutex)}
}

func (c *globalPriceInfoCache) GetLatestPriceInfos(symbol string) conf.PriceInfos {
	c.m.RLock()
	defer c.m.RUnlock()
	if len(c.gPriceInfosCache.PriceInfosCache) == 0 || len(c.gPriceInfosCache.PriceInfosCache[symbol]) == 0 {
		return conf.PriceInfos{}
	}
	latestInfos := c.gPriceInfosCache.PriceInfosCache[symbol][len(c.gPriceInfosCache.PriceInfosCache[symbol])-1]
	return latestInfos
}

func (c *globalPriceInfoCache) GetPriceInfosEqualTimestamp(symbol string, timestamp int64) (bool, conf.PriceInfos) {
	bMemory := false
	var cacheInfo conf.PriceInfos
	c.m.RLock()
	defer c.m.RUnlock()
	//latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	if len(c.gPriceInfosCache.PriceInfosCache) != 0 && len(c.gPriceInfosCache.PriceInfosCache[symbol]) != 0 {
		for i := len(c.gPriceInfosCache.PriceInfosCache[symbol]) - 1; i >= 0; i-- {
			info := c.gPriceInfosCache.PriceInfosCache[symbol][i]
			if len(info.PriceInfos) == 0 {
				continue
			}
			if info.PriceInfos[0].TimeStamp == timestamp {
				//use memory
				bMemory = true
				cacheInfo = c.gPriceInfosCache.PriceInfosCache[symbol][i]
			}
		}
	}

	return bMemory, cacheInfo
}

func (c *globalPriceInfoCache) GetPriceInfosByRange(symbol string, start, end int) conf.PriceInfosCache {
	tmpRetData := conf.PriceInfosCache{PriceInfosCache: map[string][]conf.PriceInfos{}}
	c.m.RLock()
	defer c.m.RUnlock()
	if start < len(c.gPriceInfosCache.PriceInfosCache[symbol]) {
		if end < len(c.gPriceInfosCache.PriceInfosCache[symbol]) {
			tmpRetData.PriceInfosCache[symbol] = c.gPriceInfosCache.PriceInfosCache[symbol][start:end]
		} else {
			tmpRetData.PriceInfosCache[symbol] = c.gPriceInfosCache.PriceInfosCache[symbol][start:]
		}
	}

	return tmpRetData
}

func (c *globalPriceInfoCache) GetCacheLength() int {
	c.m.Lock()
	defer c.m.Unlock()
	infoLen := len(c.gPriceInfosCache.PriceInfosCache)
	return infoLen
}

func (c *globalPriceInfoCache) GetSymbolCacheLength(symbol string) int {
	c.m.Lock()
	defer c.m.Unlock()
	infoLen := len(c.gPriceInfosCache.PriceInfosCache[symbol])
	return infoLen
}

func (c *globalPriceInfoCache) UpdateCachePrice(symbol string, infos conf.PriceInfos, maxMemTime int) {
	symbol = strings.Replace(symbol, "-", "", -1)
	c.m.Lock()
	defer c.m.Unlock()
	c.gPriceInfosCache.PriceInfosCache[symbol] = append(c.gPriceInfosCache.PriceInfosCache[symbol], infos)
	if len(c.gPriceInfosCache.PriceInfosCache[symbol]) > maxMemTime {
		c.gPriceInfosCache.PriceInfosCache[symbol] = c.gPriceInfosCache.PriceInfosCache[symbol][1:]
	}

}
