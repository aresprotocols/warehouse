package cache

import (
	conf "price_api/price_server/config"
	"sync"
)

type GlobalPriceInfoCache struct {
	gPriceInfosCache conf.PriceInfosCache
	m                *sync.RWMutex
}

func NewGlobalPriceInfoCache() *GlobalPriceInfoCache {
	return &GlobalPriceInfoCache{
		gPriceInfosCache: conf.PriceInfosCache{},
		m:                new(sync.RWMutex)}
}

func (c *GlobalPriceInfoCache) GetLatestPriceInfos() conf.PriceInfos {
	c.m.RLock()
	latestInfos := c.gPriceInfosCache.PriceInfosCache[len(c.gPriceInfosCache.PriceInfosCache)-1]
	c.m.RUnlock()
	return latestInfos
}

func (c *GlobalPriceInfoCache) GetPriceInfosEqualTimestamp(timestamp int64) (bool, conf.PriceInfos) {
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

func (c *GlobalPriceInfoCache) GetPriceInfosByRange(start, end int) conf.PriceInfosCache {
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

func (c *GlobalPriceInfoCache) UpdateSymbolWeight(symbol, exchange string, weight int) {
	c.m.Lock()
	for i, confTemp := range conf.GRequestPriceConfs[symbol] {
		if confTemp.Name == exchange {
			conf.GRequestPriceConfs[symbol][i].Weight = int64(weight)
			break
		}
	}
	c.m.Unlock()
}

func (c *GlobalPriceInfoCache) GetCacheLength() int {
	c.m.RLock()
	infoLen := len(c.gPriceInfosCache.PriceInfosCache)
	c.m.RUnlock()
	return infoLen
}

func (c *GlobalPriceInfoCache) UpdateCachePrice(infos conf.PriceInfos, maxMemTime int) {
	c.m.Lock()
	c.gPriceInfosCache.PriceInfosCache = append(c.gPriceInfosCache.PriceInfosCache, infos)
	if len(c.gPriceInfosCache.PriceInfosCache) > maxMemTime {
		c.gPriceInfosCache.PriceInfosCache = c.gPriceInfosCache.PriceInfosCache[1:]
	}
	c.m.Unlock()
}
