package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Address string `envconfig:"ADDRESS" default:":8000"`
}

func Load() Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
