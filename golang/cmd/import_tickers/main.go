package main

import (
	"fmt"
	"sync"

	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
)

const (
	prodDRFServerURL  = "http://localhost:9000"
	localDRFServerURL = "http://localhost:8000"

	prodGolangServerURL  = "http://localhost:9080"
	localGolangServerURL = "http://localhost:8080"

	postProcessNum   = 10
	deleteProcessNum = 10
)

func main() {
	localCfg := getLocalConfig()
	prodCfg := getProdConfig()

	localTickers, err := getTickers(localCfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to get tickers: %v", err))
	}

	deleteErrs := deleteTickerProcess(deleteTickerIDs(localTickers))

	localTickersAfterDeleteProcess, err := getTickers(localCfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to get tickers after delete process: %v", err))
	}

	prodTickers, err := getTickers(prodCfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to get production tickers: %v", err))
	}

	convertedTickers := convertTickers(prodTickers)

	postErrs := postTickerProcess(convertedTickers)

	fmt.Println("Ticker import completed successfully.")
	fmt.Printf("Total tickers imported: %d\n", len(convertedTickers))

	fmt.Printf("Tickers after delete process: %d\n", len(localTickersAfterDeleteProcess))

	for _, err := range deleteErrs {
		fmt.Printf("Failed to delete ticker with ID %d: %v\n", err.ID, err.Err)
	}

	for _, err := range postErrs {
		fmt.Printf("Failed to post ticker with ID %d: %v\n", err.ID, err.Err)
	}
}

func getProdConfig() config.Config {
	return config.Config{
		ServerURL: config.ServerURL{
			DRFServer:    prodDRFServerURL,
			GolangServer: prodGolangServerURL,
		},
	}
}

func getLocalConfig() config.Config {
	return config.Config{
		ServerURL: config.ServerURL{
			DRFServer:    localDRFServerURL,
			GolangServer: localGolangServerURL,
		},
	}
}

func deleteTickerIDs(tickers []api.GetTickerFromDRFResponse) []int {
	tickerIDs := make([]int, len(tickers))

	for i, ticker := range tickers {
		tickerIDs[i] = ticker.ID
	}

	return tickerIDs
}

func getTickers(cfg config.Config) ([]api.GetTickerFromDRFResponse, error) {
	drfAPI := api.NewDRFAPI(cfg)
	tickers, err := drfAPI.GetBitFlyerTickers()
	if err != nil {
		return nil, fmt.Errorf("failed to get tickers: %w", err)
	}
	return tickers, nil
}

type ProcessError struct {
	Err error
	ID  int
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

func postTickerProcess(tickers []api.PostTickerDRFRequest) []ProcessError {
	drfAPI := api.NewDRFAPI(getLocalConfig())

	var wg sync.WaitGroup
	wg.Add(postProcessNum)
	errChan := make(chan ProcessError, len(tickers))
	processErrs := make([]ProcessError, 0)

	postTickers := func(processNum int) {
		for i, ticker := range tickers {
			if i%postProcessNum == processNum {
				if err := drfAPI.PostBitFlyerTicker(ticker); err != nil {
					fmt.Printf("Process %d: Failed to post ticker: %v\n", processNum, err)
					errChan <- ProcessError{Err: err, ID: i}
				}
			}
		}
	}

	for i := 0; i < postProcessNum; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
				fmt.Printf("post process %d completed\n", i)
			}()
			fmt.Printf("Starting post process %d\n", i)
			postTickers(i)
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		processErrs = append(processErrs, err)
	}

	return processErrs
}

func deleteTickerProcess(tickerIDs []int) []ProcessError {
	drfAPI := api.NewDRFAPI(getLocalConfig())

	var wg sync.WaitGroup
	wg.Add(deleteProcessNum)
	errChan := make(chan ProcessError, len(tickerIDs))
	processErrors := make([]ProcessError, 0)

	deleteTickers := func(processNum int) {
		for i, tickerID := range tickerIDs {
			if i%deleteProcessNum == processNum {
				if err := drfAPI.DeleteBitFlyerTicker(tickerID); err != nil {
					fmt.Printf("Process %d: Failed to delete ticker with ID %d: %v\n", processNum, tickerID, err)
					errChan <- ProcessError{Err: err, ID: tickerID}
				}
			}
		}
	}

	for i := 0; i < deleteProcessNum; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
				fmt.Printf("delete Process %d completed\n", i)
			}()
			fmt.Printf("Starting delete process %d\n", i)
			deleteTickers(i)
		}(i)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		processErrors = append(processErrors, err)
	}

	return processErrors
}

func convertTickers(tickers []api.GetTickerFromDRFResponse) []api.PostTickerDRFRequest {
	converted := make([]api.PostTickerDRFRequest, len(tickers))
	for i, ticker := range tickers {
		converted[i] = getDRFFromPostDRF(ticker)
	}
	return converted
}
