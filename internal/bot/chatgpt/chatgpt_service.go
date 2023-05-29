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
	messages := messageToChatCompletionMessages(history)
	messages = append(
		messages,
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	)

	systemMessage := "For this conversation, your name is Chatty. Respond to that name please. Also respond in the language that the last message from the user was in. "

	messagesWithSystemMessage := make([]openai.ChatCompletionMessage, len(messages)+1)
	messagesWithSystemMessage[0] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemMessage,
	}
	copy(messagesWithSystemMessage[1:], messages)

	return c.OpenAI.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messagesWithSystemMessage,
		},
	)
}
