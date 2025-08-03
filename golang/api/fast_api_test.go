package api

import (
	"bitcoin-app-golang/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewFastAPI(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name string
		args args
		want IFastAPI
	}{
		{
			name: "正常系 - FastAPI インスタンス作成",
			args: args{
				cfg: config.Config{
					ServerURL: config.ServerURL{
						FastAPIServer: "http://localhost:8000",
					},
				},
			},
			want: &FastAPI{
				Config: config.Config{
					ServerURL: config.ServerURL{
						FastAPIServer: "http://localhost:8000",
					},
				},
				API: NewAPI(),
			},
		},
		{
			name: "正常系 - 空のConfig",
			args: args{
				cfg: config.Config{},
			},
			want: &FastAPI{
				Config: config.Config{},
				API:    NewAPI(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFastAPI(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFastAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFastAPI_GetHealthcheck(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.serverFunc()
			defer server.Close()

			cfg := tt.fields.Config
			cfg.ServerURL.FastAPIServer = server.URL

			f := &FastAPI{
				Config: cfg,
				API:    tt.fields.API,
			}
			if err := f.GetHealthcheck(); (err != nil) != tt.wantErr {
				t.Errorf("FastAPI.GetHealthcheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
