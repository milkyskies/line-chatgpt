package chatgpt

import (
	"context"
	"errors"
	"strings"

	"github.com/milkyskies/line-chatgpt/internal/database"
	openai "github.com/sashabaranov/go-openai"
)

var (
	ErrReplyGenerationFailed = errors.New("reply generation failed")
)

func (c *ChatGPT) GenerateReply(prompt string, history []database.Message) (string, error) {
	resp, err := c.createChatCompletion(prompt, history)
	if err != nil {
		return "", err
	}

	response := resp.Choices[0].Message.Content
	trimmedResponse := strings.TrimLeft(response, "\n")

	return trimmedResponse, nil
}

func messageToChatCompletionMessages(messages []database.Message) []openai.ChatCompletionMessage {
	var chatCompletionMessages []openai.ChatCompletionMessage
	for _, message := range messages {
		var role string
		if message.SenderID == "chatgpt" {
			role = openai.ChatMessageRoleAssistant
		} else {
			role = openai.ChatMessageRoleUser
		}
		chatCompletionMessage := openai.ChatCompletionMessage{
			Role:    role,
			Content: message.MessageText,
		}
		chatCompletionMessages = append(chatCompletionMessages, chatCompletionMessage)
	}
	return chatCompletionMessages
}

func (c *ChatGPT) createChatCompletion(prompt string, history []database.Message) (openai.ChatCompletionResponse, error) {
	messageHistory := messageToChatCompletionMessages(history)
	messageHistory = append(messageHistory, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	return c.OpenAI.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4,
			Messages: messageHistory,
		},
	)
}
