package chat

type Chat interface {
    SendMessage(userID string, message string) error
    SendAudioMessage(userID string, id string) error
    ReceiveMessage(userID string) (string, error)
}
