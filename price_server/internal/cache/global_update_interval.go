package cache

import (
	"sync"
)

type GlobalUpdateIntervalCache interface {
	UpdateSymbolInterval(symbol string, interval int)
	GetSymbolInterval(symbol string) int
}

func NewGlobalUpdateIntervalCache() GlobalUpdateIntervalCache {
	return &globalUpdateIntervalCache{
		updateIntervalCache: map[string]int{},
		m:                   new(sync.RWMutex)}
}

type globalUpdateIntervalCache struct {
	updateIntervalCache map[string]int
	m                   *sync.RWMutex
}

func (g *globalUpdateIntervalCache) UpdateSymbolInterval(symbol string, interval int) {
	g.m.Lock()
	defer g.m.Unlock()
	g.updateIntervalCache[symbol] = interval
}

func (g *globalUpdateIntervalCache) GetSymbolInterval(symbol string) int {
	g.m.RLock()
	defer g.m.RUnlock()
	return g.updateIntervalCache[symbol]
}
