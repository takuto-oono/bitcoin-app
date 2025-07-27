package api

import (
	"reflect"
	"testing"
)

func TestConvertTicker(t *testing.T) {
	type args struct {
		golangTicker TickerFromGolangServer
	}
	tests := []struct {
		name string
		args args
		want PostTickerDRFRequest
	}{
		{
			name: "success",
			args: args{
				golangTicker: TickerFromGolangServer{
					TickID:          1,
					ProductCode:     "BTC_JPY",
					State:           "RUNNING",
					Timestamp:       "2023-10-01T00:00:00Z",
					BestBid:         5000000,
					BestAsk:         5100000,
					BestBidSize:     0.1,
					BestAskSize:     0.2,
					TotalBidDepth:   1000000,
					TotalAskDepth:   2000000,
					MarketBidSize:   0.5,
					MarketAskSize:   0.6,
					Ltp:             5050000,
					Volume:          100,
					VolumeByProduct: 50,
				},
			},
			want: PostTickerDRFRequest{
				TickID:          1,
				ProductCode:     "BTC_JPY",
				State:           "RUNNING",
				Timestamp:       "2023-10-01T00:00:00Z",
				BestBid:         5000000,
				BestAsk:         5100000,
				BestBidSize:     0.1,
				BestAskSize:     0.2,
				TotalBidDepth:   1000000,
				TotalAskDepth:   2000000,
				MarketBidSize:   0.5,
				MarketAskSize:   0.6,
				Ltp:             5050000,
				Volume:          100,
				VolumeByProduct: 50,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertTickerFromGolang(tt.args.golangTicker); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertTicker() = %v, want %v", got, tt.want)
			}
		})
	}
}
