package api

import (
	"net/http"

	"bitcoin-app-golang/config"
)

type IGolangServerAPI interface {
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
