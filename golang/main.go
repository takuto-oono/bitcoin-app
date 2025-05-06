package main

import (
	"flag"
	"fmt"

	"bitcoin-app/golang/config"
	"bitcoin-app/golang/router"
)

func main() {
	tomlFilePath := flag.String("toml", "toml/local.toml", "toml file path")
	envFilePath := flag.String("env", "env/.env.local", "env file path")
	flag.Parse()

	cfg, err := config.NewConfig(*tomlFilePath, *envFilePath)
	if err != nil {
		panic(err)
	}

	router := router.NewRouter(cfg)

	if err := router.Run(fmt.Sprintf(":%s", cfg.GeneralSetting.Port)); err != nil {
		panic(err)
	}
}
