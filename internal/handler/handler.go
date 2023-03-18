package handler

import (
	"errors"

	"github.com/milkyskies/line-chatgpt/internal/bot"
	"github.com/milkyskies/line-chatgpt/internal/chat"
	"github.com/milkyskies/line-chatgpt/internal/database"
)

type MessageHandler struct {
	ChatServices *chat.ChatServices
	BotServices  *bot.BotServices
	Database	*database.Database
}

func NewMessageHandler(chatServices *chat.ChatServices, botServices *bot.BotServices, database *database.Database) *MessageHandler {
	return &MessageHandler{
		ChatServices: chatServices,
		BotServices:  botServices,
		Database: database,
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

	msg := database.NewMessage(userID, userID, messageText)
	mh.Database.PostMessage(msg)

	history, err := mh.Database.GetMessages(userID)
	if err != nil {
		return errors.New("failed to get history")
	}

	reply, err := botService.GenerateReply(messageText, history)
	if err != nil {
		return errors.New("failed to generate reply")
	}

	rpl := database.NewMessage("chatgpt", userID, reply)
	mh.Database.PostMessage(rpl)

	if err := chatService.SendMessage(userID, reply); err != nil {
		return errors.New("failed to send message")
	}

	return nil
}