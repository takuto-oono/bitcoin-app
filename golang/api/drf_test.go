package api

import (
	"bitcoin-app-golang/config"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDRFAPI_PostBitFlyerTicker(t *testing.T) {
	// モックティッカーデータ
	mockTicker := PostTickerDRFRequest{
		TickID:          12345,
		ProductCode:     "BTC_JPY",
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
		ticker PostTickerDRFRequest
	}
	tests := []struct {
		name       string
		serverFunc func() *httptest.Server
		fields     fields
		args       args
		wantErr    bool
	}{
		{
			name: "正常系 - BTC_JPY",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// リクエストパスを確認
					if r.URL.Path != "/api/tickers/" {
						t.Errorf("Expected path '/api/tickers/', got %s", r.URL.Path)
					}

					// リクエストメソッドを確認
					if r.Method != http.MethodPost {
						t.Errorf("Expected method POST, got %s", r.Method)
					}

					// リクエストボディを確認
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Errorf("Failed to read request body: %v", err)
					}
					defer r.Body.Close()

					var receivedTicker PostTickerDRFRequest
					if err := json.Unmarshal(body, &receivedTicker); err != nil {
						t.Errorf("Failed to unmarshal request body: %v", err)
					}

					if receivedTicker.ProductCode != "BTC_JPY" {
						t.Errorf("Expected product_code 'BTC_JPY', got %s", receivedTicker.ProductCode)
					}

					// 正常なレスポンスを返す
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(map[string]string{"status": "success"})
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			args: args{
				ticker: mockTicker,
			},
			wantErr: false,
		},
		{
			name: "正常系 - ETH_JPY",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// リクエストパスを確認
					if r.URL.Path != "/api/tickers/" {
						t.Errorf("Expected path '/api/tickers/', got %s", r.URL.Path)
					}

					// リクエストメソッドを確認
					if r.Method != http.MethodPost {
						t.Errorf("Expected method POST, got %s", r.Method)
					}

					// リクエストボディを確認
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Errorf("Failed to read request body: %v", err)
					}
					defer r.Body.Close()

					var receivedTicker PostTickerDRFRequest
					if err := json.Unmarshal(body, &receivedTicker); err != nil {
						t.Errorf("Failed to unmarshal request body: %v", err)
					}

					if receivedTicker.ProductCode != "ETH_JPY" {
						t.Errorf("Expected product_code 'ETH_JPY', got %s", receivedTicker.ProductCode)
					}

					// 正常なレスポンスを返す
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					json.NewEncoder(w).Encode(map[string]string{"status": "success"})
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			args: args{
				ticker: PostTickerDRFRequest{
					TickID:          12345,
					ProductCode:     "ETH_JPY",
					State:           "RUNNING",
					Timestamp:       "2025-05-18T17:00:00",
					BestBid:         300000.0,
					BestAsk:         301000.0,
					BestBidSize:     0.1,
					BestAskSize:     0.2,
					TotalBidDepth:   1000.0,
					TotalAskDepth:   1200.0,
					MarketBidSize:   0.5,
					MarketAskSize:   0.6,
					Ltp:             300500.0,
					Volume:          100.0,
					VolumeByProduct: 90.0,
				},
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
			args: args{
				ticker: mockTicker,
			},
			wantErr: true,
		},
		{
			name: "異常系 - 無効なリクエスト",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			args: args{
				ticker: PostTickerDRFRequest{
					ProductCode: "INVALID_CODE",
					// 他のフィールドは省略
				},
			},
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
			cfg.ServerURL.DRFServer = server.URL

			d := &DRFAPI{
				Config: cfg,
				API:    tt.fields.API,
			}
			if err := d.PostBitFlyerTicker(tt.args.ticker); (err != nil) != tt.wantErr {
				t.Errorf("DRFAPI.PostBitFlyerTicker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
