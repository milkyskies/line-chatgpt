package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/milkyskies/line-chatgpt/internal/bot"
	"github.com/milkyskies/line-chatgpt/internal/bot/chatgpt"
	"github.com/milkyskies/line-chatgpt/internal/chat"
	"github.com/milkyskies/line-chatgpt/internal/chat/line"
	"github.com/milkyskies/line-chatgpt/internal/handler"
	"github.com/milkyskies/line-chatgpt/internal/transport/http"
	"github.com/milkyskies/line-chatgpt/internal/transport/webhook"
)

func Run() error {
	fmt.Println("Starting LINE ChatGPT Bot")
	
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return err
	}

	chatServices := chat.NewChatServices()
	botServices := bot.NewBotServices()

	chatGPT := chatgpt.NewChatGPT(os.Getenv("OPENAI_API_KEY"))
	botServices.Register(bot.ChatGPT, chatGPT)

	// Create Line instance and MessageHandler instance
	lineChat, _ := line.NewLineChat(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	chatServices.Register(chat.LineChat, lineChat)

	messageHandler := handler.NewMessageHandler(chatServices, botServices)

	// Register the LINE webhook handler
	lineWebhookHandler := webhook.NewLineWebhookHandler(lineChat, messageHandler)

	// Initialize the HTTP handler
	httpHandler, err := http.NewHandler(lineWebhookHandler)
	if err != nil {
		log.Fatalf("Failed to initialize HTTP handler: %v", err)
	}

	// Serve the HTTP handler
	if err := httpHandler.Serve(); err != nil {
		log.Fatalf("Failed to serve HTTP handler: %v", err)
	}


	return nil
}

func main() {
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
