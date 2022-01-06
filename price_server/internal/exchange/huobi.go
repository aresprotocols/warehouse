package exchange

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

//{"ch":"market.btcusdt.detail.merged","status":"ok","ts":1629098174462,"tick":{"id":271006617773,"version":271006617773,"open":46266.36,"close":47323.59,"low":45477.01,"high":48050.0,"amount":17729.513874233984,"vol":8.266551922141494E8,"count":549488,"bid":[47323.58,1.55321],"ask":[47323.59,0.06133773059609383]}}
type HuobiPriceInfo struct {
	Status   string   `json:"status"`
	Tick     TickInfo `json:"tick"`
	ErrorMsg string   `json:"err-msg"`
}

type TickInfo struct {
	Ask []float64 `json:"ask"`
}

func parseHuobiPrice(priceJson string) (float64, error) {
	var huobiPriceInfo HuobiPriceInfo

	err := json.Unmarshal([]byte(priceJson), &huobiPriceInfo)
	if err != nil {
		return 0, err
	}

	if huobiPriceInfo.Status == "error" {
		if huobiPriceInfo.ErrorMsg == "invalid symbol" {
			return 0, nil
		} else {
			return 0, xerrors.New("some error")
		}
	} else {
		if len(huobiPriceInfo.Tick.Ask) == 0 {
			logger.Infoln("response:", huobiPriceInfo.Status, " Tick:", huobiPriceInfo.Tick)

		}
		return huobiPriceInfo.Tick.Ask[0], nil
	}
}
