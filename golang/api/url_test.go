package api

import (
	"net/url"
	"testing"

	"bitcoin-app-golang/consts"
)

func TestBitFlyerURL_GetTicker(t *testing.T) {
	type args struct {
		productCode string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{productCode: consts.ProductCodeBTCJPY},
			want:    "https://api.bitflyer.com/v1/getticker/?product_code=BTC_JPY",
			wantErr: false,
		},
		{
			name: "success productCode is empty",
			args: args{
				productCode: "",
			},
			want:    "https://api.bitflyer.com/v1/getticker/",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BitFlyerURL(BitFlyerBaseURL).GetTicker(tt.args.productCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitFlyerURL.GetTicker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BitFlyerURL.GetTicker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitFlyerURL_SendChildOrder(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			want:    "https://api.bitflyer.com/v1/me/sendchildorder/",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BitFlyerURL(BitFlyerBaseURL).SendChildOrder()
			if err != nil {
				t.Errorf("BitFlyerURL.SendChildOrder() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("BitFlyerURL.SendChildOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGolangServerURL_GetTicker(t *testing.T) {
	type args struct {
		productCode string
	}
	tests := []struct {
		name    string
		g       GolangServerURL
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			g:    GolangServerURL("https://localhost:8080"),
			args: args{
				productCode: consts.ProductCodeBTCJPY,
			},
			want:    "https://localhost:8080/bitflyer/ticker/?product_code=BTC_JPY",
			wantErr: false,
		},
		{
			name: "success productCode is empty",
			g:    GolangServerURL("https://localhost:8080"),
			args: args{
				productCode: "",
			},
			want:    "https://localhost:8080/bitflyer/ticker/",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.GetTicker(tt.args.productCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GolangServerURL.GetTicker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GolangServerURL.GetTicker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDRFServerURL_PostTicker(t *testing.T) {
	tests := []struct {
		name    string
		d       DRFServerURL
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			d:       DRFServerURL("https://localhost:8080"),
			want:    "https://localhost:8080/api/tickers/",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.PostTicker()
			if (err != nil) != tt.wantErr {
				t.Errorf("DRFServerURL.PostTicker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DRFServerURL.PostTicker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createUrl(t *testing.T) {
	type args struct {
		baseUrl string
		p       string
		qVal    url.Values
		el      []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				baseUrl: "https://localhost:8080",
				p:       "",
				qVal:    url.Values{},
			},
			want:    "https://localhost:8080/",
			wantErr: false,
		},
		{
			name: "success p and qVal",
			args: args{
				baseUrl: "https://localhost:8080",
				p:       "hoge/path",
				qVal: url.Values{
					"hoge": []string{"fuga"},
				},
			},
			want:    "https://localhost:8080/hoge/path/?hoge=fuga",
			wantErr: false,
		},
		{
			name: "baseUrl is empty",
			args: args{
				baseUrl: "",
				p:       "",
				qVal:    url.Values{},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createUrl(tt.args.baseUrl, tt.args.p, tt.args.qVal, tt.args.el...)
			if (err != nil) != tt.wantErr {
				t.Errorf("createUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_withSuffixSlash(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				s: "https://localhost:8080",
			},
			want: "https://localhost:8080/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := withSuffixSlash(tt.args.s); got != tt.want {
				t.Errorf("withSuffixSlash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractPort(t *testing.T) {
	type args struct {
		urlString string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				urlString: "https://localhost:8080",
			},
			want:    "8080",
			wantErr: false,
		},
		{
			name: "url is empty",
			args: args{
				urlString: "",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid url",
			args: args{
				urlString: "hogehogehogehoge",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "port is empty",
			args: args{
				urlString: "https://localhost",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractPort(tt.args.urlString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDRFServerURL_GetTickers(t *testing.T) {
	tests := []struct {
		name    string
		g       DRFServerURL
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			g:       DRFServerURL("https://localhost:8080"),
			want:    "https://localhost:8080/api/tickers/",
			wantErr: false,
		},
		{
			name:    "success with different port",
			g:       DRFServerURL("https://localhost:3000"),
			want:    "https://localhost:3000/api/tickers/",
			wantErr: false,
		},
		{
			name:    "success with different host",
			g:       DRFServerURL("https://api.example.com"),
			want:    "https://api.example.com/api/tickers/",
			wantErr: false,
		},
		{
			name:    "success with trailing slash",
			g:       DRFServerURL("https://localhost:8080/"),
			want:    "https://localhost:8080/api/tickers/",
			wantErr: false,
		},
		{
			name:    "error empty base URL",
			g:       DRFServerURL(""),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.GetTickers()
			if (err != nil) != tt.wantErr {
				t.Errorf("DRFServerURL.GetTickers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DRFServerURL.GetTickers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDRFServerURL_DeleteTicker(t *testing.T) {
	tests := []struct {
		name    string
		g       DRFServerURL
		id      int
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			g:       DRFServerURL("https://localhost:8080"),
			id:      123,
			want:    "https://localhost:8080/api/tickers/123/",
			wantErr: false,
		},
		{
			name:    "error invalid ID",
			g:       DRFServerURL("https://localhost:8080"),
			id:      -1,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.DeleteTicker(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DRFServerURL.DeleteTicker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DRFServerURL.DeleteTicker() = %v, want %v", got, tt.want)
			}
		})
	}
}
