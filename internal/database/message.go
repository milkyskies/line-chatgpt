package database

import (
	"fmt"
	"time"

	"github.com/surrealdb/surrealdb.go"
	//uuid "github.com/satori/go.uuid"
)

type Message struct {
	SenderID    string
	RoomID      string
	MessageText string
	SentAt      time.Time
}

func NewMessage(SenderID string, RoomID string, MessageText string) Message {
	return Message{
		SenderID:    SenderID,
		RoomID:      RoomID,
		MessageText: MessageText,
		SentAt:      time.Now(),
	}
}

func (db *Database) PostMessage(msg Message) error {
	_, err := db.Client.Create("messages", msg)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetMessages(roomId string) ([]Message, error) {
	sql := fmt.Sprintf("SELECT * FROM messages WHERE RoomID = '%s'", roomId)

	res, err := db.Client.Query(sql, nil)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := surrealdb.Unmarshal(res, &response); err != nil {
		return nil, err
	}

	var messages []Message
	if err := surrealdb.Unmarshal(response.Result, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

type MessagesResponse struct {
	Messages []Message `json:"result"`
	Status string   `json:"status"`
	Time   string   `json:"time"`
}