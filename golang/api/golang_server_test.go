package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"bitcoin-app-golang/config"
	"bitcoin-app-golang/consts"
)

func TestGolangServerAPI_GetBitFlyerTicker(t *testing.T) {
	// モックティッカーデータ
	mockTicker := TickerFromGolangServer{
		TickID:          12345,
		ProductCode:     consts.ProductCodeBTCJPY,
		State:           "RUNNING",
		Timestamp:       "2025-05-18T17:00:00",
		BestBid:         5000000.0,
		BestAsk:         5010000.0,
		BestBidSize:     0.1,
		BestAskSize:     0.2,
		TotalBidDepth:   1000.0,
		TotalAskDepth:   1200.0,
		MarketBidSize:   0.5,
		MarketAskSize:   0.6,
		Ltp:             5005000.0,
		Volume:          100.0,
		VolumeByProduct: 90.0,
	}

	type fields struct {
		Config config.Config
		API    *API
	}
	type args struct {
		productCode string
	}

	tests := []struct {
		name       string
		serverFunc func() *httptest.Server
		fields     fields
		args       args
		want       TickerFromGolangServer
		wantErr    bool
	}{
		{
			name: "正常系 - BTC_JPY",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// リクエストパラメータを確認
					if r.URL.Path != "/bitflyer/ticker/" {
						t.Errorf("Expected path '/bitflyer/ticker/', got %s", r.URL.Path)
					}

					productCode := r.URL.Query().Get("product_code")
					if productCode != consts.ProductCodeBTCJPY {
						t.Errorf("Expected product_code 'BTC_JPY', got %s", productCode)
					}

					// 正常なレスポンスを返す
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(mockTicker)
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			args: args{
				productCode: consts.ProductCodeBTCJPY,
			},
			want:    mockTicker,
			wantErr: false,
		},
		{
			name: "正常系 - ETH_JPY",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// リクエストパラメータを確認
					if r.URL.Path != "/bitflyer/ticker/" {
						t.Errorf("Expected path '/bitflyer/ticker/', got %s", r.URL.Path)
					}

					productCode := r.URL.Query().Get("product_code")
					if productCode != consts.ProductCodeETHJPY {
						t.Errorf("Expected product_code 'ETH_JPY', got %s", productCode)
					}

					// 正常なレスポンスを返す
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(mockTicker)
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			args: args{
				productCode: consts.ProductCodeETHJPY,
			},
			want:    mockTicker,
			wantErr: false,
		},
		{
			name: "異常系 - サーバーエラー",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error"))
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			args: args{
				productCode: consts.ProductCodeBTCJPY,
			},
			want:    TickerFromGolangServer{},
			wantErr: true,
		},
		{
			name: "異常系 - 無効なProductCode",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// 無効なProductCodeの場合でもサーバーは正常に動作するが、
					// GolangServerAPI.GetBitFlyerTickerメソッド内でエラーが発生する
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode("Invalid product code")
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			args: args{
				productCode: "INVALID_CODE", // 無効なProductCode
			},
			want:    TickerFromGolangServer{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用のサーバーを作成
			server := tt.serverFunc()
			defer server.Close()

			// Configを設定
			cfg := tt.fields.Config
			cfg.ServerURL.GolangServer = server.URL

			g := &GolangServerAPI{
				Config: cfg,
				API:    tt.fields.API,
			}
			got, err := g.GetBitFlyerTicker(tt.args.productCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GolangServerAPI.GetBitFlyerTicker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GolangServerAPI.GetBitFlyerTicker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGolangServerAPI_GetHealthcheck(t *testing.T) {
	type fields struct {
		Config config.Config
		API    *API
	}
	tests := []struct {
		name       string
		serverFunc func() *httptest.Server
		fields     fields
		wantErr    bool
	}{
		{
			name: "正常系 - ヘルスチェック成功",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/healthcheck/" {
						t.Errorf("Expected path '/healthcheck/', got %s", r.URL.Path)
					}

					if r.Method != http.MethodGet {
						t.Errorf("Expected method GET, got %s", r.Method)
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			wantErr: false,
		},
		{
			name: "異常系 - サーバーエラー",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error"))
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			wantErr: true,
		},
		{
			name: "異常系 - サービス利用不可",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Write([]byte("Service Unavailable"))
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			wantErr: true,
		},
		{
			name: "異常系 - 404 Not Found",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("Not Found"))
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			wantErr: true,
		},
		{
			name: "異常系 - タイムアウト",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusRequestTimeout)
					w.Write([]byte("Request Timeout"))
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.serverFunc()
			defer server.Close()

			cfg := tt.fields.Config
			cfg.ServerURL.GolangServer = server.URL

			g := &GolangServerAPI{
				Config: cfg,
				API:    tt.fields.API,
			}
			if err := g.GetHealthcheck(); (err != nil) != tt.wantErr {
				t.Errorf("GolangServerAPI.GetHealthcheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
