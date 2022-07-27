package common

import (
	"time"

	"github.com/spf13/viper"
)

const (
	LocalEnv      = "local"
	DevEnv        = "dev"
	StagingEnv    = "staging"
	ProductionEnv = "production"
)

type HttpConf struct {
	Host             string
	Port             int
	Mode             string
	CORSAllowOrigins []string
	X509CertFile     string
	X509KeyFile      string
}

type Secrets struct {
	JwtSecret string
}

type DBconf struct {
	Host            string
	Port            int
	Username        string
	Password        string
	DBname          string
	DBbackend       string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	SSLMode         string
}

type Redisconf struct {
	Network  string
	Host     string
	Port     int
	Username string
	Password string
}

type Mailerconf struct {
	SmtpHost string
	SmtpPort int
	Address  string
	Password string
}

type TemporalCliconf struct {
	Host      string
	Port      int
	Namespace string
	Secret    string
}

type CasbinCnf struct {
	ModelPath  string
	PolicyPath string // only support policy file for now
}

type EtherumConfig struct {
	BlockchainRPC string
	BlockTime     uint64
	BlockOffSet   int64
}

type WelupsConfig struct {
	Nodes         []string
	BlockTime     uint64
	BlockOffSet   int64
	ClientTimeout int64
}

func WithDefault[A any](key string, df A) A {
	viper.SetTypeByDefaultValue(true) // make sure viper.Get would always Get() value as A
	viper.SetDefault(key, df)
	return viper.Get(key).(A)
}
