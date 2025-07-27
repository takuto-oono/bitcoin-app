package usecase

import (
	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
	"net/http"
	"reflect"
	"testing"
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

func TestBitFlyerUsecase_BuyOrder(t *testing.T) {
	type fields struct {
		Config      config.Config
		BitFlyerAPI api.IBitFlyerAPI
	}
	type args struct {
		dto BuyOrderDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    api.SendChildOrderResponse
		want1   int
		wantErr bool
	}{
		// 注意: IsDryをfalseにすると実際の購入APIが実行されます
		{
			name: "success - LIMIT order",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: BuyOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: ChildOrderTypeLimit,
					Price:          1000000,
					Size:           0.001,
					MinuteToExpire: 43200,
					TimeInForce:    TimeInForceGTC,
					IsDry:          true, // 注意: falseにすると実際の購入APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{ChildOrderAcceptanceID: ""},
			want1:   http.StatusOK,
			wantErr: false,
		},
		{
			name: "success - MARKET order",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: BuyOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: ChildOrderTypeMarket,
					Price:          0,
					Size:           0.001,
					MinuteToExpire: 1,
					TimeInForce:    TimeInForceIOC,
					IsDry:          true, // 注意: falseにすると実際の購入APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{ChildOrderAcceptanceID: ""},
			want1:   http.StatusOK,
			wantErr: false,
		},
		{
			name: "invalid child order type",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: BuyOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: "INVALID",
					Price:          1000000,
					Size:           0.001,
					MinuteToExpire: 43200,
					TimeInForce:    TimeInForceGTC,
					IsDry:          true, // 注意: falseにすると実際の購入APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{},
			want1:   http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "invalid time in force",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: BuyOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: ChildOrderTypeLimit,
					Price:          1000000,
					Size:           0.001,
					MinuteToExpire: 43200,
					TimeInForce:    "INVALID",
					IsDry:          true, // 注意: falseにすると実際の購入APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{},
			want1:   http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "invalid minute to expire - too small",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: BuyOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: ChildOrderTypeLimit,
					Price:          1000000,
					Size:           0.001,
					MinuteToExpire: 0,
					TimeInForce:    TimeInForceGTC,
					IsDry:          true, // 注意: falseにすると実際の購入APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{},
			want1:   http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "invalid minute to expire - too large",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: BuyOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: ChildOrderTypeLimit,
					Price:          1000000,
					Size:           0.001,
					MinuteToExpire: 50000,
					TimeInForce:    TimeInForceGTC,
					IsDry:          true, // 注意: falseにすると実際の購入APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{},
			want1:   http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "invalid price for LIMIT order",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: BuyOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: ChildOrderTypeLimit,
					Price:          0,
					Size:           0.001,
					MinuteToExpire: 43200,
					TimeInForce:    TimeInForceGTC,
					IsDry:          true, // 注意: falseにすると実際の購入APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{},
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
			got, got1, err := b.BuyOrder(tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyerUsecase.BuyOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitFlyerUsecase.BuyOrder() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BitFlyerUsecase.BuyOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBitFlyerUsecase_SellOrder(t *testing.T) {
	type fields struct {
		Config      config.Config
		BitFlyerAPI api.IBitFlyerAPI
	}
	type args struct {
		dto SellOrderDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    api.SendChildOrderResponse
		want1   int
		wantErr bool
	}{
		// 注意: IsDryをfalseにすると実際の売却APIが実行されます
		{
			name: "success - LIMIT order",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: SellOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: ChildOrderTypeLimit,
					Price:          1000000,
					Size:           0.001,
					MinuteToExpire: 43200,
					TimeInForce:    TimeInForceGTC,
					IsDry:          true, // 注意: falseにすると実際の売却APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{ChildOrderAcceptanceID: ""},
			want1:   http.StatusOK,
			wantErr: false,
		},
		{
			name: "success - MARKET order",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: SellOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: ChildOrderTypeMarket,
					Price:          0,
					Size:           0.001,
					MinuteToExpire: 1,
					TimeInForce:    TimeInForceIOC,
					IsDry:          true, // 注意: falseにすると実際の売却APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{ChildOrderAcceptanceID: ""},
			want1:   http.StatusOK,
			wantErr: false,
		},
		{
			name: "invalid child order type",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: SellOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: "INVALID",
					Price:          1000000,
					Size:           0.001,
					MinuteToExpire: 43200,
					TimeInForce:    TimeInForceGTC,
					IsDry:          true, // 注意: falseにすると実際の売却APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{},
			want1:   http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "invalid time in force",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: SellOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: ChildOrderTypeLimit,
					Price:          1000000,
					Size:           0.001,
					MinuteToExpire: 43200,
					TimeInForce:    "INVALID",
					IsDry:          true, // 注意: falseにすると実際の売却APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{},
			want1:   http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "invalid price for LIMIT order",
			fields: fields{
				Config:      TestConfig,
				BitFlyerAPI: api.NewBitFlyerAPI(TestConfig),
			},
			args: args{
				dto: SellOrderDTO{
					ProductCode:    "BTC_JPY",
					ChildOrderType: ChildOrderTypeLimit,
					Price:          0,
					Size:           0.001,
					MinuteToExpire: 43200,
					TimeInForce:    TimeInForceGTC,
					IsDry:          true, // 注意: falseにすると実際の売却APIが実行されます
				},
			},
			want:    api.SendChildOrderResponse{},
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
			got, got1, err := b.SellOrder(tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyerUsecase.SellOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitFlyerUsecase.SellOrder() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BitFlyerUsecase.SellOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestChildOrderType_validate(t *testing.T) {
	tests := []struct {
		name    string
		c       ChildOrderType
		wantErr bool
	}{
		{
			name:    "valid LIMIT order type",
			c:       ChildOrderTypeLimit,
			wantErr: false,
		},
		{
			name:    "valid MARKET order type",
			c:       ChildOrderTypeMarket,
			wantErr: false,
		},
		{
			name:    "invalid order type",
			c:       "INVALID",
			wantErr: true,
		},
		{
			name:    "empty order type",
			c:       "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.validate(); (err != nil) != tt.wantErr {
				t.Errorf("ChildOrderType.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimeInForce_validate(t *testing.T) {
	tests := []struct {
		name    string
		tr      TimeInForce
		wantErr bool
	}{
		{
			name:    "valid GTC time in force",
			tr:      TimeInForceGTC,
			wantErr: false,
		},
		{
			name:    "valid IOC time in force",
			tr:      TimeInForceIOC,
			wantErr: false,
		},
		{
			name:    "valid FOK time in force",
			tr:      TimeInForceFOK,
			wantErr: false,
		},
		{
			name:    "invalid time in force",
			tr:      "INVALID",
			wantErr: true,
		},
		{
			name:    "empty time in force",
			tr:      "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.validate(); (err != nil) != tt.wantErr {
				t.Errorf("TimeInForce.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMinuteToExpire_validate(t *testing.T) {
	tests := []struct {
		name    string
		m       MinuteToExpire
		wantErr bool
	}{
		{
			name:    "valid minimum minute to expire",
			m:       MinuteToExpire(MinMinuteToExpire),
			wantErr: false,
		},
		{
			name:    "valid maximum minute to expire",
			m:       MinuteToExpire(MaxMinuteToExpire),
			wantErr: false,
		},
		{
			name:    "valid middle value minute to expire",
			m:       MinuteToExpire(1440), // 1 day
			wantErr: false,
		},
		{
			name:    "invalid minute to expire - too small",
			m:       MinuteToExpire(0),
			wantErr: true,
		},
		{
			name:    "invalid minute to expire - negative",
			m:       MinuteToExpire(-1),
			wantErr: true,
		},
		{
			name:    "invalid minute to expire - too large",
			m:       MinuteToExpire(MaxMinuteToExpire + 1),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.validate(); (err != nil) != tt.wantErr {
				t.Errorf("MinuteToExpire.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductCode_validate(t *testing.T) {
	tests := []struct {
		name    string
		p       ProductCode
		wantErr bool
	}{
		{
			name:    "valid BTC_JPY",
			p:       ProductCodeBTCJPY,
			wantErr: false,
		},
		{
			name:    "valid XRP_JPY",
			p:       ProductCodeXRPJPY,
			wantErr: false,
		},
		{
			name:    "valid ETH_JPY",
			p:       ProductCodeETHJPY,
			wantErr: false,
		},
		{
			name:    "valid XLM_JPY",
			p:       ProductCodeXLMJPY,
			wantErr: false,
		},
		{
			name:    "valid MONA_JPY",
			p:       ProductCodeMONAJPY,
			wantErr: false,
		},
		{
			name:    "valid ETH_BTC",
			p:       ProductCodeETHBTC,
			wantErr: false,
		},
		{
			name:    "valid BCH_BTC",
			p:       ProductCodeBCHBTC,
			wantErr: false,
		},
		{
			name:    "valid FX_BTC_JPY",
			p:       ProductCodeFXBTCJPY,
			wantErr: false,
		},
		{
			name:    "invalid product code",
			p:       ProductCode("INVALID_CODE"),
			wantErr: true,
		},
		{
			name:    "empty product code",
			p:       ProductCode(""),
			wantErr: true,
		},
		{
			name:    "lowercase product code",
			p:       ProductCode("btc_jpy"),
			wantErr: true,
		},
		{
			name:    "partial product code",
			p:       ProductCode("BTC"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.validate(); (err != nil) != tt.wantErr {
				t.Errorf("ProductCode.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewProductCode(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    ProductCode
		wantErr bool
	}{
		{
			name: "valid BTC_JPY",
			args: args{
				code: "BTC_JPY",
			},
			want:    ProductCodeBTCJPY,
			wantErr: false,
		},
		{
			name: "valid XRP_JPY",
			args: args{
				code: "XRP_JPY",
			},
			want:    ProductCodeXRPJPY,
			wantErr: false,
		},
		{
			name: "valid ETH_JPY",
			args: args{
				code: "ETH_JPY",
			},
			want:    ProductCodeETHJPY,
			wantErr: false,
		},
		{
			name: "valid XLM_JPY",
			args: args{
				code: "XLM_JPY",
			},
			want:    ProductCodeXLMJPY,
			wantErr: false,
		},
		{
			name: "valid MONA_JPY",
			args: args{
				code: "MONA_JPY",
			},
			want:    ProductCodeMONAJPY,
			wantErr: false,
		},
		{
			name: "valid ETH_BTC",
			args: args{
				code: "ETH_BTC",
			},
			want:    ProductCodeETHBTC,
			wantErr: false,
		},
		{
			name: "valid BCH_BTC",
			args: args{
				code: "BCH_BTC",
			},
			want:    ProductCodeBCHBTC,
			wantErr: false,
		},
		{
			name: "valid FX_BTC_JPY",
			args: args{
				code: "FX_BTC_JPY",
			},
			want:    ProductCodeFXBTCJPY,
			wantErr: false,
		},
		{
			name: "invalid product code",
			args: args{
				code: "INVALID_CODE",
			},
			want:    ProductCode(""),
			wantErr: true,
		},
		{
			name: "empty product code",
			args: args{
				code: "",
			},
			want:    ProductCode(""),
			wantErr: true,
		},
		{
			name: "lowercase product code",
			args: args{
				code: "btc_jpy",
			},
			want:    ProductCode(""),
			wantErr: true,
		},
		{
			name: "partial product code",
			args: args{
				code: "BTC",
			},
			want:    ProductCode(""),
			wantErr: true,
		},
		{
			name: "product code with spaces",
			args: args{
				code: " BTC_JPY ",
			},
			want:    ProductCode(""),
			wantErr: true,
		},
		{
			name: "mixed case product code",
			args: args{
				code: "Btc_Jpy",
			},
			want:    ProductCode(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProductCode(tt.args.code)
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
