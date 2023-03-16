package linebot

import (
	"fmt"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/milkyskies/line-chatgpt/internal/chatgpt"
)

type LineBot struct {
	Client *linebot.Client
}

func NewBot() (*LineBot, error) {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	bot, err := linebot.New(channelSecret, channelToken)

	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &LineBot{bot}, nil
}

func (b *LineBot) HandleRequest(r *http.Request) error {
	events, err := b.Client.ParseRequest(r)
	if err != nil {
		return fmt.Errorf("failed to parse request: %w", err)
	}

	return b.handleEvents(events)
}

func (b *LineBot) handleEvents(events []*linebot.Event) error {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			if err := b.handleMessageEvent(event); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *LineBot) handleMessageEvent(event *linebot.Event) error {
	message, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		return nil
	}

	openai, err := chatgpt.NewChatGPT()
	if err != nil {
		return fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	openaiResponse, err := openai.GetResponse(message.Text)
	if err != nil {
		return fmt.Errorf("could not get response: %w", err)
	}

	newMessage := linebot.NewTextMessage(openaiResponse)
	_, err = b.Client.PushMessage(event.Source.UserID, newMessage).Do()

	return err
}