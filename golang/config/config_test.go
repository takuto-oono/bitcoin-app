package config

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
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
				GeneralSetting: GeneralSetting{
					Port: "8080",
				},
			},
			wantErr: false,
		},
		{
			name: "success prod toml",
			args: args{
				tomlFilePath: "../toml/prod.toml",
			},
			want: Config{
				GeneralSetting: GeneralSetting{
					Port: "7080",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.tomlFilePath)
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

func TestConfig_mustCheck(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "success",
			config: &Config{
				GeneralSetting: GeneralSetting{
					Port: "8080",
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
			name: "fail port is empty",
			config: &Config{
				GeneralSetting: GeneralSetting{
					Port: "",
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
