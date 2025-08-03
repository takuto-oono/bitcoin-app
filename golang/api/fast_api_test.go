package api

import (
	"bitcoin-app-golang/config"
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
		// TODO: Add test cases.
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
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FastAPI{
				Config: tt.fields.Config,
				API:    tt.fields.API,
			}
			if err := f.GetHealthcheck(); (err != nil) != tt.wantErr {
				t.Errorf("FastAPI.GetHealthcheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
