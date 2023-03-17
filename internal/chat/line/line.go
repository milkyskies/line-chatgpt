package line

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/milkyskies/line-chatgpt/internal/chat"
)

type LineChat struct {
    Client *linebot.Client
}

func NewLineChat(channelSecret, channelAccessToken string) (*LineChat, error) {
    bot, err := linebot.New(channelSecret, channelAccessToken)
    if err != nil {
        return nil, err
    }
    return &LineChat{Client: bot}, nil
}

var _ chat.Chat = (*LineChat)(nil)