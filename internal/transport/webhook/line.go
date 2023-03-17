package webhook

import (
	"errors"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/milkyskies/line-chatgpt/internal/bot"
	"github.com/milkyskies/line-chatgpt/internal/chat"
	"github.com/milkyskies/line-chatgpt/internal/chat/line"
	"github.com/milkyskies/line-chatgpt/internal/handler"
)

var (
	ErrInvalidMessage = errors.New("invalid message")
)

type LineWebhookHandler struct {
	LineChat          *line.LineChat
	MessageHandler *handler.MessageHandler
}

func NewLineWebhookHandler(lineChat *line.LineChat, messageHandler *handler.MessageHandler) *LineWebhookHandler {
	return &LineWebhookHandler{LineChat: lineChat, MessageHandler: messageHandler}
}

func (lwh *LineWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	events, err := lwh.LineChat.Client.ParseRequest(r)
	if err != nil {
		http.Error(w, "Failed to parse request", http.StatusBadRequest)
		return
	}

	lwh.handleEvents(events)
}

func (lwh *LineWebhookHandler) handleEvents(events []*linebot.Event) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			lwh.handleMessageEvent(event)
		}
	}
}

func (lwh *LineWebhookHandler) handleMessageEvent(event *linebot.Event) error {
	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		return ErrInvalidMessage
	}


	if err := lwh.MessageHandler.HandleMessage(chat.LineChat, bot.ChatGPT, event.Source.UserID, message.Text); err != nil {
		return err
	}

	return nil
}