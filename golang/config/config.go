package config

import (
	"errors"

	"github.com/BurntSushi/toml"
)

type GeneralSetting struct {
	Port string `toml:"port"`
}

type Config struct {
	GeneralSetting `toml:"general"`
}

func NewConfig(tomlFilePath string) (Config, error) {
	var cfg Config

	if _, err := toml.DecodeFile(tomlFilePath, &cfg); err != nil {
		return Config{}, err
	}

	if err := cfg.mustCheck(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) mustCheck() error {
	if c == nil {
		return errors.New("config is nil")
	}

	if c.GeneralSetting.Port == "" {
		return errors.New("port is empty")
	}

	return nil
}
