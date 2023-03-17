package chat

import (
	"errors"
)

var (
    ErrServiceNotFound = errors.New("chat service not found")
    ErrServiceExists = errors.New("chat service already registered")
)

type ChatServiceName int

const (
	LineChat ChatServiceName = iota
	CustomChat
)

type ChatServices struct {
	services map[ChatServiceName]Chat
}

func NewChatServices() *ChatServices {
	return &ChatServices{
		services: make(map[ChatServiceName]Chat),
	}
}

func (cs *ChatServices) Register(name ChatServiceName, service Chat) error {
	if _, exists := cs.services[name]; exists {
		return ErrServiceExists
	}
	cs.services[name] = service
	return nil
}

func (cs *ChatServices) Get(name ChatServiceName) (Chat, error) {
	service, exists := cs.services[name]
	if !exists {
		return nil, ErrServiceNotFound
	}
	return service, nil
}
