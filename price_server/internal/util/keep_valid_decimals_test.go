package util

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func TestKeepValidDecimals(t *testing.T) {
	keep := 6

	args := []float64{0, 1, 10, 11, 58909, 658900, 658901, 658901.01, 658901.1, 658901.001234566, 658901.000234466, 0.1, 0.01, 0.0001, 0.01234556, 0.01234454, 0.000000012345, 0.0012344, 0.000000012345789}
	wants := []float64{0, 1, 10, 11, 58909, 658900, 658901, 658901.01, 658901.1, 658901.001235, 658901.000234, 0.1, 0.01, 0.0001, 0.0123456, 0.0123445, 0.000000012345, 0.0012344, 0.0000000123458}

	for i, tt := range args {
		dec := decimal.NewFromFloat(tt)
		t.Run(fmt.Sprintf("test %s", dec.String()), func(t *testing.T) {
			if got := KeepValidDecimals(dec, keep); got != wants[i] {
				t.Errorf("KeepValidDecimals() = %v, want %v", got, wants[i])
			}
		})
	}
}
