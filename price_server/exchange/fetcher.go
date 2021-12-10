package exchange

import (
	"errors"
	"math/big"
	conf "price_api/price_server/config"
	"price_api/price_server/dex"
	"sync"
	"time"
)

var (
	fetcher       *Fetcher
	errTerminated = errors.New("terminated")
)

type DexPrice struct {
	Uni     *BasePrice
	Pancake *BasePrice
}

type BasePrice struct {
	Price     *big.Float
	TimeStamp int64
}

type announce struct {
	Price     *big.Float
	TimeStamp int64
	UniSwap   bool
	err       error
}

type announceCMC struct {
	CASI AresShowInfo
	err  error
}

type Fetcher struct {
	quit       chan struct{}
	blockMutex *sync.Mutex //block mutex
	Price      *DexPrice
	notify     chan *announce
	proxy      string
	CASI       AresShowInfo
}

func InitFetcher(cfg conf.Config, proxy string) *Fetcher {
	fetcher = &Fetcher{
		notify:     make(chan *announce),
		Price:      &DexPrice{},
		blockMutex: new(sync.Mutex),
		quit:       make(chan struct{}),
		proxy:      proxy,
	}
	return fetcher
}

func (f *Fetcher) GetDexPrice() *DexPrice {
	f.blockMutex.Lock()
	defer f.blockMutex.Unlock()

	return f.Price
}

// Start boots up the announcement based synchroniser, accepting and processing
// hash notifications and block fetches until termination requested.
func (f *Fetcher) Start() {
	go f.calUniswapPrice(true)
	go f.calUniswapPrice(false)
	go f.loop()
}

// Stop terminates the announcement based synchroniser, canceling all pending
// operations.
func (f *Fetcher) Stop() {
	close(f.quit)
}

// Loop is the main fetcher loop, checking and processing various notification
// events.
func (f *Fetcher) loop() {
	// Iterate the block fetching until a quit is requested

	for {
		// Clean up any expired block fetches
		// Import any queued blocks that could potentially fit
		// Wait for an outside event to occur
		select {
		case <-f.quit:
			// Fetcher terminating, abort all operations
			return

			// At least one block's timer ran out, check for needing retrieval
		case notification := <-f.notify:
			if notification.err == nil {
				base := &BasePrice{
					Price:     new(big.Float).Set(notification.Price),
					TimeStamp: notification.TimeStamp,
				}

				if notification.UniSwap {
					f.Price.Uni = base
				} else {
					f.Price.Uni = base
					f.Price.Pancake = base
				}
			}
		}
	}
}

func (f *Fetcher) calUniswapPrice(uniswap bool) {
	timer1 := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer1.C:
			var ann = &announce{
				UniSwap: uniswap,
			}
			var value *big.Float
			var err error
			for i := 0; i < 2; i++ {
				if uniswap {
					value, err = dex.GetUniswapAresPrice()
				} else {
					value, err = dex.GetPancakeAresPrice()
				}
				if value != nil {
					ann.Price = value
					ann.TimeStamp = time.Now().Unix()
					f.notify <- ann
					break
				}
				if i == 1 && err != nil {
					ann.err = err
					ann.TimeStamp = time.Now().Unix()
					f.notify <- ann
				}
			}
		case <-f.quit:
			return
		}
	}
}

func (f *Fetcher) calGateCMCAresInfo() {
	timer1 := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-timer1.C:
			GetAresInfo(f.proxy)
		case <-f.quit:
			return
		}
	}
}
