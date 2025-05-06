package handler

import (
	"github.com/gin-gonic/gin"

	"bitcoin-app-golang/config"
	"bitcoin-app-golang/usecase"
)

type IBitFlyerHandler interface {
	GetTickerFromBitFlyer(ctx *gin.Context)
}

type BitFlyerHandler struct {
	Config config.Config

	UseCase usecase.IBitFlyerUsecase
}

func NewBitFlyerHandler(cfg config.Config) IBitFlyerHandler {
	usecase := usecase.NewBitFlyerUsecase(cfg)

	return &BitFlyerHandler{
		Config:  cfg,
		UseCase: usecase,
	}
}

func (h *BitFlyerHandler) GetTickerFromBitFlyer(ctx *gin.Context) {
	productCode := ctx.Request.URL.Query().Get("product_code")

	ticker, statusCode, err := h.UseCase.GetTicker(productCode)
	if err != nil {
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(statusCode, ticker)
}
