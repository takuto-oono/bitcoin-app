package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"bitcoin-app/golang/config"
)

func main() {
	tomlFilePath := flag.String("toml", "toml/local.toml", "toml file path")
	flag.Parse()

	cfg, err := config.NewConfig(*tomlFilePath)
	if err != nil {
		panic(err)
	}

	// TODO: ginのdefaultは色々いけないらしいので改良する。
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "hello golang server")
	})

	if err := router.Run(fmt.Sprintf(":%s", cfg.GeneralSetting.Port)); err != nil {
		panic(err)
	}
}
