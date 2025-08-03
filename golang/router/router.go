package router

import (
	"fmt"
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
	lineHandler, err := handler.NewLineHandler(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to create Line handler: %w", err))
	}

	r.GET("/healthcheck/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	bitflyer := r.Group("/bitflyer")
	bitflyer.GET("/ticker", bitFlyerHandler.GetTickerFromBitFlyer)
	bitflyer.POST("/order/buy", bitFlyerHandler.BuyOrder)
	bitflyer.POST("/order/sell", bitFlyerHandler.SellOrder)

	line := r.Group("/line")
	line.POST("/message", lineHandler.PostMessage)
	line.POST("/callback", lineHandler.CallbackMessage) // LINEのグループIDを取得するために実装したエンドポイントを一応残しておく

	return r
}
