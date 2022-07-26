package cron

import (
	"github.com/robfig/cron/v3"
	logger "github.com/sirupsen/logrus"
	"price_api/price_server/internal/service"
	"time"
)

func StartCron() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.WithError(err).Errorf("load lcoation Asia/Shanghai occur error")
		return
	}
	requestInfoService := service.Svc.RequestInfo()
	coinHistoryService := service.Svc.CoinHistory()
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), cron.WithLocation(location))
	_, _ = c.AddFunc("0 0 */1 * * *", func() {
		requestInfoService.DeleteOld()
		coinHistoryService.DeleteOld()
	})
	c.Start()
}
