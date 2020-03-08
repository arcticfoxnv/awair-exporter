package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	AccessToken string `toml:"access_token"`
	Tier        string
}

func LoadConfig(filename string) (*Config, error) {
	config := &Config{}
	if _, err := toml.DecodeFile(filename, config); err != nil {
		return nil, err
	}

	return config, nil
}
