package api

import (
	"errors"
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"bitcoin-app-golang/config"
)

func NewLinebot(cfg config.Config) (*linebot.Client, error) {
	if cfg.Line.ChannelToken == "" || cfg.Line.ChannelSecret == "" {
		return nil, errors.New("line channel token or secret is empty")
	}

	return linebot.New(string(cfg.Line.ChannelSecret), string(cfg.Line.ChannelToken))
}

type ILineAPI interface {
	PostMessage(message string) error
}

type LineAPI struct {
	Config config.Config
	Bot    *linebot.Client
}

func NewLineAPI(cfg config.Config) (ILineAPI, error) {
	bot, err := NewLinebot(cfg)
	if err != nil {
		return nil, err
	}

	return &LineAPI{
		Config: cfg,
		Bot:    bot,
	}, nil
}

func (l *LineAPI) PostMessage(message string) error {
	if l.Config.Line.GroupID == "" {
		return errors.New("line group ID is empty")
	}

	if message == "" {
		return errors.New("message is empty")
	}

	if l.Bot == nil {
		return errors.New("line bot client is not initialized")
	}

	res, err := l.Bot.PushMessage(string(l.Config.Line.GroupID), linebot.NewTextMessage(message)).Do()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	log.Printf("Message sent successfully: %v", res)
	return nil
}
