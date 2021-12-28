package util

import (
	"github.com/shopspring/decimal"
	conf "price_api/price_server/config"
	"price_api/price_server/internal/vo"
	"sort"
	"strings"
)

//@param bAverage     get average not cointain lowest and highest
//@return bool     symbol find?
func PartyPrice(infos []conf.PriceInfo, symbol string, bAverage bool) (bool, vo.PartyPriceInfo) {
	var symbolPriceInfo []conf.PriceInfo
	for _, info := range infos {
		if strings.EqualFold(info.Symbol, symbol) {
			symbolPriceInfo = append(symbolPriceInfo, info)
		}
	}

	infosLen := len(symbolPriceInfo)
	if infosLen == 0 {
		return false, vo.PartyPriceInfo{}
	}

	sort.Slice(symbolPriceInfo, func(i, j int) bool {
		if symbolPriceInfo[i].Price > symbolPriceInfo[j].Price {
			return true
		} else {
			return false
		}
	})

	if infosLen > 2 && bAverage {
		symbolPriceInfo = symbolPriceInfo[1 : infosLen-1]
	}

	var partyPriceInfo vo.PartyPriceInfo
	totalPrice := decimal.NewFromFloat(0)
	totalWeight := decimal.NewFromFloat(0)
	for _, info := range symbolPriceInfo {
		//totalPrice += info.Price * float64(info.Weight)
		totalPrice = totalPrice.Add(decimal.NewFromFloat(info.Price).Mul(decimal.NewFromInt(info.Weight)))
		//totalWeight += info.Weight
		totalWeight = totalWeight.Add(decimal.NewFromInt(info.Weight))

		partyPriceInfo.Infos = append(partyPriceInfo.Infos, vo.WeightInfo{Price: info.Price, Weight: info.Weight, ExchangeName: info.PriceOrigin})
	}
	partyPriceInfo.Price = KeepValidDecimals(totalPrice.Div(totalWeight), 6)
	partyPriceInfo.Timestamp = symbolPriceInfo[0].TimeStamp

	return true, partyPriceInfo
}
