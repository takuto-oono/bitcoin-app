package usecase

import (
	"errors"
	"net/http"

	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
)

type ILineUsecase interface {
	SendMessageToGroup(dto PostLineMessageDTO) (int, error)
}

type LineUsecase struct {
	Config   config.Config
	ILineAPI api.ILineAPI
}

type PostLineMessageDTO struct {
	Message string `json:"message"`
}

func NewLineUsecase(cfg config.Config) (ILineUsecase, error) {
	lineAPI, err := api.NewLineAPI(cfg)
	if err != nil {
		return nil, err
	}

	return &LineUsecase{
		Config:   cfg,
		ILineAPI: lineAPI,
	}, nil
}

func (l *LineUsecase) SendMessageToGroup(dto PostLineMessageDTO) (int, error) {
	if dto.Message == "" {
		return http.StatusBadRequest, errors.New("message cannot be empty")
	}

	if err := l.ILineAPI.PostMessage(dto.Message); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
