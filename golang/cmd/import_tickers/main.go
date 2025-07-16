package main

import (
	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
	"fmt"
	"time"
)

func main() {
	localCfg := getLocalConfig()
	prodCfg := getProdConfig()

	tickers, err := getTickers(localCfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to get tickers: %v", err))
	}

	convertedTickers := convertTickers(tickers)

	for _, ticker := range convertedTickers {
		time.Sleep(1 * time.Second)

		err := postTicker(prodCfg, ticker)
		if err != nil {
			fmt.Printf("Failed to post ticker %s: %v\n", ticker.ProductCode, err)
		} else {
			fmt.Printf("Successfully posted ticker %s\n", ticker.ProductCode)
		}
	}
}

func getProdConfig() config.Config {
	cfg, err := config.NewConfig("../../toml/prod.toml", "../../env/.env.prod")
	if err != nil {
		panic(err)
	}

	cfg.ServerURL.DRFServer = "http://localhost:7000" // Override for production

	return cfg
}

func getLocalConfig() config.Config {
	cfg, err := config.NewConfig("../../toml/local.toml", "../../env/.env.local")
	if err != nil {
		panic(err)
	}
	return cfg
}

func getTickers(cfg config.Config) ([]api.GetTickerFromDRFResponse, error) {
	drfAPI := api.NewDRFAPI(cfg)
	tickers, err := drfAPI.GetBitFlyerTickers()
	if err != nil {
		return nil, fmt.Errorf("failed to get tickers: %w", err)
	}
	return tickers, nil
}

func postTicker(cfg config.Config, ticker api.PostTickerDRFRequest) error {
	drfAPI := api.NewDRFAPI(cfg)
	err := drfAPI.PostBitFlyerTicker(ticker)
	if err != nil {
		return fmt.Errorf("failed to post ticker: %w", err)
	}
	return nil
}

func getDRFFromPostDRF(ticker api.GetTickerFromDRFResponse) api.PostTickerDRFRequest {
	return api.PostTickerDRFRequest{
		TickID:          ticker.TickID,
		ProductCode:     ticker.ProductCode,
		State:           ticker.State,
		Timestamp:       ticker.Timestamp,
		BestBid:         ticker.BestBid,
		BestAsk:         ticker.BestAsk,
		BestBidSize:     ticker.BestBidSize,
		BestAskSize:     ticker.BestAskSize,
		TotalBidDepth:   ticker.TotalBidDepth,
		TotalAskDepth:   ticker.TotalAskDepth,
		MarketBidSize:   ticker.MarketBidSize,
		MarketAskSize:   ticker.MarketAskSize,
		Ltp:             ticker.Ltp,
		Volume:          ticker.Volume,
		VolumeByProduct: ticker.VolumeByProduct,
	}
}

func convertTickers(tickers []api.GetTickerFromDRFResponse) []api.PostTickerDRFRequest {
	converted := make([]api.PostTickerDRFRequest, len(tickers))
	for i, ticker := range tickers {
		converted[i] = getDRFFromPostDRF(ticker)
	}
	return converted
}
