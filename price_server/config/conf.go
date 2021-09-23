package conf

import (
	"errors"
	"strconv"

	"github.com/pelletier/go-toml"
)

type MysqlConfig struct {
	Server   string
	Port     int64
	Db       string
	Name     string
	Password string
}

type ExchangeConfig struct {
	Name   string
	Weight int64
	Url    string
}

type Config struct {
	Interval       int64
	Port           int64
	Proxy          string
	InsertInterval int64 `toml:"insertInterval"`
	MaxVolume      int64 `toml:"maxVolume"`
	PageSize       int64 `toml:"pageSize"`
	Mysql          MysqlConfig
	Exchanges      []ExchangeConfig
	Symbols        []string
}

// func GetConfig() (Config, error) {
// 	var config Config
// 	if _, err := toml.DecodeFile("./conf.toml", &config); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(config.Symbol)
// 	return config, nil
// }

func GetConfig() (Config, error) {
	var retConfig Config
	var ok bool

	config, err := toml.LoadFile("./conf.toml")
	if err != nil {
		return Config{}, err
	}

	retConfig.Interval, ok = config.Get("interval").(int64)
	if !ok {
		return Config{}, errors.New("parse key interval error")
	}

	retConfig.Port, ok = config.Get("port").(int64)
	if !ok {
		return Config{}, errors.New("parse key port error")
	}

	retConfig.MaxVolume, ok = config.Get("maxVolume").(int64)
	if !ok {
		return Config{}, errors.New("parse key maxVolume error")
	}

	retConfig.InsertInterval, ok = config.Get("insertInterval").(int64)
	if !ok {
		return Config{}, errors.New("parse key insertInterval error")
	}

	retConfig.PageSize, ok = config.Get("pageSize").(int64)
	if !ok {
		return Config{}, errors.New("parse key pageSize error")
	}

	retConfig.Proxy, ok = config.Get("proxy").(string)
	if !ok {
		return Config{}, errors.New("parse key proxy error")
	}

	sysbols := config.Get("symbols")
	if sysbols == nil {
		return Config{}, errors.New("get symbols error")
	}

	for _, symbolIdx := range sysbols.([]interface{}) {
		retConfig.Symbols = append(retConfig.Symbols, symbolIdx.(string))
	}

	retConfig.Mysql.Server, ok = config.Get("mysql.server").(string)
	if !ok {
		return Config{}, errors.New("parse key mysql.server error")
	}

	retConfig.Mysql.Port, ok = config.Get("mysql.port").(int64)
	if !ok {
		return Config{}, errors.New("parse key mysql.port error")
	}

	retConfig.Mysql.Db, ok = config.Get("mysql.db").(string)
	if !ok {
		return Config{}, errors.New("parse key mysql.db error")
	}

	retConfig.Mysql.Name, ok = config.Get("mysql.name").(string)
	if !ok {
		return Config{}, errors.New("parse key mysql.name error")
	}

	retConfig.Mysql.Password, ok = config.Get("mysql.password").(string)
	if !ok {
		return Config{}, errors.New("parse key mysql.password error")
	}

	index := 1
	exchange := "exchange."
	for {
		session := exchange + strconv.Itoa(index)
		if !config.Has(session) {
			break
		}

		var exchangeConfig ExchangeConfig
		exchangeConfig.Name, ok = config.Get(session + ".name").(string)
		if !ok {
			return Config{}, errors.New("parse key " + session + ".name error")
		}

		exchangeConfig.Url, ok = config.Get(session + ".url").(string)
		if !ok {
			return Config{}, errors.New("parse key " + session + ".url error")
		}

		exchangeConfig.Weight, ok = config.Get(session + ".weight").(int64)
		if !ok {
			return Config{}, errors.New("parse key " + session + ".weight error")
		}
		retConfig.Exchanges = append(retConfig.Exchanges, exchangeConfig)
		index++
	}

	return retConfig, nil
}
