package handler

import (
	"errors"

	"github.com/milkyskies/line-chatgpt/internal/bot"
	"github.com/milkyskies/line-chatgpt/internal/chat"
)

type MessageHandler struct {
	ChatServices *chat.ChatServices
	BotServices  *bot.BotServices
}

func NewMessageHandler(chatServices *chat.ChatServices, botServices *bot.BotServices) *MessageHandler {
	return &MessageHandler{
		ChatServices: chatServices,
		BotServices:  botServices,
	}
}

func (mh *MessageHandler) HandleMessage(chatServiceName chat.ChatServiceName, botServiceName bot.BotServiceName, userID, messageText string) error {
	chatService, err := mh.ChatServices.Get(chatServiceName)
	if err != nil {
		return errors.New("chat service not found")
	}

	botService, err := mh.BotServices.Get(botServiceName)
	if err != nil {
		return errors.New("bot service not found")
	}

	reply, err := botService.GenerateReply(messageText)
	if err != nil {
		return errors.New("failed to generate reply")
	}

	if err := chatService.SendMessage(userID, reply); err != nil {
		return errors.New("failed to send message")
	}

	return nil
}