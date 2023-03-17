package bot

import "errors"

var (
    ErrServiceNotFound = errors.New("bot service not found")
    ErrServiceExists = errors.New("bot service already registered")
)

type BotServiceName int

const (
    ChatGPT BotServiceName = iota
)

type BotServices struct {
    services map[BotServiceName]Bot
}

func NewBotServices() *BotServices {
    return &BotServices{
        services: make(map[BotServiceName]Bot),
    }
}

func (bs *BotServices) Register(name BotServiceName, service Bot) error {
    if _, exists := bs.services[name]; exists {
        return ErrServiceExists
    }
    
    bs.services[name] = service
    return nil
}

func (bs *BotServices) Get(name BotServiceName) (Bot, error) {
    service, ok := bs.services[name]
    if !ok {
        return nil, ErrServiceNotFound
    }
    return service, nil
}
