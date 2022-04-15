package config

import (
	"bridge/common"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	HttpConfig         common.HttpConf
	DBconfig           common.DBconf
	RedisConfig        common.Redisconf
	Secrets            common.Secrets
	EtherumConf        common.EtherumConfig
	WelupsConf         common.WelupsConfig
	WelContractAddress []string
	EthContractAddress []string
	Mailerconf         common.Mailerconf
	TemporalCliConfig  common.TemporalCliconf
}

func parseEnv() Env {

	// read env
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		viper.AutomaticEnv()
	}
	CORSstring := common.WithDefault("APP_CORS", "*")
	CORS := strings.Split(CORSstring, ":")

	return Env{

		HttpConfig: common.HttpConf{
			Host:             common.WithDefault("APP_HOST", ""),
			Port:             common.WithDefault("APP_PORT", 8001),
			Mode:             common.WithDefault("APP_MODE", "debug"), // "release", "test"
			CORSAllowOrigins: CORS,
		},

		DBconfig: common.DBconf{
			Host:            common.WithDefault("APP_DB_HOST", "localhost"),
			Port:            common.WithDefault("APP_DB_PORT", 5432),
			Username:        common.WithDefault("APP_DB_USERNAME", "welbridge"),
			Password:        common.WithDefault("APP_DB_PASSWORD", ""),
			DBname:          common.WithDefault("APP_DB_NAME", "welbridge_weleth"),
			DBbackend:       common.WithDefault("APP_DB_BACKEND", "postgres"),
			MaxOpenConns:    common.WithDefault("APP_DB_CONN_MAX", 10),
			MaxIdleConns:    common.WithDefault("APP_DB_CONN_MAX_IDLE", 5),
			ConnMaxLifetime: common.WithDefault("APP_DB_CONN_MAX_LIFETIME", 30*time.Second),
			SSLMode:         common.WithDefault("APP_DB_SSL_MODE", "disable"),
		},

		RedisConfig: common.Redisconf{
			Network:  common.WithDefault("APP_REDIS_NET", "tcp"),
			Host:     common.WithDefault("APP_REDIS_HOST", "localhost"),
			Port:     common.WithDefault("APP_REDIS_PORT", 6379),
			Username: common.WithDefault("APP_REDIS_USERNAME", ""),
			Password: common.WithDefault("APP_REDIS_PASSWORD", ""),
		},

		TemporalCliConfig: common.TemporalCliconf{
			Host:      common.WithDefault("APP_TEMPORAL_HOST", "localhost"),
			Port:      common.WithDefault("APP_TEMPORAL_POST", 7233),
			Namespace: common.WithDefault("APP_TEMPORAL_NAMESPACE", "default"), // "devWelbridge", "prodWelbridge"
			// Ideally this should be retrieved from some secret manager
			Secret: common.WithDefault("APP_TEMPORAL_SECRET", "411ab14d42f1f5cf668db2d6ebd73937"), // 16,24,32 bytes long for AES-128,192,256 respectively
		},

		Secrets: common.Secrets{
			// Ideally this should be retrieved from some secret manager
			JwtSecret: common.WithDefault("APP_JWT_SECRET", "keepcalmandstaypositive"),
		},

		EtherumConf: common.EtherumConfig{
			BlockchainRPC: common.WithDefault("ETH_BLOCKCHAIN_RPC", "https://rinkeby.infura.io/v3/4e7b43c678a14cffbe49ed691311ff1a"),
			BlockTime:     common.WithDefault("ETH_BLOCK_TIME", uint64(14)),
			BlockOffSet:   common.WithDefault("ETH_BLOCK_OFFSET", int64(5)),
		},
		EthContractAddress: common.WithDefault("ETH_CONTRACT_ADDRESS", []string{"0xE5a7d2F508579C22238688AD0d90db20f708e2A5"}),

		WelupsConf: common.WelupsConfig{
			Nodes:         common.WithDefault("WEL_NODES", []string{"54.179.208.1:16669"}),
			BlockTime:     common.WithDefault("WEL_BLOCK_TIME", uint64(3)),
			ClientTimeout: common.WithDefault("WEL_CLIENT_TIMEOUT", int64(5)),
			BlockOffSet:   common.WithDefault("WEL_BLOCK_OFFSET", int64(20)),
		},
		WelContractAddress: common.WithDefault("WEL_CONTRACT_ADDRESS", []string{"WUbnXM9M4QYEkksG3ADmSan2kY5xiHTr1E"}),

		Mailerconf: common.Mailerconf{
			SmtpHost: common.WithDefault("APP_MAILER_SMTP_HOST", "smtp.gmail.com"),
			SmtpPort: common.WithDefault("APP_MAILER_SMTP_PORT", 587),
			Address:  common.WithDefault("APP_MAILER_ADDRESS", "bridgemail.welups@gmail.com"),
			Password: common.WithDefault("APP_MAILER_PASSWORD", "showmethemoney11!1"),
		},
	}
}

type Flags struct {
	Structured bool
}

func parseFlags() Flags {

	// output structured log
	structured := flag.Bool("structuredLog", false, "structured log")

	// parse all flags
	flag.Parse()

	return Flags{
		Structured: *structured,
	}
}

type TokensMap []struct {
	Eth     string `json:"eth"`
	Wel     string `json:"wel"`
	EthName string `json:"eth_name"` // not actually used right now, but it's nice to have some clarity
	WelName string `json:"wel_name"` // not actually used right now, but it's nice to have some clarity
}

func ParseTokensMap() TokensMap {
	mapfile, err := os.Open("tokens-map.json")
	if err != nil {
		fmt.Println("[config] Unable to load tokens map, error: ", err.Error())
		panic(err)
	}
	defer mapfile.Close()

	var tkMaps TokensMap
	decoder := json.NewDecoder(mapfile)
	if err := decoder.Decode(&tkMaps); err != nil {
		fmt.Println("[config] Unable to parse tokens map, error: ", err.Error())
		panic(err)
	}
	return tkMaps
}

type Config struct {
	Env
	Flags
	TokensMap
}

var cnf *Config

func Load() {
	if cnf != nil {
		return
	}
	// parse flags
	flags := parseFlags()

	// parse env
	env := parseEnv()

	// tokens map
	tkmap := ParseTokensMap()

	// init config
	cnf = &Config{
		Env:       env,
		Flags:     flags,
		TokensMap: tkmap,
	}
	return
}

func Get() *Config {
	if cnf == nil {
		Load()
	}
	return cnf
}
