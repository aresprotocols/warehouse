package dex

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"price_api/price_server/dex/erc20"
	pair "price_api/price_server/dex/uniswapV2Pair"
	"price_api/price_server/util"
	"time"
)

func calAresEthPrice(pairAddr string, client *ethclient.Client) (*big.Float, error) {
	ethPairAddr := common.HexToAddress(pairAddr)

	eth, err := pair.NewPair(ethPairAddr, client)
	if err != nil {
		fmt.Println("err", err)
	}
	res, err := eth.GetReserves(nil)
	if err != nil {
		fmt.Println("GetReserves err", err)
	}
	fmt.Println(time.Unix(int64(res.BlockTimestampLast), 0))

	token0, err := eth.Token0(nil)
	if err != nil {
		fmt.Println("GetReserves err", err)
	}
	weth := printErc20(token0, client)

	token1, err := eth.Token1(nil)
	if err != nil {
		fmt.Println("GetReserves err", err)
		return nil, err
	}

	usdt := printErc20(token1, client)
	if weth.Symbol == "ARES" {
		ethVal := util.ToDecimalsEth(res.Reserve0, weth.Decimals.Int64())
		usdtVal := util.ToDecimalsEth(res.Reserve1, usdt.Decimals.Int64())
		result := new(big.Float).Quo(usdtVal, ethVal)
		fmt.Println("ares", ethVal, " usdtVal ", usdtVal, " result ", result)
		return result, nil
	} else if usdt.Symbol == "ARES" {
		ethVal := util.ToDecimalsEth(res.Reserve1, weth.Decimals.Int64())
		usdtVal := util.ToDecimalsEth(res.Reserve0, usdt.Decimals.Int64())
		result := new(big.Float).Quo(usdtVal, ethVal)
		fmt.Println("ares", ethVal, " usdtVal ", usdtVal, " result ", result)
		return result, nil
	}

	return nil, errors.New("not find")
}

func calEthPrice(pairAddr string, client *ethclient.Client) (*big.Float, error) {
	ethPairAddr := common.HexToAddress(pairAddr)

	eth, err := pair.NewPair(ethPairAddr, client)
	if err != nil {
		fmt.Println("err", err)
	}
	res, err := eth.GetReserves(nil)
	if err != nil {
		fmt.Println("GetReserves err", err)
	}

	token0, err := eth.Token0(nil)
	if err != nil {
		fmt.Println("GetReserves err", err)
	}
	weth := printErc20(token0, client)

	token1, err := eth.Token1(nil)
	if err != nil {
		fmt.Println("GetReserves err", err)
		return nil, err
	}

	usdt := printErc20(token1, client)
	if weth.Symbol == "WETH" {
		ethVal := util.ToDecimalsEth(res.Reserve0, weth.Decimals.Int64())
		usdtVal := util.ToDecimalsEth(res.Reserve1, usdt.Decimals.Int64())
		result := new(big.Float).Quo(usdtVal, ethVal)
		fmt.Println("ethVal", ethVal, " usdtVal ", usdtVal, " result ", result)
		return result, nil
	} else if usdt.Symbol == "WBNB" {
		ethVal := util.ToDecimalsEth(res.Reserve1, weth.Decimals.Int64())
		usdtVal := util.ToDecimalsEth(res.Reserve0, usdt.Decimals.Int64())
		result := new(big.Float).Quo(usdtVal, ethVal)
		fmt.Println("ethVal", ethVal, " usdtVal ", usdtVal, " result ", result)
		return result, nil
	}
	return nil, errors.New("not find")
}

type Erc20 struct {
	Name     string
	Symbol   string
	Decimals *big.Int
}

func printErc20(addr common.Address, client *ethclient.Client) (erc Erc20) {

	ens, err := erc20.NewToken(addr, client)
	if err != nil {
		fmt.Printf("can't NewContract: %v\n", err)
	}

	// Set ourself as the owner of the name.
	name, err := ens.Name(nil)
	if err != nil {
		fmt.Println("Failed to retrieve token ", "name: %v", err)
	}
	erc.Name = name

	// Set ourself as the owner of the name.
	symbol, err := ens.Symbol(nil)
	if err != nil {
		fmt.Println("Failed to retrieve token ", "name: %v", err)
	}
	erc.Symbol = symbol

	decimals, err := ens.Decimals(nil)
	if err != nil {
		fmt.Println("Failed to retrieve token ", "name: %v", err)
	}
	fmt.Println("addr", addr, " Token name:", name, " Token symbol:", symbol, " Token decimals:", decimals)
	erc.Decimals = decimals
	return erc
}

func dialConn(url string) (*ethclient.Client, string) {
	ip := "165.227.99.131"
	port := 8545

	//url = "https://ethrpc.truescan.network"
	//url = "https://kovan.poa.network/"

	if url == "" {
		url = fmt.Sprintf("http://%s", fmt.Sprintf("%s:%d", ip, port))
	}
	// Create an IPC based RPC connection to a remote node
	// "http://39.100.97.129:8545"
	conn, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal("dialConn", "Failed to connect to the ethereum client: %v", err)
	}
	return conn, url
}
