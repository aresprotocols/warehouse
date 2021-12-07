package util

import "math/big"

func ToEth(val *big.Int) *big.Float {
	baseUnit := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	fbaseUnit := new(big.Float).SetFloat64(float64(baseUnit.Int64()))
	return new(big.Float).Quo(new(big.Float).SetInt(val), fbaseUnit)
}

func ToDecimalsEth(val *big.Int, decimals int64) *big.Float {
	baseUnit := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	fbaseUnit := new(big.Float).SetFloat64(float64(baseUnit.Int64()))
	return new(big.Float).Quo(new(big.Float).SetInt(val), fbaseUnit)
}
