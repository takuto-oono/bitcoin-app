package usecase

import (
	"net/http"
	"reflect"
	"testing"

	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
)

func TestNewBitFlyerUsecase(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name string
		args args
		want IBitFlyerUsecase
	}{
		{
			name: "success",
			args: args{
				cfg: TestConfig,
			},
			want: &BitFlyerUsecase{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBitFlyerUsecase(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBitFlyerUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitFlyerUsecase_GetTicker(t *testing.T) {
	type fields struct {
		Config      config.Config
		BitFlyerAPI api.IBitFlyerAPI
	}
	type args struct {
		productCode string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantCheckFunc func(want api.TickerFromBitFlyer) bool
		want1         int
		wantErr       bool
	}{
		{
			name: "success",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				productCode: "BTC_JPY",
			},
			wantCheckFunc: func(want api.TickerFromBitFlyer) bool {
				return want.TickID > 0 && want.ProductCode == "BTC_JPY"
			},
			want1:   http.StatusOK,
			wantErr: false,
		},
		{
			name: "productCode is empty",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				productCode: "",
			},
			wantCheckFunc: func(want api.TickerFromBitFlyer) bool {
				return want.TickID == 0 && want.ProductCode == ""
			},
			want1:   http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "productCode is invalid",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				productCode: "invalid",
			},
			wantCheckFunc: func(want api.TickerFromBitFlyer) bool {
				return want.TickID == 0 && want.ProductCode == ""
			},
			want1:   http.StatusBadRequest,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BitFlyerUsecase{
				Config:      tt.fields.Config,
				BitFlyerAPI: tt.fields.BitFlyerAPI,
			}
			got, got1, err := b.GetTicker(tt.args.productCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyerUsecase.GetTicker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantCheckFunc(got) {
				t.Errorf("BitFlyerUsecase.GetTicker() = %v, want %v", got, tt.wantCheckFunc(got))
			}
			if got1 != tt.want1 {
				t.Errorf("BitFlyerUsecase.GetTicker() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
