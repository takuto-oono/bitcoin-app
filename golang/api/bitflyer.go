package api

import (
	"net/http"

	"bitcoin-app-golang/config"
)

const BitFlyerBaseURL = "https://api.bitflyer.com"

type IBitFlyerAPI interface {
	GetTicker(ProductCode) (TickerFromBitFlyer, error)
}

type BitFlyerAPI struct {
	Config config.Config
	API    *API
}

func NewBitFlyerAPI(cfg config.Config) IBitFlyerAPI {
	return &BitFlyerAPI{
		Config: cfg,
		API:    NewAPI(),
	}
}

func (b *BitFlyerAPI) GetTicker(productCode ProductCode) (TickerFromBitFlyer, error) {
	url, err := BitFlyerURL(BitFlyerBaseURL).GetTicker(productCode)
	if err != nil {
		return TickerFromBitFlyer{}, err
	}

	resModel := TickerFromBitFlyer{}
	if err := b.API.Do(http.MethodGet, nil, &resModel, url, nil); err != nil {
		return TickerFromBitFlyer{}, err
	}
	return resModel, nil
}
