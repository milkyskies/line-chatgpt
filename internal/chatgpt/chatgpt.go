package chatgpt

import (
	"context"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type ChatGPT struct {
	Client *openai.Client
}

func NewChatGPT() (*ChatGPT) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	return &ChatGPT{client}
}

func (c *ChatGPT) GetResponse(prompt string) (string, error) {
	resp, err := c.createChatCompletion(prompt)
	if err != nil {
		return "", err
	}

	response := resp.Choices[0].Message.Content
	trimmedResponse := strings.TrimLeft(response, "\n")

	return trimmedResponse, nil
}

func (c *ChatGPT) createChatCompletion(prompt string) (openai.ChatCompletionResponse, error) {
	return c.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
}
