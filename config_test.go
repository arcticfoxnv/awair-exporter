package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPreflightCheckOK(t *testing.T) {
	cfg, _ := loadConfig()

	assert.Nil(t, preflightCheck(cfg))
}

func TestPrelightCheckErrAccessToken(t *testing.T) {
	cfg, _ := loadConfig()
	cfg.Set(CFG_ACCESS_TOKEN, "")

	assert.NotNil(t, preflightCheck(cfg))
}

func TestLoadConfig(t *testing.T) {
	_, err := loadConfig()

	assert.Nil(t, err)
}

func TestLoadConfigFile(t *testing.T) {
	os.Setenv("AWAIR_CONFIG_FILE", "awair.toml")
	_, err := loadConfig()

	assert.Nil(t, err)
}

func TestLoadConfigError(t *testing.T) {
	os.Setenv("AWAIR_CONFIG_FILE", ".does.not.exist.toml")
	_, err := loadConfig()

	assert.NotNil(t, err)
}
