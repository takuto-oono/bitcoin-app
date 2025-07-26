package main

import (
	"reflect"
	"testing"

	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
)

func Test_getDRFFromPostDRF(t *testing.T) {
	type args struct {
		ticker api.GetTickerFromDRFResponse
	}
	tests := []struct {
		name string
		args args
		want api.PostTickerDRFRequest
	}{
		{
			name: "正常なティッカーデータの変換",
			args: args{
				ticker: api.GetTickerFromDRFResponse{
					ID:              1,
					TickID:          12345,
					ProductCode:     "BTC_JPY",
					State:           "RUNNING",
					Timestamp:       "2023-01-01T00:00:00.000Z",
					BestBid:         4500000.0,
					BestAsk:         4500100.0,
					BestBidSize:     0.1,
					BestAskSize:     0.2,
					TotalBidDepth:   10.5,
					TotalAskDepth:   15.3,
					MarketBidSize:   5.0,
					MarketAskSize:   7.2,
					Ltp:             4500050.0,
					Volume:          100.5,
					VolumeByProduct: 200.8,
				},
			},
			want: api.PostTickerDRFRequest{
				TickID:          12345,
				ProductCode:     "BTC_JPY",
				State:           "RUNNING",
				Timestamp:       "2023-01-01T00:00:00.000Z",
				BestBid:         4500000.0,
				BestAsk:         4500100.0,
				BestBidSize:     0.1,
				BestAskSize:     0.2,
				TotalBidDepth:   10.5,
				TotalAskDepth:   15.3,
				MarketBidSize:   5.0,
				MarketAskSize:   7.2,
				Ltp:             4500050.0,
				Volume:          100.5,
				VolumeByProduct: 200.8,
			},
		},
		{
			name: "ゼロ値を含むティッカーデータの変換",
			args: args{
				ticker: api.GetTickerFromDRFResponse{
					ID:              2,
					TickID:          0,
					ProductCode:     "ETH_JPY",
					State:           "STOP",
					Timestamp:       "",
					BestBid:         0.0,
					BestAsk:         0.0,
					BestBidSize:     0.0,
					BestAskSize:     0.0,
					TotalBidDepth:   0.0,
					TotalAskDepth:   0.0,
					MarketBidSize:   0.0,
					MarketAskSize:   0.0,
					Ltp:             0.0,
					Volume:          0.0,
					VolumeByProduct: 0.0,
				},
			},
			want: api.PostTickerDRFRequest{
				TickID:          0,
				ProductCode:     "ETH_JPY",
				State:           "STOP",
				Timestamp:       "",
				BestBid:         0.0,
				BestAsk:         0.0,
				BestBidSize:     0.0,
				BestAskSize:     0.0,
				TotalBidDepth:   0.0,
				TotalAskDepth:   0.0,
				MarketBidSize:   0.0,
				MarketAskSize:   0.0,
				Ltp:             0.0,
				Volume:          0.0,
				VolumeByProduct: 0.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDRFFromPostDRF(tt.args.ticker); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDRFFromPostDRF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertTickers(t *testing.T) {
	type args struct {
		tickers []api.GetTickerFromDRFResponse
	}
	tests := []struct {
		name string
		args args
		want []api.PostTickerDRFRequest
	}{
		{
			name: "空のスライスの変換",
			args: args{
				tickers: []api.GetTickerFromDRFResponse{},
			},
			want: []api.PostTickerDRFRequest{},
		},
		{
			name: "単一のティッカーデータの変換",
			args: args{
				tickers: []api.GetTickerFromDRFResponse{
					{
						ID:              1,
						TickID:          12345,
						ProductCode:     "BTC_JPY",
						State:           "RUNNING",
						Timestamp:       "2023-01-01T00:00:00.000Z",
						BestBid:         4500000.0,
						BestAsk:         4500100.0,
						BestBidSize:     0.1,
						BestAskSize:     0.2,
						TotalBidDepth:   10.5,
						TotalAskDepth:   15.3,
						MarketBidSize:   5.0,
						MarketAskSize:   7.2,
						Ltp:             4500050.0,
						Volume:          100.5,
						VolumeByProduct: 200.8,
					},
				},
			},
			want: []api.PostTickerDRFRequest{
				{
					TickID:          12345,
					ProductCode:     "BTC_JPY",
					State:           "RUNNING",
					Timestamp:       "2023-01-01T00:00:00.000Z",
					BestBid:         4500000.0,
					BestAsk:         4500100.0,
					BestBidSize:     0.1,
					BestAskSize:     0.2,
					TotalBidDepth:   10.5,
					TotalAskDepth:   15.3,
					MarketBidSize:   5.0,
					MarketAskSize:   7.2,
					Ltp:             4500050.0,
					Volume:          100.5,
					VolumeByProduct: 200.8,
				},
			},
		},
		{
			name: "複数のティッカーデータの変換",
			args: args{
				tickers: []api.GetTickerFromDRFResponse{
					{
						ID:              1,
						TickID:          12345,
						ProductCode:     "BTC_JPY",
						State:           "RUNNING",
						Timestamp:       "2023-01-01T00:00:00.000Z",
						BestBid:         4500000.0,
						BestAsk:         4500100.0,
						BestBidSize:     0.1,
						BestAskSize:     0.2,
						TotalBidDepth:   10.5,
						TotalAskDepth:   15.3,
						MarketBidSize:   5.0,
						MarketAskSize:   7.2,
						Ltp:             4500050.0,
						Volume:          100.5,
						VolumeByProduct: 200.8,
					},
					{
						ID:              2,
						TickID:          67890,
						ProductCode:     "ETH_JPY",
						State:           "RUNNING",
						Timestamp:       "2023-01-01T01:00:00.000Z",
						BestBid:         300000.0,
						BestAsk:         300100.0,
						BestBidSize:     1.0,
						BestAskSize:     2.0,
						TotalBidDepth:   50.0,
						TotalAskDepth:   60.0,
						MarketBidSize:   25.0,
						MarketAskSize:   30.0,
						Ltp:             300050.0,
						Volume:          500.0,
						VolumeByProduct: 1000.0,
					},
				},
			},
			want: []api.PostTickerDRFRequest{
				{
					TickID:          12345,
					ProductCode:     "BTC_JPY",
					State:           "RUNNING",
					Timestamp:       "2023-01-01T00:00:00.000Z",
					BestBid:         4500000.0,
					BestAsk:         4500100.0,
					BestBidSize:     0.1,
					BestAskSize:     0.2,
					TotalBidDepth:   10.5,
					TotalAskDepth:   15.3,
					MarketBidSize:   5.0,
					MarketAskSize:   7.2,
					Ltp:             4500050.0,
					Volume:          100.5,
					VolumeByProduct: 200.8,
				},
				{
					TickID:          67890,
					ProductCode:     "ETH_JPY",
					State:           "RUNNING",
					Timestamp:       "2023-01-01T01:00:00.000Z",
					BestBid:         300000.0,
					BestAsk:         300100.0,
					BestBidSize:     1.0,
					BestAskSize:     2.0,
					TotalBidDepth:   50.0,
					TotalAskDepth:   60.0,
					MarketBidSize:   25.0,
					MarketAskSize:   30.0,
					Ltp:             300050.0,
					Volume:          500.0,
					VolumeByProduct: 1000.0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertTickers(tt.args.tickers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertTickers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getProdConfig(t *testing.T) {
	tests := []struct {
		name string
		want config.Config
	}{
		{
			name: "プロダクション環境の設定取得",
			want: config.Config{
				ServerURL: config.ServerURL{
					DRFServer:    prodDRFServerURL,
					GolangServer: prodGolangServerURL,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getProdConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getProdConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLocalConfig(t *testing.T) {
	tests := []struct {
		name string
		want config.Config
	}{
		{
			name: "ローカル環境の設定取得",
			want: config.Config{
				ServerURL: config.ServerURL{
					DRFServer:    localDRFServerURL,
					GolangServer: localGolangServerURL,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLocalConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLocalConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
