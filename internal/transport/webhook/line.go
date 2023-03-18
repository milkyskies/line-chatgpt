package webhook

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/milkyskies/line-chatgpt/internal/bot"
	"github.com/milkyskies/line-chatgpt/internal/bot/chatgpt"
	"github.com/milkyskies/line-chatgpt/internal/chat"
	"github.com/milkyskies/line-chatgpt/internal/chat/line"
	"github.com/milkyskies/line-chatgpt/internal/handler"
	"github.com/sashabaranov/go-openai"
)

var (
	ErrInvalidMessage = errors.New("invalid message")
)

// TODO: remove chatgpt from here
type LineWebhookHandler struct {
	LineChat       *line.LineChat
	ChatGPT        *chatgpt.ChatGPT
	MessageHandler *handler.MessageHandler
}

func NewLineWebhookHandler(lineChat *line.LineChat,  chatGPT *chatgpt.ChatGPT, messageHandler *handler.MessageHandler) *LineWebhookHandler {
	return &LineWebhookHandler{LineChat: lineChat, ChatGPT: chatGPT, MessageHandler: messageHandler}
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
	switch event.Message.(type) {
	case *linebot.TextMessage:
		return lwh.handleTextMessageEvent(event)
	case *linebot.AudioMessage:
		return lwh.handleAudioMessageEvent(event)
	}

	return nil
}

func (lwh *LineWebhookHandler) handleTextMessageEvent(event *linebot.Event) error {
	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		return ErrInvalidMessage
	}

	if err := lwh.MessageHandler.HandleMessage(chat.LineChat, bot.ChatGPT, event.Source.UserID, message.Text); err != nil {
		return err
	}

	return nil
}

func (lwh *LineWebhookHandler) handleAudioMessageEvent(event *linebot.Event) error {
	message, ok := event.Message.(*linebot.AudioMessage)
	if !ok {
		return ErrInvalidMessage
	}

	content, err := lwh.LineChat.Client.GetMessageContent(message.ID).Do()
	if err != nil {
		return err
	}
	//defer content.Content.Close()

	if err := saveAsM4A(content.Content, fmt.Sprintf("%s.m4a", message.ID)); err != nil {
		return err
	}

	req := openai.AudioRequest{
		FilePath: filepath.Join("content/line/audio", fmt.Sprintf("%s.m4a", message.ID)),
		Model:    openai.Whisper1,
	}

	ctx := context.Background()
	res, err := lwh.ChatGPT.Client.CreateTranscription(ctx, req)
	if err != nil {
		return err
	}

	if err := lwh.MessageHandler.HandleAudioMessage(chat.LineChat, bot.ChatGPT, event.Source.UserID, res.Text); err != nil {
		return err
	}

	return nil
}

// TODO: move this
func saveAsM4A(r io.ReadCloser, fileName string) error {
	outputDir := "content/line/audio"

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}
	outputFilePath := filepath.Join(outputDir, fileName)
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Copy data from the ReadCloser to the output file
	_, err = io.Copy(outputFile, r)
	if err != nil {
		return err
	}

	// Close the ReadCloser
	if err := r.Close(); err != nil {
		return err
	}

	return nil
}