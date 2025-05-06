package usecase

import (
	"net/http"

	"bitcoin-app/golang/api"
	"bitcoin-app/golang/config"
)

type IBitFlyerUsecase interface {
	GetTicker(productCode string) (api.TickerFromBitFlyer, int, error)
}

type BitFlyerUsecase struct {
	Config      config.Config
	BitFlyerAPI api.IBitFlyerAPI
}

func NewBitFlyerUsecase(cfg config.Config) IBitFlyerUsecase {
	return &BitFlyerUsecase{
		Config:      cfg,
		BitFlyerAPI: api.NewBitFlyerAPI(cfg),
	}
}

func (b *BitFlyerUsecase) GetTicker(productCode string) (api.TickerFromBitFlyer, int, error) {
	pc, err := api.NewProductCode(productCode)
	if err != nil {
		return api.TickerFromBitFlyer{}, http.StatusBadRequest, err
	}

	res, err := b.BitFlyerAPI.GetTicker(pc)
	if err != nil {
		return api.TickerFromBitFlyer{}, http.StatusInternalServerError, err
	}

	return res, http.StatusOK, nil
}
