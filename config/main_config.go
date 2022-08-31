package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress     string `env:"SERVER_ADDRESS"`
	BaseURL           string `env:"BASE_URL"`
	DatabaseDsn       string `env:"DATABASE_DSN"`
	FilepathToDBDump  string
	RequestTimeoutSec int64
	Auth              Auth
	Callback          Callback
}

type Auth struct {
	AdminAuthToken string
}

type Callback struct {
	MaxAttempts     uint
	LimitUnitOfTime uint
	CronSpec        string
}

var cfg Config

func GetConfigSettings() Config {
	const (
		ServerAddress           = ":8080"
		DatabaseDSN             = "postgres://postgres:433370@localhost:5432/license_server"
		RequestTimeoutSec       = 60
		AdminAuthToken          = "54d1ba805e2a4891aeac9299b618945e"
		CallbackMaxAttempts     = 3
		CallbackLimitUnitOfTime = 50
		CallbackCronSpec        = "1 * * * * *"
	)

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "The address of the local server")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base address of the result short url")
	flag.StringVar(&cfg.DatabaseDsn, "d", cfg.DatabaseDsn, "Database port")
	flag.Parse()
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ServerAddress
	}

	if cfg.DatabaseDsn == "" {
		cfg.DatabaseDsn = DatabaseDSN
	}

	if cfg.RequestTimeoutSec == 0 {
		cfg.RequestTimeoutSec = RequestTimeoutSec
	}

	cfg.Auth.AdminAuthToken = AdminAuthToken

	cfg.Callback.MaxAttempts = CallbackMaxAttempts
	cfg.Callback.LimitUnitOfTime = CallbackLimitUnitOfTime
	cfg.Callback.CronSpec = CallbackCronSpec

	return cfg
}
