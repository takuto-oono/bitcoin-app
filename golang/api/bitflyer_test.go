package api

import (
	"bitcoin-app/golang/config"
	"reflect"
	"testing"
)

func TestNewBitFlyerAPI(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name string
		args args
		want IBitFlyerAPI
	}{
		{
			name: "success",
			args: args{
				cfg: testConfig,
			},
			want: &BitFlyerAPI{
				Config: testConfig,
				API:    NewAPI(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBitFlyerAPI(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBitFlyerAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitFlyerAPI_GetTicker(t *testing.T) {
	type fields struct {
		Config config.Config
		API    *API
	}
	type args struct {
		productCode ProductCode
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantCheckFunc func(want TickerFromBitFlyer) bool
		wantErr       bool
	}{
		{
			name: "success",
			fields: fields{
				Config: testConfig,
				API:    NewAPI(),
			},
			args: args{
				productCode: "BTC_JPY",
			},
			wantCheckFunc: func(want TickerFromBitFlyer) bool {
				return want.TickID > 0 && want.ProductCode == ProductCodeBTCJPY
			},
			wantErr: false,
		},
		{
			name: "fail",
			fields: fields{
				Config: testConfig,
				API:    NewAPI(),
			},
			args: args{
				productCode: "",
			},
			wantCheckFunc: func(want TickerFromBitFlyer) bool {
				return want.TickID > 0 && want.ProductCode == ProductCodeBTCJPY
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BitFlyerAPI{
				Config: tt.fields.Config,
				API:    tt.fields.API,
			}
			got, err := b.GetTicker(tt.args.productCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyerAPI.GetTicker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantCheckFunc(got) {
				t.Errorf("BitFlyerAPI.GetTicker() = %v", got)
			}
		})
	}
}
