package bot

import "errors"

var (
	ErrServiceNotFound = errors.New("bot service not found")
	ErrServiceExists   = errors.New("bot service already registered")
)

type ServiceName int

const (
	ChatGPT ServiceName = iota
)

type Services struct {
	services map[ServiceName]Bot
}

func NewBotServices() *Services {
	return &Services{
		services: make(map[ServiceName]Bot),
	}
}

func (bs *Services) Register(name ServiceName, service Bot) error {
	if _, exists := bs.services[name]; exists {
		return ErrServiceExists
	}

	bs.services[name] = service
	return nil
}

func (bs *Services) Get(name ServiceName) (Bot, error) {
	service, ok := bs.services[name]
	if !ok {
		return nil, ErrServiceNotFound
	}
	return service, nil
}
