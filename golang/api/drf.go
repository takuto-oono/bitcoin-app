package api

import (
	"net/http"

	"bitcoin-app-golang/config"
)

type IDRFAPI interface {
	GetBitFlyerTickers() ([]GetTickerFromDRFResponse, error)
	PostBitFlyerTicker(ticker PostTickerDRFRequest) error
	DeleteBitFlyerTicker(id int) error
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

func (d *DRFAPI) GetBitFlyerTickers() ([]GetTickerFromDRFResponse, error) {
	url, err := DRFServerURL(d.Config.ServerURL.DRFServer).GetTickers()
	if err != nil {
		return nil, err
	}
	var tickers []GetTickerFromDRFResponse
	if err := d.API.Do(http.MethodGet, nil, &tickers, url, nil); err != nil {
		return nil, err
	}
	return tickers, nil
}

func (d *DRFAPI) PostBitFlyerTicker(ticker PostTickerDRFRequest) error {
	url, err := DRFServerURL(d.Config.ServerURL.DRFServer).PostTicker()
	if err != nil {
		return err
	}

	return d.API.Do(http.MethodPost, ticker, nil, url, nil)
}

func (d *DRFAPI) DeleteBitFlyerTicker(id int) error {
	url, err := DRFServerURL(d.Config.ServerURL.DRFServer).DeleteTicker(id)
	if err != nil {
		return err
	}

	return d.API.Do(http.MethodDelete, nil, nil, url, nil)
}
