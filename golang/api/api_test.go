package api

import (
	"reflect"
	"testing"

	"bitcoin-app-golang/config"
)

var testConfig config.Config

func init() {
	var err error
	testConfig, err = config.NewConfig("../toml/local.toml", "../env/.env.test")
	if err != nil {
		panic(err)
	}
}

func Test_marshalJson(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Test marshalJson",
			args: args{
				v: map[string]any{
					"key": "value",
				},
			},
			want:    []byte("{\"key\":\"value\"}"),
			wantErr: false,
		},
		{
			name: "Test marshalJson is nil",
			args: args{
				v: nil,
			},
			want:    []byte{},
			wantErr: false,
		},
		{
			name: "Test marshalJson is empty",
			args: args{
				v: "",
			},
			want:    []byte{},
			wantErr: false,
		},
		{
			name: "Test marshalJson is emptm map",
			args: args{
				v: map[string]any{},
			},
			want:    []byte("{}"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshalJson(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshalJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("marshalJson() = %v, want %v", got, tt.want)
			}
		})
	}
}
