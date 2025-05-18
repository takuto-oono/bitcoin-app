package api

import (
	"net/url"
	"testing"
)

func TestBitFlyerURL_GetTicker(t *testing.T) {
	type args struct {
		productCode ProductCode
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{productCode: "BTC_JPY"},
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

func TestGolangServerURL_GetTicker(t *testing.T) {
	type args struct {
		productCode ProductCode
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
				productCode: "BTC_JPY",
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
