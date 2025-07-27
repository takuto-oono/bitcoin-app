package config

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Credential string

func (c Credential) String() string {
	return string("************")
}

func (c Credential) GoString() string {
	return "************"
}

type ServerURL struct {
	GolangServer string `toml:"golangServer"`
	DRFServer    string `toml:"drfServer"`
}

type BitFlyer struct {
	ApiKey    Credential
	ApiSecret Credential
}

type TickerBatch struct {
	BatchIntervalSec int `toml:"batchIntervalSec"`
}

type Config struct {
	ServerURL `toml:"serverURL"`
	BitFlyer
	TickerBatch `toml:"tickerBatch"`
}

func NewConfig(tomlFilePath, envFilePath string) (Config, error) {
	var cfg Config

	if err := cfg.setFromToml(tomlFilePath); err != nil {
		return Config{}, err
	}

	if err := cfg.setFromEnv(envFilePath); err != nil {
		return Config{}, err
	}

	if err := cfg.mustCheck(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) setFromToml(tomlFilePath string) error {
	_, err := toml.DecodeFile(tomlFilePath, c)
	return err
}

func (c *Config) setFromEnv(envFilePath string) error {
	if err := godotenv.Load(envFilePath); err != nil {
		return err
	}

	c.BitFlyer.ApiKey = Credential(os.Getenv("BITFLYER_API_KEY"))
	c.BitFlyer.ApiSecret = Credential(os.Getenv("BITFLYER_API_SECRET"))

	return nil
}

func (c *Config) mustCheck() error {
	if c == nil {
		return errors.New("config is nil")
	}

	if c.ServerURL.GolangServer == "" {
		return errors.New("golang server is empty")
	}

	if c.ServerURL.DRFServer == "" {
		return errors.New("drf server is empty")
	}

	if c.BitFlyer.ApiKey == "" {
		return errors.New("bitflyer api key is empty")
	}

	if c.BitFlyer.ApiSecret == "" {
		return errors.New("bitflyer api secret is empty")
	}

	if c.TickerBatch.BatchIntervalSec <= 0 {
		return errors.New("ticker batch interval must be greater than 0")
	}

	return nil
}
