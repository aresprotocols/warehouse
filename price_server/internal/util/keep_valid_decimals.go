package util

import (
	"github.com/shopspring/decimal"
	"strconv"
)

func KeepValidDecimals(num decimal.Decimal, keep int) float64 {
	if num.IsZero() {
		return 0
	}
	if num.GreaterThanOrEqual(decimal.NewFromInt(1)) {
		floatNum, _ := strconv.ParseFloat(num.StringFixed(int32(keep)), 64)
		return floatNum
	} else {
		bigNum := num.Mul(decimal.New(1, 18))
		length := len(bigNum.String())
		bigNum = bigNum.Round(0 - int32(length-keep))
		bigNum = bigNum.Div(decimal.New(1, 18))
		floatNum, _ := strconv.ParseFloat(bigNum.String(), 64)
		return floatNum
	}
}
