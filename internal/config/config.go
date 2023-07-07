package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	RunAddr     string `env:"RUN_ADDRESS" envDefault:":8080"`
	AccrualAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURI string `env:"DATABASE_URI"`
	Seed        string `env:"SEED" envDefault:"b4952c3809196592c026529df00774e46bfb5be0"`
	LogLevel    int    `env:"LOG_LEVEL" envDefault:"0"`
}

var config ServerConfig

func ParseFlags() (*ServerConfig, error) {
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	flag.StringVar(&config.RunAddr, "a", config.RunAddr, "address and port to run server")
	flag.StringVar(&config.AccrualAddr, "r", config.AccrualAddr, "Accrual System Address")
	flag.StringVar(&config.DatabaseURI, "d", config.DatabaseURI, "Data Source Name (DSN)")
	flag.StringVar(&config.Seed, "s", config.Seed, "seed")
	flag.IntVar(&config.LogLevel, "l", config.LogLevel, "0, 1, 2, 3")
	flag.Parse()

	return &config, nil
}
