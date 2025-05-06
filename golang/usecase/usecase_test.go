package usecase

import "bitcoin-app/golang/config"


var TestConfig config.Config

func init() {
	var err error

	TestConfig, err = config.NewConfig("../toml/local.toml", "../env/.env.test")
	if err != nil {
		panic(err)
	}
}
