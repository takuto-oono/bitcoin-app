package api

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"bitcoin-app-golang/config"
)

func TestNewLinebot(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid config",
			args: args{
				cfg: config.Config{
					Line: config.Line{
						ChannelToken:  "test_token",
						ChannelSecret: "test_secret",
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "empty channel token",
			args: args{
				cfg: config.Config{
					Line: config.Line{
						ChannelToken:  "",
						ChannelSecret: "test_secret",
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "empty channel secret",
			args: args{
				cfg: config.Config{
					Line: config.Line{
						ChannelToken:  "test_token",
						ChannelSecret: "",
					},
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLinebot(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLinebot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.want {
				t.Errorf("NewLinebot() got = %v, want %v", got != nil, tt.want)
			}
		})
	}
}

func TestNewLineAPI(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid config",
			args: args{
				cfg: config.Config{
					Line: config.Line{
						ChannelToken:  "test_token",
						ChannelSecret: "test_secret",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty channel token",
			args: args{
				cfg: config.Config{
					Line: config.Line{
						ChannelToken:  "",
						ChannelSecret: "test_secret",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty channel secret",
			args: args{
				cfg: config.Config{
					Line: config.Line{
						ChannelToken:  "test_token",
						ChannelSecret: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "both empty",
			args: args{
				cfg: config.Config{
					Line: config.Line{
						ChannelToken:  "",
						ChannelSecret: "",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLineAPI(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLineAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("NewLineAPI() returned nil when expecting valid LineAPI")
					return
				}
				lineAPI, ok := got.(*LineAPI)
				if !ok {
					t.Errorf("NewLineAPI() returned wrong type, expected *LineAPI")
					return
				}
				if lineAPI.Bot == nil {
					t.Errorf("NewLineAPI() Bot is nil")
				}
				if !reflect.DeepEqual(lineAPI.Config, tt.args.cfg) {
					t.Errorf("NewLineAPI() Config = %v, want %v", lineAPI.Config, tt.args.cfg)
				}
			}
		})
	}
}

func TestLineAPI_PostMessage(t *testing.T) {
	cfg, err := config.NewConfig("../toml/local.toml", "../env/.env.local")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	bot, err := NewLinebot(cfg)
	if err != nil {
		t.Fatalf("Failed to create Line bot: %v", err)
	}

	sendingMessage := fmt.Sprintf("Golang Unit Test Message %s", time.Now().Format(time.RFC3339))

	type fields struct {
		Config config.Config
		Bot    *linebot.Client
	}
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful message post",
			fields: fields{
				Config: cfg,
				Bot:    bot,
			},
			args: args{
				message: sendingMessage,
			},
			wantErr: false,
		},
		{
			name: "empty message",
			fields: fields{
				Config: cfg,
				Bot:    bot,
			},
			args: args{
				message: "",
			},
			wantErr: true,
		},
		{
			name: "empty group ID",
			fields: fields{
				Config: config.Config{
					Line: config.Line{
						GroupID: "",
					},
				},
				Bot: bot,
			},
			args: args{
				message: "Test message",
			},
			wantErr: true,
		},
		{
			name: "uninitialized bot",
			fields: fields{
				Config: cfg,
				Bot:    nil,
			},
			args: args{
				message: "Test message",
			},
			wantErr: true,
		},
		{
			name: "invalid group ID",
			fields: fields{
				Config: config.Config{
					Line: config.Line{
						GroupID: "invalid_group_id",
					},
				},
				Bot: bot,
			},
			args: args{
				message: "Test message",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LineAPI{
				Config: tt.fields.Config,
				Bot:    tt.fields.Bot,
			}
			if err := l.PostMessage(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("LineAPI.PostMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
