package database

import (
	"fmt"
	"time"

	"github.com/surrealdb/surrealdb.go"
)

type Message struct {
	SenderID    string
	RoomID      string
	MessageText string
	SentAt      time.Time
}

func NewMessage(senderID string, roomID string, messageText string) Message {
	return Message{
		SenderID:    senderID,
		RoomID:      roomID,
		MessageText: messageText,
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

func (db *Database) GetMessages(roomID string) ([]Message, error) {
	sql := fmt.Sprintf("SELECT * FROM messages WHERE RoomID = '%s' ORDER BY SentAt DESC LIMIT 15", roomID)

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

	for i := 0; i < len(messages)/2; i++ {
		j := len(messages) - i - 1
		messages[i], messages[j] = messages[j], messages[i]
	}

	// for _, message := range messages {
	// 	truncatedMessage := message.MessageText
	// 	if len(truncatedMessage) > 50 {
	// 		truncatedMessage = truncatedMessage[:50] + "..."
	// 	}
	// 	fmt.Printf("(%s) %s %s: %s\n", message.RoomID, message.SentAt, message.SenderID, truncatedMessage)
	// }

	return messages, nil
}

type MessagesResponse struct {
	Messages []Message `json:"result"`
	Status   string    `json:"status"`
	Time     string    `json:"time"`
}
