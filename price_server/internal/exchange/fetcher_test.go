package exchange

import (
	"fmt"
	"price_api/price_server/internal/config"
	"testing"
	"time"
)

func TestFetcher(t *testing.T) {
	fet := InitFetcher(conf.Config{})
	fet.Start()
	time.Sleep(40 * time.Second)
	fmt.Println(fet.GetDexPrice())
	fet.Stop()
}

func TestFetcherCMC(t *testing.T) {
	fet := InitFetcher(conf.Config{})
	fet.Start()
	time.Sleep(40 * time.Second)
	fmt.Println(fet.GetCMCInfo())
	fet.Stop()
}
