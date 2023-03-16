package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"github.com/milkyskies/line-chatgpt/internal/messenger"
	"github.com/milkyskies/line-chatgpt/internal/chatgpt"
	transportHttp "github.com/milkyskies/line-chatgpt/internal/transport/http"
)

func Run() error {
	fmt.Println("Starting LINE ChatGPT Bot")
	
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return err
	}

	// rename to services later
	gptService := chatgpt.NewChatGPT()

	msgnService, err := messenger.NewLineBot(*gptService) 
	if err != nil {
		fmt.Println("Error creating LINE bot")
		return err
	}


	httpHandler, err := transportHttp.NewHandler(*msgnService)
	if err != nil {
		fmt.Println("Error creating HTTP handler")
		return err
	}

	if err := httpHandler.Serve(); err != nil {
		return err
	}

	// bot, err := linebot.NewBot()
	// if err != nil {
	// 	fmt.Println("Error creating LINE bot")
	// 	return err
	// }

	return nil
}

func main() {
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
