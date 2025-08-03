package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"

	"bitcoin-app-golang/config"
	"bitcoin-app-golang/usecase"
)

type ILineHandler interface {
	PostMessage(ctx *gin.Context)
	CallbackMessage(ctx *gin.Context)
}

type LineHandler struct {
	Config       config.Config
	ILineUsecase usecase.ILineUsecase
}

func NewLineHandler(cfg config.Config) (ILineHandler, error) {
	lineUsecase, err := usecase.NewLineUsecase(cfg)
	if err != nil {
		return nil, err
	}

	return &LineHandler{
		Config:       cfg,
		ILineUsecase: lineUsecase,
	}, nil
}

func (h *LineHandler) PostMessage(ctx *gin.Context) {
	var dto usecase.PostLineMessageDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	statusCode, err := h.ILineUsecase.SendMessageToGroup(dto)
	if err != nil {
		ctx.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(statusCode, gin.H{"status": "Message sent successfully"})
}

func (h *LineHandler) CallbackMessage(c *gin.Context) {
	// 本来Usecase層で処理するべきだが、利用されない想定のコードなのでここに残しておく。

	bot, err := linebot.New(string(h.Config.Line.ChannelSecret), string(h.Config.Line.ChannelToken))
	if err != nil {
		log.Println("Error creating Line bot:", err)
		c.String(http.StatusInternalServerError, "Internal server error")
		return
	}

	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.String(http.StatusBadRequest, "Invalid signature")
		} else {
			c.String(http.StatusInternalServerError, "Parse error")
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			if _, ok := event.Message.(*linebot.TextMessage); ok {
				if event.Source.Type == linebot.EventSourceTypeGroup {
					groupID := event.Source.GroupID
					replyText := "このグループのIDは: " + groupID

					_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do()
					if err != nil {
						log.Println("Reply error:", err)
					}
				} else {
					// グループ以外のケース（任意対応）
					_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("グループ内で使ってください")).Do()
					if err != nil {
						log.Println("Reply error:", err)
					}
				}
			}
		}
	}

	c.Status(http.StatusOK)
}
