package openai

import (
	openai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	Client *openai.Client
}

func NewOpenAI(apiKey string) *OpenAI {
	client := openai.NewClient(apiKey)
	return &OpenAI{Client: client}
}
