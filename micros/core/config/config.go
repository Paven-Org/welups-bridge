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
	HttpConfig  common.HttpConf
	DBconfig    common.DBconf
	RedisConfig common.Redisconf
	Secrets     common.Secrets
	Mailerconf  common.Mailerconf
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
	fmt.Println(CORS)

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

		Secrets: common.Secrets{
			JwtSecret: common.WithDefault("APP_JWT_SECRET", "keepcalmandstaypositive"),
		},

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
