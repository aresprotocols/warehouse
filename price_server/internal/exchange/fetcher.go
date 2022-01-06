package exchange

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"math/big"
	"price_api/price_server/internal/config"
	"price_api/price_server/internal/dex"
	"sync"
	"time"
)

var (
	fetcher       *Fetcher
	errTerminated = errors.New("terminated")
)

type DexPrice struct {
	Uni     *BasePrice `json:"uni"`
	Pancake *BasePrice `json:"pancake"`
}

type BasePrice struct {
	Price     string `json:"price"`
	TimeStamp int64  `json:"timestamp"`
}

type announce struct {
	Price     *big.Float
	TimeStamp int64
	UniSwap   bool
	err       error
}

type announceCMC struct {
	CASI      AresShowInfo
	TimeStamp int64
	err       error
}

type Fetcher struct {
	quit       chan struct{}
	blockMutex *sync.Mutex //block mutex
	price      *DexPrice
	notify     chan *announce
	notifyCMC  chan *announceCMC
	casi       AresShowInfo
	cfg        conf.Config
}

func InitFetcher(cfg conf.Config) *Fetcher {
	fetcher = &Fetcher{
		notify:     make(chan *announce),
		notifyCMC:  make(chan *announceCMC),
		price:      &DexPrice{},
		blockMutex: new(sync.Mutex),
		quit:       make(chan struct{}),
		cfg:        cfg,
	}
	return fetcher
}

func (f *Fetcher) GetDexPrice() *DexPrice {
	f.blockMutex.Lock()
	defer f.blockMutex.Unlock()

	return f.price
}

func (f *Fetcher) GetCMCInfo() *AresShowInfo {
	f.blockMutex.Lock()
	defer f.blockMutex.Unlock()

	return &f.casi
}

// Start boots up the announcement based synchroniser, accepting and processing
// hash notifications and block fetches until termination requested.
func (f *Fetcher) Start() {
	go f.calUniswapPriceTimer(true)
	go f.calUniswapPriceTimer(false)
	go f.calGateCMCAresInfoTimer()
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
					Price:     notification.Price.String(),
					TimeStamp: notification.TimeStamp,
				}

				if notification.UniSwap {
					f.price.Uni = base
				} else {
					f.price.Pancake = base
				}
			}
		case notification := <-f.notifyCMC:
			if notification.err == nil {
				f.casi = notification.CASI
			}
		}
	}
}

func (f *Fetcher) calUniswapPriceTimer(uniswap bool) {
	go f.calUniswapPrice(uniswap)

	timer1 := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-timer1.C:
			go f.calUniswapPrice(uniswap)
		case <-f.quit:
			return
		}
	}
}

func (f *Fetcher) calUniswapPrice(uniswap bool) {
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
}

func (f *Fetcher) calGateCMCAresInfoTimer() {
	go f.calGateCMCAresInfo()

	timer1 := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-timer1.C:
			go f.calGateCMCAresInfo()
		case <-f.quit:
			return
		}
	}
}

func (f *Fetcher) calGateCMCAresInfo() {
	for i := 0; i < 2; i++ {
		var ann = &announceCMC{}
		logger.Info("cal gate CMC ares info")
		info, err := getCMCAresInfo(f.cfg.Proxy)
		if err == nil {
			ann.CASI = info
			ann.TimeStamp = time.Now().Unix()
			f.notifyCMC <- ann
			break
		}
		if i == 1 {
			ann.err = err
			ann.TimeStamp = time.Now().Unix()
			f.notifyCMC <- ann
		}
	}
}
