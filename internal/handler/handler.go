package handler

import (
	"errors"
	"fmt"

	"github.com/milkyskies/line-chatgpt/internal/bot"
	"github.com/milkyskies/line-chatgpt/internal/chat"
	"github.com/milkyskies/line-chatgpt/internal/database"
)

type MessageHandler struct {
	ChatServices *chat.Services
	BotServices  *bot.Services
	Database     *database.Database
}

func NewMessageHandler(chatServices *chat.Services, botServices *bot.Services, database *database.Database) *MessageHandler {
	return &MessageHandler{
		ChatServices: chatServices,
		BotServices:  botServices,
		Database:     database,
	}
}

func (mh *MessageHandler) HandleMessage(chatServiceName chat.ServiceName, botServiceName bot.ServiceName, userID, messageText string) error {
	chatService, err := mh.ChatServices.Get(chatServiceName)
	if err != nil {
		return errors.New("chat service not found")
	}

	botService, err := mh.BotServices.Get(botServiceName)
	if err != nil {
		return errors.New("bot service not found")
	}

	msg := database.NewMessage(userID, userID, messageText)
	if err := mh.Database.PostMessage(msg); err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	history, err := mh.Database.GetMessages(userID)
	if err != nil {
		return errors.New("failed to get history")
	}

	reply, err := botService.GenerateReply(messageText, history)
	if err != nil {
		return errors.New("failed to generate reply")
	}

	rpl := database.NewMessage("chatgpt", userID, reply)
	if err := mh.Database.PostMessage(rpl); err != nil {
		return fmt.Errorf("failed to post reply: %w", err)
	}

	if err := chatService.SendMessage(userID, reply); err != nil {
		return errors.New("failed to send message")
	}

	return nil
}

// TODO: CLEAN THIS UP
func (mh *MessageHandler) HandleAudioMessage(chatServiceName chat.ServiceName, botServiceName bot.ServiceName, userID, messageText string) (string, error) {
	_, err := mh.ChatServices.Get(chatServiceName)
	if err != nil {
		return "", errors.New("chat service not found")
	}

	botService, err := mh.BotServices.Get(botServiceName)
	if err != nil {
		return "", errors.New("bot service not found")
	}

	msg := database.NewMessage(userID, userID, messageText)
	if err := mh.Database.PostMessage(msg); err != nil {
		return "", fmt.Errorf("failed to post message: %w", err)
	}

	history, err := mh.Database.GetMessages(userID)
	if err != nil {
		return "", errors.New("failed to get history")
	}

	reply, err := botService.GenerateReply(messageText, history)
	if err != nil {
		return "", errors.New("failed to generate reply")
	}

	rpl := database.NewMessage("chatgpt", userID, reply)
	err = mh.Database.PostMessage(rpl)
	if err != nil {
		return "", fmt.Errorf("failed to post reply: %w", err)
	}

	return reply, nil
}
