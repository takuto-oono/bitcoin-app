package config

import (
	"reflect"
	"testing"
)

const (
	TestBitFlyerAPIKey    = "BITFLYER_API_KEY_HOGE_HOGE"
	TestBitFlyerAPISecret = "BITFLYER_API_SECRET_HOGE_HOGE"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		tomlFilePath string
		envFilePath  string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "success local toml",
			args: args{
				tomlFilePath: "../toml/local.toml",
				envFilePath:  "../env/.env.test",
			},
			want: Config{
				ServerURL: ServerURL{
					GolangServer: "http://localhost:8080",
					DRFServer:    "http://localhost:8000",
				},
				BitFlyer: BitFlyer{
					ApiKey:    TestBitFlyerAPIKey,
					ApiSecret: TestBitFlyerAPISecret,
				},
			},
			wantErr: false,
		},
		{
			name: "success prod toml",
			args: args{
				tomlFilePath: "../toml/prod.toml",
				envFilePath:  "../env/.env.test",
			},
			want: Config{
				ServerURL: ServerURL{
					GolangServer: "http://localhost:7080",
					DRFServer:    "http://localhost:7000",
				},
				BitFlyer: BitFlyer{
					ApiKey:    TestBitFlyerAPIKey,
					ApiSecret: TestBitFlyerAPISecret,
				},
			},
			wantErr: false,
		},
		{
			name: "fail toml",
			args: args{
				tomlFilePath: "toml/fail.toml",
			},
			want:    Config{},
			wantErr: true,
		},
		{
			name: "fail env",
			args: args{
				tomlFilePath: "../toml/local.toml",
				envFilePath:  "env/fail.env",
			},
			want:    Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.tomlFilePath, tt.args.envFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_setFromToml(t *testing.T) {
	type args struct {
		tomlFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "success local toml",
			args: args{
				tomlFilePath: "../toml/local.toml",
			},
			want: Config{
				ServerURL: ServerURL{
					GolangServer: "http://localhost:8080",
					DRFServer:    "http://localhost:8000",
				},
				BitFlyer: BitFlyer{},
			},
			wantErr: false,
		},
		{
			name: "success prod toml",
			args: args{
				tomlFilePath: "../toml/prod.toml",
			},
			want: Config{
				ServerURL: ServerURL{
					GolangServer: "http://localhost:7080",
					DRFServer:    "http://localhost:7000",
				},
				BitFlyer: BitFlyer{},
			},
			wantErr: false,
		},
		{
			name: "fail toml",
			args: args{
				tomlFilePath: "toml/fail.toml",
			},
			want:    Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config
			if err := cfg.setFromToml(tt.args.tomlFilePath); (err != nil) != tt.wantErr {
				t.Errorf("Config.setFromToml() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(cfg, tt.want) {
				t.Errorf("Config.setFromToml() = %v, want %v", cfg, tt.want)
			}
		})
	}
}

func TestConfig_setFromEnv(t *testing.T) {
	type args struct {
		envFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				envFilePath: "../env/.env.test",
			},
			want: Config{
				ServerURL: ServerURL{},
				BitFlyer: BitFlyer{
					ApiKey:    TestBitFlyerAPIKey,
					ApiSecret: TestBitFlyerAPISecret,
				},
			},
			wantErr: false,
		},
		{
			name: "fail env",
			args: args{
				envFilePath: "env/fail.env",
			},
			want:    Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg Config
			if err := cfg.setFromEnv(tt.args.envFilePath); (err != nil) != tt.wantErr {
				t.Errorf("Config.setFromEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(cfg, tt.want) {
				t.Errorf("Config.setFromEnv() = %v, want %v", cfg, tt.want)
			}
		})
	}
}

func TestConfig_mustCheck(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "success",
			config: &Config{
				ServerURL: ServerURL{
					GolangServer: "http://localhost:8080",
					DRFServer:    "http://localhost:8000",
				},
				BitFlyer: BitFlyer{
					ApiKey:    TestBitFlyerAPIKey,
					ApiSecret: TestBitFlyerAPISecret,
				},
			},
			wantErr: false,
		},
		{
			name:    "fail config is nil",
			config:  nil,
			wantErr: true,
		},
		{
			name: "fail golang server is empty",
			config: &Config{
				ServerURL: ServerURL{
					GolangServer: "",
				},
				BitFlyer: BitFlyer{
					ApiKey:    TestBitFlyerAPIKey,
					ApiSecret: TestBitFlyerAPISecret,
				},
			},
			wantErr: true,
		},
		{
			name: "fail drf server is empty",
			config: &Config{
				ServerURL: ServerURL{
					GolangServer: "http://localhost:8080",
					DRFServer: "",
				},
				BitFlyer: BitFlyer{
					ApiKey:    TestBitFlyerAPIKey,
					ApiSecret: TestBitFlyerAPISecret,
				},
			},
			wantErr: true,
		},
		{
			name: "fail bitflyer api key is empty",
			config: &Config{
				ServerURL: ServerURL{
					GolangServer: "http://localhost:8080",
					DRFServer:    "http://localhost:8000",
				},
				BitFlyer: BitFlyer{
					ApiKey:    "",
					ApiSecret: TestBitFlyerAPISecret,
				},
			},
			wantErr: true,
		},
		{
			name: "fail bitflyer api secret is empty",
			config: &Config{
				ServerURL: ServerURL{
					GolangServer: "http://localhost:8080",
					DRFServer:    "http://localhost:8000",
				},
				BitFlyer: BitFlyer{
					ApiKey:    TestBitFlyerAPIKey,
					ApiSecret: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.config.mustCheck(); (err != nil) != tt.wantErr {
				t.Errorf("Config.mustCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
