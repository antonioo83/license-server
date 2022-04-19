package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type Config struct {
	ServerAddress    string `env:"SERVER_ADDRESS"`
	BaseURL          string `env:"BASE_URL"`
	DatabaseDsn      string `env:"DATABASE_DSN"`
	FilepathToDBDump string
	Auth             Auth
	DeleteShotURL    DeleteShotURL
}

type Auth struct {
	Alg            string
	RememberMeTime time.Duration
	SignKey        []byte
	TokenName      string
}

type DeleteShotURL struct {
	WorkersCount int
	ChunkLength  int
}

var cfg Config

func GetConfigSettings() Config {
	const ServerAddress string = ":8080"
	const DatabaseDSN = "postgres://postgres:433370@localhost:5433/postgres"
	const AuthEncodeAlgorithm = "HS256"
	const AuthRememberMeTime = 60 * 30 * time.Second
	const AuthSignKey = "secret"
	const AuthTokenName = "token"

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

	cfg.Auth.Alg = AuthEncodeAlgorithm
	cfg.Auth.RememberMeTime = AuthRememberMeTime
	cfg.Auth.SignKey = []byte(AuthSignKey)
	cfg.Auth.TokenName = AuthTokenName

	cfg.DeleteShotURL.ChunkLength = 10
	cfg.DeleteShotURL.WorkersCount = 1

	return cfg
}
