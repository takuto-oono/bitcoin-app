package api

import (
	"bitcoin-app-golang/config"
	"bitcoin-app-golang/consts"
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
		productCode string
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
				productCode: consts.ProductCodeBTCJPY,
			},
			wantCheckFunc: func(want TickerFromBitFlyer) bool {
				return want.TickID > 0 && want.ProductCode == consts.ProductCodeBTCJPY
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
				return want.TickID > 0 && want.ProductCode == consts.ProductCodeBTCJPY
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

func TestBitFlyerAPI_SendChildOrder(t *testing.T) {
	type fields struct {
		Config config.Config
		API    *API
	}
	type args struct {
		args SendChildOrderRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    SendChildOrderResponse
		wantErr bool
	}{
		{
			name: "success with valid request",
			fields: fields{
				Config: testConfig,
				API:    NewAPI(),
			},
			args: args{
				args: SendChildOrderRequest{
					ProductCode:    consts.ProductCodeBTCJPY,
					ChildOrderType: consts.ChildOrderTypeLimit,
					Side:           consts.SideBuy,
					Price:          1000000,
					Size:           0.01,
					MinuteToExpire: 43200,
					TimeInForce:    consts.TimeInForceGTC,
				},
			},
			want:    SendChildOrderResponse{},
			wantErr: false,
		},
		{
			name: "success with market order",
			fields: fields{
				Config: testConfig,
				API:    NewAPI(),
			},
			args: args{
				args: SendChildOrderRequest{
					ProductCode:    consts.ProductCodeBTCJPY,
					ChildOrderType: consts.ChildOrderTypeMarket,
					Side:           consts.SideSell,
					Size:           0.01,
					MinuteToExpire: 43200,
					TimeInForce:    consts.TimeInForceIOC,
				},
			},
			want:    SendChildOrderResponse{},
			wantErr: false,
		},
		{
			name: "success with different product code",
			fields: fields{
				Config: testConfig,
				API:    NewAPI(),
			},
			args: args{
				args: SendChildOrderRequest{
					ProductCode:    consts.ProductCodeETHJPY,
					ChildOrderType: consts.ChildOrderTypeLimit,
					Side:           consts.SideBuy,
					Price:          500000,
					Size:           0.1,
					MinuteToExpire: 43200,
					TimeInForce:    consts.TimeInForceFOK,
				},
			},
			want:    SendChildOrderResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BitFlyerAPI{
				Config: tt.fields.Config,
				API:    tt.fields.API,
			}

			isDry := true // falseにすると本当に注文APIが実行されるので注意

			got, err := b.SendChildOrder(tt.args.args, isDry)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyerAPI.SendChildOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitFlyerAPI.SendChildOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitFlyerAPI_privateRequestHeader(t *testing.T) {
	type fields struct {
		Config config.Config
		API    *API
	}
	type args struct {
		timeStamp string
		method    string
		url       string
		body      []byte
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         bool
		checkHeaderFunc func(header map[string]any) bool
	}{
		{
			name: "success with GET request",
			fields: fields{
				Config: testConfig,
				API:    NewAPI(),
			},
			args: args{
				timeStamp: "1640995200",
				method:    "GET",
				url:       "https://api.bitflyer.com/v1/getticker",
				body:      nil,
			},
			wantErr: false,
			checkHeaderFunc: func(header map[string]any) bool {
				return header["ACCESS-KEY"] == testConfig.BitFlyer.ApiKey &&
					header["ACCESS-TIMESTAMP"] == "1640995200" &&
					header["ACCESS-SIGN"] != nil &&
					header["ACCESS-SIGN"] != ""
			},
		},
		{
			name: "success with POST request with body",
			fields: fields{
				Config: testConfig,
				API:    NewAPI(),
			},
			args: args{
				timeStamp: "1640995200",
				method:    "POST",
				url:       "https://api.bitflyer.com/v1/me/sendchildorder",
				body:      []byte(`{"product_code":"BTC_JPY","child_order_type":"LIMIT","side":"BUY","price":1000000,"size":0.01}`),
			},
			wantErr: false,
			checkHeaderFunc: func(header map[string]any) bool {
				return header["ACCESS-KEY"] == testConfig.BitFlyer.ApiKey &&
					header["ACCESS-TIMESTAMP"] == "1640995200" &&
					header["ACCESS-SIGN"] != nil &&
					header["ACCESS-SIGN"] != ""
			},
		},
		{
			name: "success with empty body",
			fields: fields{
				Config: testConfig,
				API:    NewAPI(),
			},
			args: args{
				timeStamp: "1640995200",
				method:    "POST",
				url:       "https://api.bitflyer.com/v1/me/getbalance",
				body:      []byte{},
			},
			wantErr: false,
			checkHeaderFunc: func(header map[string]any) bool {
				return header["ACCESS-KEY"] == testConfig.BitFlyer.ApiKey &&
					header["ACCESS-TIMESTAMP"] == "1640995200" &&
					header["ACCESS-SIGN"] != nil &&
					header["ACCESS-SIGN"] != ""
			},
		},
		{
			name: "error with invalid URL",
			fields: fields{
				Config: testConfig,
				API:    NewAPI(),
			},
			args: args{
				timeStamp: "1640995200",
				method:    "GET",
				url:       "://invalid-url",
				body:      nil,
			},
			wantErr: true,
			checkHeaderFunc: func(header map[string]any) bool {
				return header == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &BitFlyerAPI{
				Config: tt.fields.Config,
				API:    tt.fields.API,
			}
			got, err := api.privateRequestHeader(tt.args.timeStamp, tt.args.method, tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyerAPI.privateRequestHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.checkHeaderFunc(got) {
				t.Errorf("BitFlyerAPI.privateRequestHeader() = %v", got)
			}
		})
	}
}

func Test_extractPath(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success with full URL",
			args: args{
				u: "https://api.bitflyer.com/v1/getticker",
			},
			want:    "/v1/getticker",
			wantErr: false,
		},
		{
			name: "success with query parameters",
			args: args{
				u: "https://api.bitflyer.com/v1/getticker?product_code=BTC_JPY",
			},
			want:    "/v1/getticker",
			wantErr: false,
		},
		{
			name: "success with root path",
			args: args{
				u: "https://api.bitflyer.com/",
			},
			want:    "/",
			wantErr: false,
		},
		{
			name: "success with no path",
			args: args{
				u: "https://api.bitflyer.com",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "success with nested path",
			args: args{
				u: "https://api.bitflyer.com/v1/me/sendchildorder",
			},
			want:    "/v1/me/sendchildorder",
			wantErr: false,
		},
		{
			name: "success with localhost",
			args: args{
				u: "http://localhost:8080/api/test",
			},
			want:    "/api/test",
			wantErr: false,
		},
		{
			name: "error with invalid URL",
			args: args{
				u: "://invalid-url",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "error with malformed URL",
			args: args{
				u: "http://[::1:80",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "success with empty string",
			args: args{
				u: "",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractPath(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
