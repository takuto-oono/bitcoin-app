package api

import (
	"net/http"

	"bitcoin-app-golang/config"
)

type IDRFAPI interface {
	PostBitFlyerTicker(ticker PostTickerDRFRequest) error
}

type DRFAPI struct {
	Config config.Config
	API    *API
}

func NewDRFAPI(cfg config.Config) IDRFAPI {
	return &DRFAPI{
		Config: cfg,
		API:    NewAPI(),
	}
}

func (d *DRFAPI) PostBitFlyerTicker(ticker PostTickerDRFRequest) error {
	url, err := DRFServerURL(d.Config.ServerURL.DRFServer).PostTicker()
	if err != nil {
		return err
	}

	return d.API.Do(http.MethodPost, ticker, nil, url, nil)
}
