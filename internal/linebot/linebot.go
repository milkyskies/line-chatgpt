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
	// prompt := "Hi! How are you doing?"

	openai, err := chatgpt.NewOpenAI()
	if err != nil {
		return fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	// message := linebot.NewTextMessage(response)
	// b.Client.BroadcastMessage(message).Do()

	events, err := b.Client.ParseRequest(r)
	if err != nil {
		return fmt.Errorf("failed to parse request: %w", err)
	}

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				text := message.Text

				openaiResponse, err := openai.GetResponse(text)
				if err != nil {
					return fmt.Errorf("could not get response: %w", err)
				}

				newMessage := linebot.NewTextMessage(openaiResponse)

				b.Client.PushMessage(event.Source.UserID, newMessage).Do()
			}
		}
	}

	return nil
}
