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
	"github.com/milkyskies/line-chatgpt/internal/database"
	"github.com/milkyskies/line-chatgpt/internal/handler"
	"github.com/milkyskies/line-chatgpt/internal/openai"
	"github.com/milkyskies/line-chatgpt/internal/transport/http"
	"github.com/milkyskies/line-chatgpt/internal/transport/webhook"
)

func Run() error {
	fmt.Println("starting LINE ChatGPT Bot")

	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading .env file")
		return err
	}

	chatServices := chat.NewChatServices()
	botServices := bot.NewBotServices()

	openai := openai.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	chatGPT := chatgpt.NewChatGPT(openai)
	if err = botServices.Register(bot.ChatGPT, chatGPT); err != nil {
		return fmt.Errorf("failed to register ChatGPT bot: %v", err)
	}

	lineChat, err := line.NewLineChat(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		return fmt.Errorf("failed to initialize LINE chat: %v", err)
	}
	if err := chatServices.Register(chat.LineChat, lineChat); err != nil {
		return fmt.Errorf("failed to register LINE chat: %v", err)
	}

	database, err := database.NewDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.Client.Close()

	messageHandler := handler.NewMessageHandler(chatServices, botServices, database)
	lineWebhookHandler := webhook.NewLineWebhookHandler(lineChat, chatGPT, messageHandler)

	// initDB := flag.Bool("init-db", false, "initialize the database")
	// flag.Parse()
	// if *initDB {
	// 	if err := database.Init(os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD")); err != nil {
	// 		log.Fatalf("failed to initialize database: %v", err)
	// 	}
	// }
	if err := database.Init(os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD")); err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	httpHandler, err := http.NewHandler(lineWebhookHandler)
	if err != nil {
		return fmt.Errorf("failed to initialize HTTP handler: %v", err)
	}

	if err := httpHandler.Serve(); err != nil {
		return fmt.Errorf("failed to serve HTTP handler: %v", err)
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
