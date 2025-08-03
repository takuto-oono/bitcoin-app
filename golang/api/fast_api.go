package api

import "bitcoin-app-golang/config"

type IFastAPI interface {
	GetHealthcheck() error
}

type FastAPI struct {
	Config config.Config
	API    *API
}

func NewFastAPI(cfg config.Config) IFastAPI {
	return &FastAPI{
		Config: cfg,
		API:    NewAPI(),
	}
}

func (f *FastAPI) GetHealthcheck() error {
	url, err := FastAPIURL(f.Config.ServerURL.FastAPIServer).GetHealthcheck()
	if err != nil {
		return err
	}

	if err := f.API.Do("GET", nil, nil, url, nil); err != nil {
		return err
	}

	return nil
}
