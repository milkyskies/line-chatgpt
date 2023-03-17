package chatgpt

import (
	openai "github.com/sashabaranov/go-openai"

	"github.com/milkyskies/line-chatgpt/internal/bot"
)

type ChatGPT struct {
	Client *openai.Client
}

func NewChatGPT(apiKey string) *ChatGPT {
    client := openai.NewClient(apiKey)
    return &ChatGPT{Client: client}
}

var _ bot.Bot = (*ChatGPT)(nil)