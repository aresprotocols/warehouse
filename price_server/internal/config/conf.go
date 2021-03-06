package conf

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml"
)

var GCfg Config

//var GRequestPriceConfs map[string][]ExchangeConfig

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
	MaxMemTime     int64 `toml:"maxMemTime"`
	PageSize       int64 `toml:"pageSize"`
	RetryCount     int64 `toml:"retryCount"`
	User           string
	Password       string
	RunByDocker    bool
	Mysql          MysqlConfig
	Exchanges      []ExchangeConfig
	Symbols        []string
	SymbolReplaces map[string]map[string]string
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

	config, err := toml.LoadFile("./configs/conf.toml")
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

	retConfig.RunByDocker, ok = config.Get("runByDocker").(bool)
	if !ok {
		return Config{}, errors.New("parse key runByDocker error")
	}

	retConfig.MaxMemTime, ok = config.Get("maxMemTime").(int64)
	if !ok {
		return Config{}, errors.New("parse key maxMemTime error")
	}

	retConfig.InsertInterval, ok = config.Get("insertInterval").(int64)
	if !ok {
		return Config{}, errors.New("parse key insertInterval error")
	}

	retConfig.PageSize, ok = config.Get("pageSize").(int64)
	if !ok {
		return Config{}, errors.New("parse key pageSize error")
	}

	retConfig.RetryCount, ok = config.Get("retryCount").(int64)
	if !ok {
		return Config{}, errors.New("parse key retryCount error")
	}

	retConfig.Proxy, ok = config.Get("proxy").(string)
	if !ok {
		return Config{}, errors.New("parse key proxy error")
	}

	retConfig.User, ok = config.Get("user").(string)
	if !ok {
		return Config{}, errors.New("parse key user error")
	}

	retConfig.Password, ok = config.Get("password").(string)
	if !ok {
		return Config{}, errors.New("parse key password error")
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

	symbolReplaces := config.Get("symbolReplaces")
	if symbolReplaces == nil {
		return Config{}, errors.New("get symbolReplaces error")
	}

	for _, replaceTemp := range symbolReplaces.([]interface{}) {
		replaceArr := strings.Split(replaceTemp.(string), ",")
		if len(replaceArr) != 3 {
			return Config{}, errors.New("incorrect symbol replace " + replaceTemp.(string))
		}
		oldSymbol := replaceArr[0]
		exchange := replaceArr[1]
		newSymbol := replaceArr[2]
		if retConfig.SymbolReplaces == nil {
			retConfig.SymbolReplaces = make(map[string]map[string]string)
		}
		if _, ok := retConfig.SymbolReplaces[oldSymbol]; !ok {
			retConfig.SymbolReplaces[oldSymbol] = make(map[string]string)
		}
		retConfig.SymbolReplaces[oldSymbol][exchange] = newSymbol
	}

	// try read mysql password from environment
	envMysqlPassword := os.Getenv("MYSQL_ROOT_PASSWORD")
	if envMysqlPassword != "" {
		retConfig.Mysql.Password = envMysqlPassword
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
