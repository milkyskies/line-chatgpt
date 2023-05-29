package line

import (
	"errors"
	"fmt"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	ErrSendMessageFailed          = errors.New("failed to send message")
	ErrReceiveMessageNotSupported = errors.New("receive message not supported for LINE")
)

func (l *Chat) SendMessage(userID string, message string) error {
	_, err := l.Client.PushMessage(userID, linebot.NewTextMessage(message)).Do()
	if err != nil {
		return ErrSendMessageFailed
	}
	return nil
}

func (l *Chat) SendAudioMessage(userID string, id string) error {
	f, err := l.FileManager.OpenAudioFile(id)
	if err != nil {
		return err
	}
	defer f.Close()

	duration, err := l.FileManager.CalculateAudioDuration(f)
	if err != nil {
		return err
	}

	fmt.Println("Duration: ", duration)

	hostname := os.Getenv("HOSTNAME")

	audioMessage := linebot.NewAudioMessage(hostname+"/audio_replies/"+id, duration)

	if _, err := l.Client.PushMessage(userID, audioMessage).Do(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
