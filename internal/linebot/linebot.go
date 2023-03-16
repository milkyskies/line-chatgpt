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
		return nil, fmt.Errorf("Failed to create bot: %w", err)
	}

	return &LineBot{bot}, nil
}

func (b *LineBot) HandleRequest(r *http.Request) error {
	prompt := "Hi! How are you doing?"

	openai, err := chatgpt.NewOpenAI()
	if err != nil {
		return fmt.Errorf("Failed to create OpenAI client: %w", err)
	}

	response, err := openai.GetResponse(prompt)

	message := linebot.NewTextMessage(response)
	b.Client.BroadcastMessage(message).Do()

	events, err := b.Client.ParseRequest(r)
	if err != nil {
		return fmt.Errorf("Failed to parse request: %w", err)
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			leftBtn := linebot.NewMessageAction("left", "left clicked")
			rightBtn := linebot.NewMessageAction("right", "right clicked")

			template := linebot.NewConfirmTemplate("Hello World", leftBtn, rightBtn)

			message := linebot.NewTemplateMessage("Sorry :(, please update your app.", template)
			fmt.Println("hello")
			b.Client.BroadcastMessage(message)

			//b.Client.PushMessage(event.Source.UserID, message)
		}
	}

	return nil
}
