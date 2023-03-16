package chatgpt

import (
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	Client *openai.Client
}

func NewOpenAI() (*OpenAI, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	return &OpenAI{client}, nil
}

func (c *OpenAI) GetResponse(prompt string) (string, error) {
	resp, err := c.Client.CreateChatCompletion(
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

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	response := resp.Choices[0].Message.Content

	trimmedResponse := strings.TrimLeft(response, "\n")

	fmt.Println(trimmedResponse)
	return trimmedResponse, nil
}
