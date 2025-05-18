package api

import (
	"fmt"
	"slices"
)

const (
	ProductCodeBTCJPY  = "BTC_JPY"
	ProductCodeXRPJPY  = "XRP_JPY"
	ProductCodeETHJPY  = "ETH_JPY"
	ProductCodeXLMJPY  = "XLM_JPY"
	ProductCOdeMONAJPY = "MONA_JPY"

	ProductCodeETHBTC   = "ETH_BTC"
	ProductCodeBCHBTC   = "BCH_BTC"
	ProductCodeFXBTCJPY = "FX_BTC_JPY"

	GetExecutionsDefaultCount = "100"

	ChildOrderTypeLimit  = "LIMIT"
	ChildOrderTypeMarket = "MARKET"

	SideBuy  = "BUY"
	SideSell = "SELL"

	TimeInForceGTC = "GTC"
	TimeInForceIOC = "IOC"
	TimeInForceFOK = "FOK"
)

type TickerFromBitFlyer struct {
	TickID          int     `json:"tick_id"`
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

type ProductCode string

func NewProductCode(productCode string) (ProductCode, error) {
	pc := ProductCode(productCode)

	if !pc.Validate() {
		return ProductCode(""), fmt.Errorf("invalid product code: %s", productCode)
	}

	return pc, nil
}

func (p ProductCode) Validate() bool {
	allowProductCodes := []string{
		ProductCodeBTCJPY,
		ProductCodeXRPJPY,
		ProductCodeETHJPY,
		ProductCodeXLMJPY,
		ProductCOdeMONAJPY,
		ProductCodeETHBTC,
		ProductCodeBCHBTC,
		ProductCodeFXBTCJPY,
	}

	return slices.Contains(allowProductCodes, string(p))
}

type TickerFromGolangServer struct {
	TickID          int     `json:"tick_id"`
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

type PostTickerDRFRequest struct {
	TickID          int     `json:"tick_id"`
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

func ConvertTicker(golangTicker TickerFromGolangServer) PostTickerDRFRequest {
	return PostTickerDRFRequest(golangTicker)
}
