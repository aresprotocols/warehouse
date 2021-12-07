package dex

import (
	"context"
	"fmt"
	"math/big"
	"testing"
)

func TestUniswapPair(t *testing.T) {
	url := "https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	client, url := dialConn(url)

	aresVal, err := calAresEthPrice("0x7a646ee13eb104853c651e1d90d143acc9e72cdb", client)
	if err != nil {
		fmt.Println("GetReserves err", err)
	}
	fmt.Println("aresVal", aresVal)

	val, err := calEthPrice("0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852", client)
	if err != nil {
		fmt.Println("GetReserves err", err)
	}
	fmt.Println("val", val)
	result := new(big.Float).Mul(aresVal, val)
	fmt.Println(" ", result, " ", result.Text('g', 4))
}

func TestPanckeswapPair(t *testing.T) {
	url := "https://bsc-dataseed1.ninicoin.io"
	client, url := dialConn(url)
	num, _ := client.BlockNumber(context.Background())
	fmt.Println("num", num)

	aresVal, err := calAresEthPrice("0x66e03400e47843ad396ee0a44dec403db8afeee0", client)
	if err != nil {
		fmt.Println("GetReserves err", err)
	}
	fmt.Println("aresVal", aresVal)

	val, err := calEthPrice("0x16b9a82891338f9ba80e2d6970fdda79d1eb0dae", client)
	if err != nil {
		fmt.Println("GetReserves err", err)
	}
	fmt.Println("val", val)
	result := new(big.Float).Mul(aresVal, val)
	fmt.Println(" ", result, " ", result.Text('g', 4))
}
