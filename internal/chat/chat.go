package chat

type Chat interface {
    SendMessage(userID string, message string) error
    ReceiveMessage(userID string) (string, error)
}
