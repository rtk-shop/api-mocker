package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port     string `env:"APP_PORT,required"`
	ApiToken string `env:"API_TOKEN,required"`
	ApiURL   string `env:"API_URL,required"`
	LogFile  string `env:"LOG_FILE"`
}

func New() *Config {

	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalln("config parse error", err)
	}

	if cfg.LogFile == "" {
		cfg.LogFile = "zap.log"
	}

	return &cfg
}
