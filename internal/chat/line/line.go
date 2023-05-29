package line

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/milkyskies/line-chatgpt/internal/chat"
	"github.com/milkyskies/line-chatgpt/internal/speech"
)

type Chat struct {
	Client      *linebot.Client
	FileManager *speech.FileManager
}

func NewLineChat(channelSecret, channelAccessToken string) (*Chat, error) {
	bot, err := linebot.New(channelSecret, channelAccessToken)
	if err != nil {
		return nil, err
	}
	return &Chat{Client: bot}, nil
}

var _ chat.Chat = (*Chat)(nil)
