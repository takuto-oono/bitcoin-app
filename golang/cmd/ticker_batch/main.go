package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
)

const (
	TickerInterval     = 1 * time.Second
	DefaultProductCode = api.ProductCodeBTCJPY
)

func main() {
	tomlFilePath := flag.String("toml", "toml/local.toml", "toml file path")
	envFilePath := flag.String("env", "env/.env.local", "env file path")
	flag.Parse()

	cfg, err := config.NewConfig(*tomlFilePath, *envFilePath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	golangServer := api.NewGolangServerAPI(cfg)
	drf := api.NewDRFAPI(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer func() {
		stop()
		log.Println("Shutting down gracefully...")
	}()

	runTickerBatch(ctx, golangServer, drf)
}

func runTickerBatch(ctx context.Context, golangServer api.IGolangServerAPI, drf api.IDRFAPI) {
	ticker := time.NewTicker(TickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Panic recovered in ticker case: %v", r)
					}
				}()
				getAndPostTicker(golangServer, drf)
			}()
		}
	}
}

func getAndPostTicker(golangServer api.IGolangServerAPI, drf api.IDRFAPI) {
	ticker, err := golangServer.GetBitFlyerTicker(api.ProductCode(DefaultProductCode))
	if err != nil {
		log.Printf("Error fetching ticker: %v", err)
		return
	}

	drfTicker := api.ConvertTickerFromGolang(ticker)
	if err := drf.PostBitFlyerTicker(drfTicker); err != nil {
		log.Printf("Error posting ticker: %v", err)
		return
	}

	log.Print("Ticker posted successfully")
}
