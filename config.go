package main

import (
	"errors"
	"github.com/spf13/viper"
	"os"
)

const (
	CFG_ACCESS_TOKEN = "access_token"
	CFG_TIER_NAME    = "tier_name"
	CFG_LISTEN_PORT  = "listen_port"
)

func loadConfig() (*viper.Viper, error) {
	v := viper.GetViper()

	// Configure viper
	v.SetConfigName("awair")
	v.SetConfigType("toml")
	v.AddConfigPath("/etc")
	v.AddConfigPath(".")
	v.SetEnvPrefix("awair")
	v.AutomaticEnv()

	if path, present := os.LookupEnv("AWAIR_CONFIG_FILE"); present {
		v.SetConfigFile(path)
	}

	// Configure defaults
	v.SetDefault(CFG_LISTEN_PORT, 8080)

	// Read config
	if err := v.ReadInConfig(); err != nil {
		return v, err
	}

	return v, nil
}

func preflightCheck(v *viper.Viper) error {
	if v.GetString(CFG_ACCESS_TOKEN) == "" {
		return errors.New("Cannot start exporter, access token is missing")
	}

	return nil
}
