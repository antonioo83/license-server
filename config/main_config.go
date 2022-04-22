package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress    string `env:"SERVER_ADDRESS"`
	BaseURL          string `env:"BASE_URL"`
	DatabaseDsn      string `env:"DATABASE_DSN"`
	FilepathToDBDump string
	Auth             Auth
}

type Auth struct {
	AdminAuthToken string
}

var cfg Config

func GetConfigSettings() Config {
	const ServerAddress string = ":8080"
	const DatabaseDSN = "postgres://postgres:433370@localhost:5433/license_server"
	const AdminAuthToken = "54d1ba805e2a4891aeac9299b618945e"

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

	cfg.Auth.AdminAuthToken = AdminAuthToken

	return cfg
}
