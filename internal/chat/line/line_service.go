package line

import (
    "github.com/line/line-bot-sdk-go/v7/linebot"
    "errors"
)

var (
    ErrSendMessageFailed = errors.New("failed to send message")
    ErrReceiveMessageNotSupported = errors.New("receive message not supported for LINE")
)

func (l *LineChat) SendMessage(userID string, message string) error {
    _, err := l.Client.PushMessage(userID, linebot.NewTextMessage(message)).Do()
    if err != nil {
        return ErrSendMessageFailed
    }
    return nil
}

func (l *LineChat) SendAudioMessage(userID string, message string) error {
    _, err := l.Client.PushMessage(userID, linebot.NewTextMessage(message)).Do()
    if err != nil {
        return ErrSendMessageFailed
    }
    return nil
}

func (l *LineChat) ReceiveMessage(userID string) (string, error) {
    // Pass the message somewhere else to handle

    return "", ErrReceiveMessageNotSupported
}
