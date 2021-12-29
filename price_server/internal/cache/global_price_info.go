package cache

import (
	conf "price_api/price_server/config"
	"sync"
)

type GlobalPriceInfoCache interface {
	GetLatestPriceInfos() conf.PriceInfos
	GetPriceInfosEqualTimestamp(timestamp int64) (bool, conf.PriceInfos)
	GetPriceInfosByRange(start, end int) conf.PriceInfosCache
	UpdateSymbolWeight(symbol, exchange string, weight int)
	GetCacheLength() int
	UpdateCachePrice(infos conf.PriceInfos, maxMemTime int)
}

type globalPriceInfoCache struct {
	gPriceInfosCache conf.PriceInfosCache
	m                *sync.RWMutex
}

func NewGlobalPriceInfoCache() GlobalPriceInfoCache {
	return &globalPriceInfoCache{
		gPriceInfosCache: conf.PriceInfosCache{},
		m:                new(sync.RWMutex)}
}

func (c *globalPriceInfoCache) GetLatestPriceInfos() conf.PriceInfos {
	c.m.RLock()
	if len(c.gPriceInfosCache.PriceInfosCache) == 0 {
		return conf.PriceInfos{}
	}
	latestInfos := c.gPriceInfosCache.PriceInfosCache[len(c.gPriceInfosCache.PriceInfosCache)-1]
	c.m.RUnlock()
	return latestInfos
}

func (c *globalPriceInfoCache) GetPriceInfosEqualTimestamp(timestamp int64) (bool, conf.PriceInfos) {
	bMemory := false
	var cacheInfo conf.PriceInfos
	c.m.RLock()
	//latestInfos := gPriceInfosCache.PriceInfosCache[len(gPriceInfosCache.PriceInfosCache)-1]
	if len(c.gPriceInfosCache.PriceInfosCache) != 0 {
		for i := len(c.gPriceInfosCache.PriceInfosCache) - 1; i >= 0; i-- {
			info := c.gPriceInfosCache.PriceInfosCache[i]
			if len(info.PriceInfos) == 0 {
				continue
			}
			if info.PriceInfos[0].TimeStamp == timestamp {
				//use memory
				bMemory = true
				cacheInfo = c.gPriceInfosCache.PriceInfosCache[i]
			}
		}
	}
	c.m.RUnlock()
	return bMemory, cacheInfo
}

func (c *globalPriceInfoCache) GetPriceInfosByRange(start, end int) conf.PriceInfosCache {
	tmpRetData := conf.PriceInfosCache{}
	c.m.RLock()
	if start < len(c.gPriceInfosCache.PriceInfosCache) {
		if end < len(c.gPriceInfosCache.PriceInfosCache) {
			tmpRetData.PriceInfosCache = c.gPriceInfosCache.PriceInfosCache[start:end]
		} else {
			tmpRetData.PriceInfosCache = c.gPriceInfosCache.PriceInfosCache[start:]
		}
	}
	c.m.RUnlock()
	return tmpRetData
}

//todo add unit test
func (c *globalPriceInfoCache) UpdateSymbolWeight(symbol, exchange string, weight int) {
	c.m.Lock()
	for i, confTemp := range conf.GRequestPriceConfs[symbol] {
		if confTemp.Name == exchange {
			conf.GRequestPriceConfs[symbol][i].Weight = int64(weight)
			break
		}
	}
	c.m.Unlock()
}

func (c *globalPriceInfoCache) GetCacheLength() int {
	c.m.RLock()
	infoLen := len(c.gPriceInfosCache.PriceInfosCache)
	c.m.RUnlock()
	return infoLen
}

func (c *globalPriceInfoCache) UpdateCachePrice(infos conf.PriceInfos, maxMemTime int) {
	c.m.Lock()
	c.gPriceInfosCache.PriceInfosCache = append(c.gPriceInfosCache.PriceInfosCache, infos)
	if len(c.gPriceInfosCache.PriceInfosCache) > maxMemTime {
		c.gPriceInfosCache.PriceInfosCache = c.gPriceInfosCache.PriceInfosCache[1:]
	}
	c.m.Unlock()
}
