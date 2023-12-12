package eventbus

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

type EventHandler interface {
	Handle(*discordgo.Session, interface{})
}

type EventBus struct {
	mux  sync.RWMutex
	subs map[string][]EventHandler
}

func New() *EventBus {
	return &EventBus{
		subs: make(map[string][]EventHandler),
	}
}

func (eb *EventBus) Subscribe(e string, callback EventHandler) {
	eb.mux.Lock()
	defer eb.mux.Unlock()
	eb.subs[e] = append(eb.subs[e], callback)
}

func (eb *EventBus) Publish(e string, s *discordgo.Session, i interface{}) {
	eb.mux.RLock()
	defer eb.mux.RUnlock()

	handlers, ok := eb.subs[e]
	if !ok {
		return
	}

	for _, callback := range handlers {
		go func(callback EventHandler) {
			callback.Handle(s, i)
		}(callback)
	}
}
