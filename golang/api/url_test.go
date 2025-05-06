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
