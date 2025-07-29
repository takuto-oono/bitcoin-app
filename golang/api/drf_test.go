package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"bitcoin-app-golang/config"
	"bitcoin-app-golang/consts"
)

func TestDRFAPI_PostBitFlyerTicker(t *testing.T) {
	// モックティッカーデータ
	mockTicker := PostTickerDRFRequest{
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

					if receivedTicker.ProductCode != consts.ProductCodeBTCJPY {
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

					if receivedTicker.ProductCode != consts.ProductCodeETHJPY {
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
					ProductCode:     consts.ProductCodeETHJPY,
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

func TestDRFAPI_GetBitFlyerTickers(t *testing.T) {
	// モックレスポンスデータ
	mockTickers := []GetTickerFromDRFResponse{
		{
			ID:              1,
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
		},
		{
			ID:              2,
			TickID:          12346,
			ProductCode:     consts.ProductCodeETHJPY,
			State:           "RUNNING",
			Timestamp:       "2025-05-18T17:01:00",
			BestBid:         300000.0,
			BestAsk:         301000.0,
			BestBidSize:     0.5,
			BestAskSize:     0.3,
			TotalBidDepth:   800.0,
			TotalAskDepth:   900.0,
			MarketBidSize:   0.4,
			MarketAskSize:   0.7,
			Ltp:             300500.0,
			Volume:          50.0,
			VolumeByProduct: 45.0,
		},
	}

	type fields struct {
		Config config.Config
		API    *API
	}
	tests := []struct {
		name       string
		serverFunc func() *httptest.Server
		fields     fields
		want       []GetTickerFromDRFResponse
		wantErr    bool
	}{
		{
			name: "正常系 - 複数のティッカーを取得",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// リクエストパスを確認
					if r.URL.Path != "/api/tickers/" {
						t.Errorf("Expected path '/api/tickers/', got %s", r.URL.Path)
					}

					// リクエストメソッドを確認
					if r.Method != http.MethodGet {
						t.Errorf("Expected method GET, got %s", r.Method)
					}

					// 正常なレスポンスを返す
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(mockTickers)
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			want:    mockTickers,
			wantErr: false,
		},
		{
			name: "正常系 - 空のティッカーリスト",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// リクエストパスを確認
					if r.URL.Path != "/api/tickers/" {
						t.Errorf("Expected path '/api/tickers/', got %s", r.URL.Path)
					}

					// リクエストメソッドを確認
					if r.Method != http.MethodGet {
						t.Errorf("Expected method GET, got %s", r.Method)
					}

					// 空のレスポンスを返す
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode([]GetTickerFromDRFResponse{})
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			want:    []GetTickerFromDRFResponse{},
			wantErr: false,
		},
		{
			name: "正常系 - 単一のティッカーを取得",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// リクエストパスを確認
					if r.URL.Path != "/api/tickers/" {
						t.Errorf("Expected path '/api/tickers/', got %s", r.URL.Path)
					}

					// リクエストメソッドを確認
					if r.Method != http.MethodGet {
						t.Errorf("Expected method GET, got %s", r.Method)
					}

					// 単一のティッカーレスポンスを返す
					singleTicker := []GetTickerFromDRFResponse{mockTickers[0]}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(singleTicker)
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			want:    []GetTickerFromDRFResponse{mockTickers[0]},
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
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系 - 不正なJSONレスポンス",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("invalid json"))
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			want:    nil,
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
			want:    nil,
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
			got, err := d.GetBitFlyerTickers()
			if (err != nil) != tt.wantErr {
				t.Errorf("DRFAPI.GetBitFlyerTickers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DRFAPI.GetBitFlyerTickers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDRFAPI_DeleteBitFlyerTicker(t *testing.T) {
	type fields struct {
		Config config.Config
		API    *API
	}
	type args struct {
		id int
	}
	tests := []struct {
		name       string
		serverFunc func() *httptest.Server
		fields     fields
		args       args
		wantErr    bool
	}{
		{
			name: "正常系 - Ticker削除",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodDelete {
						t.Errorf("Expected method DELETE, got %s", r.Method)
					}
					w.WriteHeader(http.StatusNoContent)
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			args: args{
				id: 1,
			},
			wantErr: false,
		},
		{
			name: "異常系 - 無効なID",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Invalid ID"))
				}))
			},
			fields: fields{
				API: NewAPI(),
			},
			args: args{
				id: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.serverFunc()
			defer server.Close()

			cfg := tt.fields.Config
			cfg.ServerURL.DRFServer = server.URL

			d := &DRFAPI{
				Config: cfg,
				API:    tt.fields.API,
			}
			if err := d.DeleteBitFlyerTicker(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DRFAPI.DeleteBitFlyerTicker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
