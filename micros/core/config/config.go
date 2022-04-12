package config

import (
	"bridge/common"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	HttpConfig        common.HttpConf
	DBconfig          common.DBconf
	RedisConfig       common.Redisconf
	Secrets           common.Secrets
	Mailerconf        common.Mailerconf
	TemporalCliConfig common.TemporalCliconf
	Casbin            common.CasbinCnf
	EthereumConfig    common.EtherumConfig
	EthGovContract    string
	EthImportContract string

	WelupsConfig   common.WelupsConfig
	WelGovContract string
}

func parseEnv() Env {

	// read env
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		viper.AutomaticEnv()
	}
	CORSstring := common.WithDefault("APP_CORS", "*")
	CORS := strings.Split(CORSstring, " ")
	fmt.Println(CORS)

	return Env{

		HttpConfig: common.HttpConf{
			Host:             common.WithDefault("APP_HOST", ""),
			Port:             common.WithDefault("APP_PORT", 8001),
			Mode:             common.WithDefault("APP_MODE", "debug"), // "release", "test"
			CORSAllowOrigins: CORS,

			X509CertFile: common.WithDefault("APP_X509_CERT", "./cert.pem"),
			X509KeyFile:  common.WithDefault("APP_X509_KEY", "./key.pem"),
		},

		DBconfig: common.DBconf{
			Host:            common.WithDefault("APP_DB_HOST", "localhost"),
			Port:            common.WithDefault("APP_DB_PORT", 5432),
			Username:        common.WithDefault("APP_DB_USERNAME", "welbridge"),
			Password:        common.WithDefault("APP_DB_PASSWORD", ""),
			DBname:          common.WithDefault("APP_DB_NAME", "welbridge"),
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
			JwtSecret: common.WithDefault("APP_JWT_SECRET", "keepcalmandstaypositive"),
		},

		Mailerconf: common.Mailerconf{
			SmtpHost: common.WithDefault("APP_MAILER_SMTP_HOST", "smtp.gmail.com"),
			SmtpPort: common.WithDefault("APP_MAILER_SMTP_PORT", 587),
			Address:  common.WithDefault("APP_MAILER_ADDRESS", "bridgemail.welups@gmail.com"),
			Password: common.WithDefault("APP_MAILER_PASSWORD", "showmethemoney11!1"),
		},

		Casbin: common.CasbinCnf{
			ModelPath:  common.WithDefault("APP_CASBIN_MODEL_PATH", "./config/rbac/model.conf"),
			PolicyPath: common.WithDefault("APP_CASBIN_POLICY_PATH", "./config/rbac/policy.csv"),
		},

		EthereumConfig: common.EtherumConfig{
			BlockchainRPC: common.WithDefault("ETH_BLOCKCHAIN_RPC", "https://rinkeby.infura.io/v3/4e7b43c678a14cffbe49ed691311ff1a"),
			BlockTime:     common.WithDefault("ETH_BLOCK_TIME", uint64(14)),
			BlockOffSet:   common.WithDefault("ETH_BLOCK_OFFSET", int64(5)),
		},
		EthGovContract:    common.WithDefault("ETH_GOV_CONTRACT_ADDRESS", "0x6150f59d1fd1A3C6Dc22619973D574e3CC7bA6E2"),
		EthImportContract: common.WithDefault("ETH_IMPORT_CONTRACT_ADDRESS", "0xE5a7d2F508579C22238688AD0d90db20f708e2A5"),

		WelupsConfig: common.WelupsConfig{
			Nodes:         common.WithDefault("WEL_NODES", []string{"54.179.208.1:16669"}),
			BlockTime:     common.WithDefault("WEL_BLOCK_TIME", uint64(3)),
			ClientTimeout: common.WithDefault("WEL_CLIENT_TIMEOUT", int64(5)),
			BlockOffSet:   common.WithDefault("WEL_BLOCK_OFFSET", int64(20)),
		},
		WelGovContract: common.WithDefault("WEL_GOV_CONTRACT_ADDRESS", "WE8RFVk1GA5NhK8yLxHWkpuP1E5UVqX9tQ"),
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

type Config struct {
	Env
	Flags
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

	// init config
	cnf = &Config{
		Env:   env,
		Flags: flags,
	}
	return
}

func Get() *Config {
	if cnf == nil {
		Load()
	}
	return cnf
}
