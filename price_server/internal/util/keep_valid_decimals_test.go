package util

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func TestKeepValidDecimals(t *testing.T) {

	keep := 6
	tests := []struct {
		args float64
		want float64
	}{
		{
			args: 0,
			want: 0,
		},
		{
			args: 1,
			want: 1,
		},
		{
			args: 10,
			want: 10,
		},
		{
			args: 11,
			want: 11,
		},
		{
			args: 58909,
			want: 58909,
		},
		{
			args: 658900,
			want: 658900,
		},
		{
			args: 658901,
			want: 658901,
		},
		{
			args: 658901.01,
			want: 658901.01,
		},
		{
			args: 658901.1,
			want: 658901.1,
		},
		{
			args: 658901.001234566,
			want: 658901.001235,
		},
		{
			args: 658901.000234466,
			want: 658901.000234,
		},
		{
			args: 0.1,
			want: 0.1,
		},
		{
			args: 0.01,
			want: 0.01,
		},
		{
			args: 0.0001,
			want: 0.0001,
		},
		{
			args: 0.01234556,
			want: 0.0123456,
		},
		{
			args: 0.01234454,
			want: 0.0123445,
		},
		{
			args: 0.000000012345,
			want: 0.000000012345,
		},
		{
			args: 0.0012344,
			want: 0.0012344,
		},
		{
			args: 0.000000012345789,
			want: 0.0000000123458,
		},
	}
	for _, tt := range tests {
		dec := decimal.NewFromFloat(tt.args)
		t.Run(fmt.Sprintf("test %s", dec.String()), func(t *testing.T) {
			if got := KeepValidDecimals(dec, keep); got != tt.want {
				t.Errorf("KeepValidDecimals() = %v, want %v", got, tt.want)
			}
		})
	}
}
