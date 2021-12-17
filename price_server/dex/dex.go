package dex

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	logger "github.com/sirupsen/logrus"
	"math/big"
	"price_api/price_server/dex/erc20"
	pair "price_api/price_server/dex/uniswapV2Pair"
	"price_api/price_server/util"
)

var debug = false

func GetUniswapAresPrice() (*big.Float, error) {
	logger.Info("get uniswap ares price")
	url := "https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	client, url := dialConn(url)

	aresVal, err := calAresEthPrice("0x7a646ee13eb104853c651e1d90d143acc9e72cdb", client)
	if err != nil {
		logger.WithError(err).Error("calAresEthPrice err")
		return nil, err
	}

	val, err := calEthPrice("0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852", client)
	if err != nil {
		logger.WithError(err).Error("calEthPrice err")
		return nil, err
	}
	result := new(big.Float).Mul(aresVal, val)
	return result, nil
}

func GetPancakeAresPrice() (*big.Float, error) {
	logger.Info("get pancake ares price")
	url := "https://bsc-dataseed1.ninicoin.io"
	client, url := dialConn(url)

	aresVal, err := calAresEthPrice("0x66e03400e47843ad396ee0a44dec403db8afeee0", client)
	if err != nil {
		logger.WithError(err).Error("calAresEthPrice err")
		return nil, err
	}

	val, err := calEthPrice("0x16b9a82891338f9ba80e2d6970fdda79d1eb0dae", client)
	if err != nil {
		logger.WithError(err).Error("calEthPrice err")
		return nil, err
	}
	result := new(big.Float).Mul(aresVal, val)
	return result, nil
}

func calAresEthPrice(pairAddr string, client *ethclient.Client) (*big.Float, error) {
	ethPairAddr := common.HexToAddress(pairAddr)

	eth, err := pair.NewPair(ethPairAddr, client)
	if err != nil {
		logger.WithError(err).Error("create new pair error")
	}
	res, err := eth.GetReserves(nil)
	if err != nil {
		logger.WithError(err).Error("GetReserves err")
	}
	//fmt.Println(time.Unix(int64(res.BlockTimestampLast), 0))

	token0, err := eth.Token0(nil)
	if err != nil {
		logger.WithError(err).Error("GetReserves token0 err")
	}
	weth, err := printErc20(token0, client)

	token1, err := eth.Token1(nil)
	if err != nil {
		logger.WithError(err).Error("GetReserves token1 err")
		return nil, err
	}

	usdt, err := printErc20(token1, client)
	if err != nil {
		logger.WithError(err).Error("printErc20 token1 err")
		return nil, err
	}
	if weth.Symbol == "ARES" {
		ethVal := util.ToDecimalsEth(res.Reserve0, weth.Decimals.Int64())
		usdtVal := util.ToDecimalsEth(res.Reserve1, usdt.Decimals.Int64())
		result := new(big.Float).Quo(usdtVal, ethVal)
		if debug {
			logger.Debugln("ares", ethVal, " usdtVal ", usdtVal, " result ", result)
		}
		return result, nil
	} else if usdt.Symbol == "ARES" {
		ethVal := util.ToDecimalsEth(res.Reserve1, weth.Decimals.Int64())
		usdtVal := util.ToDecimalsEth(res.Reserve0, usdt.Decimals.Int64())
		result := new(big.Float).Quo(usdtVal, ethVal)
		if debug {
			logger.Debugln("ares", ethVal, " usdtVal ", usdtVal, " result ", result)
		}
		return result, nil
	}

	return nil, errors.New("not find")
}

func calEthPrice(pairAddr string, client *ethclient.Client) (*big.Float, error) {
	ethPairAddr := common.HexToAddress(pairAddr)

	eth, err := pair.NewPair(ethPairAddr, client)
	if err != nil {
		logger.WithError(err).Error("create new pair error")
	}
	res, err := eth.GetReserves(nil)
	if err != nil {
		logger.WithError(err).Error("GetReserves err")
	}

	token0, err := eth.Token0(nil)
	if err != nil {
		logger.WithError(err).Error("GetReserves err")
	}
	weth, err := printErc20(token0, client)

	token1, err := eth.Token1(nil)
	if err != nil {
		logger.WithError(err).Error("GetReserves token1 err")
		return nil, err
	}

	usdt, err := printErc20(token1, client)
	if err != nil {
		logger.WithError(err).Error("printErc20 token1 err")
		return nil, err
	}

	if weth.Symbol == "WETH" {
		ethVal := util.ToDecimalsEth(res.Reserve0, weth.Decimals.Int64())
		usdtVal := util.ToDecimalsEth(res.Reserve1, usdt.Decimals.Int64())
		result := new(big.Float).Quo(usdtVal, ethVal)
		if debug {
			logger.Debugln("ethVal", ethVal, " usdtVal ", usdtVal, " result ", result)
		}
		return result, nil
	} else if usdt.Symbol == "WBNB" {
		ethVal := util.ToDecimalsEth(res.Reserve1, weth.Decimals.Int64())
		usdtVal := util.ToDecimalsEth(res.Reserve0, usdt.Decimals.Int64())
		result := new(big.Float).Quo(usdtVal, ethVal)
		if debug {
			logger.Debugln("ethVal", ethVal, " usdtVal ", usdtVal, " result ", result)
		}
		return result, nil
	}
	return nil, errors.New("not find")
}

type Erc20 struct {
	Name     string
	Symbol   string
	Decimals *big.Int
}

func printErc20(addr common.Address, client *ethclient.Client) (erc Erc20, err error) {

	ens, err := erc20.NewToken(addr, client)
	if err != nil {
		logger.WithError(err).Error("can't NewContract")
	}

	// Set ourself as the owner of the name.
	name, err := ens.Name(nil)
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve token name")
	}
	erc.Name = name

	// Set ourself as the owner of the name.
	symbol, err := ens.Symbol(nil)
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve token symbol")
	}
	erc.Symbol = symbol

	decimals, err := ens.Decimals(nil)
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve token decimals")
	}
	if debug {
		logger.Debugln("addr", addr, " Token name:", name, " Token symbol:", symbol, " Token decimals:", decimals)
	}
	erc.Decimals = decimals
	return erc, err
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
		logger.WithError(err).Errorln("dialConn", "Failed to connect to the ethereum client")
	}
	return conn, url
}
