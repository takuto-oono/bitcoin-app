package main

import (
	"flag"
	"fmt"

	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
	"bitcoin-app-golang/router"
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

	port, err := api.ExtractPort(cfg.ServerURL.GolangServer)
	if err != nil {
		panic(err)
	}

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		panic(err)
	}
}
