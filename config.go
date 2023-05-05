package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Url    string `toml:"url"`
	ApiKey string `toml:"api_key"`
}

func ReadConfig() (*Config, error) {
	path := os.Getenv("HOME") + "/.config/cbctl/config.toml"
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	if err := toml.Unmarshal(f, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
