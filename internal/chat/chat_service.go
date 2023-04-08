package chat

import (
	"errors"
)

var (
	ErrServiceNotFound = errors.New("chat service not found")
	ErrServiceExists   = errors.New("chat service already registered")
)

type ServiceName int

const (
	LineChat ServiceName = iota
	CustomChat
)

type Services struct {
	services map[ServiceName]Chat
}

func NewChatServices() *Services {
	return &Services{
		services: make(map[ServiceName]Chat),
	}
}

func (cs *Services) Register(name ServiceName, service Chat) error {
	if _, exists := cs.services[name]; exists {
		return ErrServiceExists
	}
	cs.services[name] = service
	return nil
}

func (cs *Services) Get(name ServiceName) (Chat, error) {
	service, exists := cs.services[name]
	if !exists {
		return nil, ErrServiceNotFound
	}
	return service, nil
}
