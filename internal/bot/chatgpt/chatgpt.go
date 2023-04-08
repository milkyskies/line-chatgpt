package chatgpt

import (
	"github.com/milkyskies/line-chatgpt/internal/bot"
	"github.com/milkyskies/line-chatgpt/internal/openai"
)

type ChatGPT struct {
	OpenAI *openai.OpenAI
}

func NewChatGPT(openai *openai.OpenAI) *ChatGPT {
	return &ChatGPT{OpenAI: openai}
}

var _ bot.Bot = (*ChatGPT)(nil)
