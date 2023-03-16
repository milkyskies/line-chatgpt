package main

import (
	"fmt"

	"github.com/joho/godotenv"

	//"github.com/milkyskies/line-chatgpt/internal/linebot"
	transportHttp "github.com/milkyskies/line-chatgpt/internal/transport/http"
)

func Run() error {
	fmt.Println("Starting LINE ChatGPT Bot")
	
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return err
	}

	httpHandler := transportHttp.NewHandler()
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
