package api

import (
	"reflect"
	"testing"
)

func TestNewProductCode(t *testing.T) {
	type args struct {
		productCode string
	}
	tests := []struct {
		name    string
		args    args
		want    ProductCode
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				productCode: ProductCodeBTCJPY,
			},
			want:    ProductCode(ProductCodeBTCJPY),
			wantErr: false,
		},
		{
			name: "validate error",
			args: args{
				productCode: "invalide",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProductCode(tt.args.productCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProductCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewProductCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductCode_Validate(t *testing.T) {
	tests := []struct {
		name string
		p    ProductCode
		want bool
	}{
		{
			name: "success",
			p:    ProductCode(ProductCodeBTCJPY),
			want: true,
		},
		{
			name: "invalid",
			p:    ProductCode("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Validate(); got != tt.want {
				t.Errorf("ProductCode.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
					ProductCode:     ProductCodeBTCJPY,
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
				ProductCode:     ProductCodeBTCJPY,
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
