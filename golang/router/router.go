package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bitcoin-app-golang/config"
	"bitcoin-app-golang/handler"
)

func NewRouter(cfg config.Config) *gin.Engine {
	r := gin.Default()

	return setRoutes(r, cfg)
}

func setRoutes(r *gin.Engine, cfg config.Config) *gin.Engine {
	bitFlyerHandler := handler.NewBitFlyerHandler(cfg)

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "hello golang server")
	})

	r.GET("/bitflyer/ticker", bitFlyerHandler.GetTickerFromBitFlyer)

	return r
}
