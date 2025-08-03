package usecase

import (
	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
	"errors"
	"net/http"
	"testing"
)

func TestNewLineUsecase(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    ILineUsecase
		wantErr bool
	}{
		{
			name: "正常なConfigでLineUsecaseが作成される",
			args: args{
				cfg: config.Config{
					Line: config.Line{
						ChannelToken:  "valid_token",
						ChannelSecret: "valid_secret",
						GroupID:       "valid_group_id",
					},
				},
			},
			want:    nil, // 実際の実装では具体的なインスタンスと比較するのは困難なため、nilチェックのみ行う
			wantErr: false,
		},
		{
			name: "空のChannelTokenでエラーが返される",
			args: args{
				cfg: config.Config{
					Line: config.Line{
						ChannelToken:  "",
						ChannelSecret: "valid_secret",
						GroupID:       "valid_group_id",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "空のChannelSecretでエラーが返される",
			args: args{
				cfg: config.Config{
					Line: config.Line{
						ChannelToken:  "valid_token",
						ChannelSecret: "",
						GroupID:       "valid_group_id",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ChannelTokenとChannelSecretが両方空でエラーが返される",
			args: args{
				cfg: config.Config{
					Line: config.Line{
						ChannelToken:  "",
						ChannelSecret: "",
						GroupID:       "valid_group_id",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLineUsecase(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLineUsecase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NewLineUsecase() = nil, want non-nil ILineUsecase")
			}
			if tt.wantErr && got != nil {
				t.Errorf("NewLineUsecase() = %v, want nil", got)
			}
		})
	}
}

// MockLineAPI はテスト用のLineAPIモック
type MockLineAPI struct {
	PostMessageFunc func(message string) error
}

func (m *MockLineAPI) PostMessage(message string) error {
	if m.PostMessageFunc != nil {
		return m.PostMessageFunc(message)
	}
	return nil
}

func TestLineUsecase_SendMessageToGroup(t *testing.T) {
	type fields struct {
		Config   config.Config
		ILineAPI api.ILineAPI
	}
	type args struct {
		dto PostLineMessageDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "正常なメッセージ送信",
			fields: fields{
				Config: config.Config{
					Line: config.Line{
						ChannelToken:  "valid_token",
						ChannelSecret: "valid_secret",
						GroupID:       "valid_group_id",
					},
				},
				ILineAPI: &MockLineAPI{
					PostMessageFunc: func(message string) error {
						return nil
					},
				},
			},
			args: args{
				dto: PostLineMessageDTO{
					Message: "テストメッセージ",
				},
			},
			want:    http.StatusOK,
			wantErr: false,
		},
		{
			name: "空のメッセージでエラーが返される",
			fields: fields{
				Config: config.Config{
					Line: config.Line{
						ChannelToken:  "valid_token",
						ChannelSecret: "valid_secret",
						GroupID:       "valid_group_id",
					},
				},
				ILineAPI: &MockLineAPI{
					PostMessageFunc: func(message string) error {
						return nil
					},
				},
			},
			args: args{
				dto: PostLineMessageDTO{
					Message: "",
				},
			},
			want:    http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "LineAPI.PostMessageでエラーが発生した場合",
			fields: fields{
				Config: config.Config{
					Line: config.Line{
						ChannelToken:  "valid_token",
						ChannelSecret: "valid_secret",
						GroupID:       "valid_group_id",
					},
				},
				ILineAPI: &MockLineAPI{
					PostMessageFunc: func(message string) error {
						return errors.New("API error")
					},
				},
			},
			args: args{
				dto: PostLineMessageDTO{
					Message: "テストメッセージ",
				},
			},
			want:    http.StatusInternalServerError,
			wantErr: true,
		},
		{
			name: "長いメッセージの送信",
			fields: fields{
				Config: config.Config{
					Line: config.Line{
						ChannelToken:  "valid_token",
						ChannelSecret: "valid_secret",
						GroupID:       "valid_group_id",
					},
				},
				ILineAPI: &MockLineAPI{
					PostMessageFunc: func(message string) error {
						return nil
					},
				},
			},
			args: args{
				dto: PostLineMessageDTO{
					Message: "これは非常に長いテストメッセージです。" +
						"このメッセージは複数の文を含んでおり、" +
						"実際のアプリケーションで使用される可能性のある" +
						"長いメッセージをシミュレートしています。",
				},
			},
			want:    http.StatusOK,
			wantErr: false,
		},
		{
			name: "特殊文字を含むメッセージの送信",
			fields: fields{
				Config: config.Config{
					Line: config.Line{
						ChannelToken:  "valid_token",
						ChannelSecret: "valid_secret",
						GroupID:       "valid_group_id",
					},
				},
				ILineAPI: &MockLineAPI{
					PostMessageFunc: func(message string) error {
						return nil
					},
				},
			},
			args: args{
				dto: PostLineMessageDTO{
					Message: "特殊文字テスト: !@#$%^&*()_+-=[]{}|;':\",./<>?",
				},
			},
			want:    http.StatusOK,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LineUsecase{
				Config:   tt.fields.Config,
				ILineAPI: tt.fields.ILineAPI,
			}
			got, err := l.SendMessageToGroup(tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("LineUsecase.SendMessageToGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LineUsecase.SendMessageToGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
