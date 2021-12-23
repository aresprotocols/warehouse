package main

import (
	logger "github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	conf "price_api/price_server/config"
	"price_api/price_server/exchange"
	"price_api/price_server/internal/handler"
	"price_api/price_server/internal/repository"
	"price_api/price_server/internal/routers"
	"price_api/price_server/internal/service"
	"strconv"
	"time"
)

func init() {
	config := DefaultConfiguration()
	err := InitLogrusLogger(config)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}

}

func main() {
	//gRequestPriceConfs = make(map[string][]conf.ExchangeConfig)

	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	cfg, err := conf.GetConfig()
	if err != nil {
		logger.Errorf("get config occur err:%v", err)
		return
	}

	conf.GCfg = cfg
	logger.Infof("config load over:%v", cfg)

	err = repository.InitMysqlDB(cfg)
	if err != nil {
		logger.Errorf("Init mysql db occur err:%v", err)
		return
	}

	logger.Info("mysql init over")
	// Init service
	service.Svc = service.New(repository.DB)

	handle := handler.InitHandle(cfg)

	conf.GRequestPriceConfs, err = exchange.InitRequestPriceConf(cfg)
	if err != nil {
		logger.Errorf("Init request price conf occur err:%v", err)
		return
	}
	logger.Info("request init over")

	showIgnoreSymbols(cfg, conf.GRequestPriceConfs)

	router := routers.NewRouter()

	go updatePrice(cfg, conf.GRequestPriceConfs)
	router.Run(":" + strconv.Itoa(int(cfg.Port)))

	abortChan := make(chan os.Signal, 1)
	signal.Notify(abortChan, os.Interrupt)

	sig := <-abortChan
	handle.Stop()
	logger.Infof("Exiting... signal %v", sig)
}

func updatePrice(cfg conf.Config, reqConf map[string][]conf.ExchangeConfig) {
	idx := 0
	time.Sleep(time.Second * 2) // run update for the first time,  need to sleep , because you have just completed initialization and have already requested data once
	priceService := service.Svc.Price()
	for {
		logger.Infof("start new round update price")
		infos, err := exchange.GetExchangePrice(reqConf, cfg)
		if err != nil {
			logger.WithError(err).Errorf("get exchange price occur error")
		} else {
			idx++
			priceService.UpdateCachePrice(infos, int(cfg.MaxMemTime))
		}

		if idx >= int(cfg.InsertInterval) {
			err = priceService.InsertPriceInfo(infos)
			if err != nil {
				logger.Errorf("insert price info occur err:%v", err)
			} else {
				idx = 0
			}
		}
		logger.Infof("end this round update price")
		time.Sleep(time.Second * time.Duration(cfg.Interval))

	}
}

func showIgnoreSymbols(cfg conf.Config, gRequestPriceConfs map[string][]conf.ExchangeConfig) {
	ignoreSymbols := make(map[string][]string)
	for _, symbol := range cfg.Symbols {
		var exchanges []string
		existSymbols, ok := gRequestPriceConfs[symbol]
		if ok {
			for _, exchangeConf := range cfg.Exchanges {
				//check config exchange if have symbol
				bFind := false
				for _, existSymbol := range existSymbols {
					if exchangeConf.Name == existSymbol.Name {
						//find it
						bFind = true
					}
				}
				if !bFind {
					exchanges = append(exchanges, exchangeConf.Name)
				}
			}
		} else {
			for _, exchangeConf := range cfg.Exchanges {
				exchanges = append(exchanges, exchangeConf.Name)
			}
		}
		ignoreSymbols[symbol] = exchanges
	}
	logger.Infoln("ignore symbols and exchange:", ignoreSymbols)
}
