package api

import (
	"net/http"

	"bitcoin-app-golang/config"
)

type IGolangServerAPI interface {
	GetHealthcheck() error
	GetBitFlyerTicker(productCode string) (TickerFromGolangServer, error)
}

type GolangServerAPI struct {
	Config config.Config
	API    *API
}

func NewGolangServerAPI(cfg config.Config) IGolangServerAPI {
	return &GolangServerAPI{
		Config: cfg,
		API:    NewAPI(),
	}
}

func (g *GolangServerAPI) GetHealthcheck() error {
	url, err := GolangServerURL(g.Config.ServerURL.GolangServer).GetHealthcheck()
	if err != nil {
		return err
	}

	if err := g.API.Do(http.MethodGet, nil, nil, url, nil); err != nil {
		return err
	}

	return nil
}

func (g *GolangServerAPI) GetBitFlyerTicker(productCode string) (TickerFromGolangServer, error) {
	url, err := GolangServerURL(g.Config.ServerURL.GolangServer).GetTicker(productCode)
	if err != nil {
		return TickerFromGolangServer{}, err
	}

	var resModel TickerFromGolangServer
	if err := g.API.Do(http.MethodGet, nil, &resModel, url, nil); err != nil {
		return TickerFromGolangServer{}, err
	}

	return resModel, nil
}
