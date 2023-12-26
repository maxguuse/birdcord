package pubsub

import (
	"fmt"
	"reflect"
	"sync"
)

type PubSub interface {
	Publish(topic string, args ...any)
	Close(topic string)
	Subscribe(topic string, fn any) error
	Unsubscribe(topic string, fn any) error
}

type handlersMap map[string][]*handler

type handler struct {
	callback reflect.Value
	queue    chan []reflect.Value
}

type messageBus struct {
	handlerQueueSize int
	mtx              sync.RWMutex
	handlers         handlersMap
}

func New(handlerQueueSize int) func() *messageBus {
	return func() *messageBus {
		if handlerQueueSize == 0 {
			panic("handlerQueueSize has to be greater then 0")
		}

		return &messageBus{
			handlerQueueSize: handlerQueueSize,
			handlers:         make(handlersMap),
		}
	}
}

func (b *messageBus) Publish(topic string, args ...any) {
	rArgs := buildHandlerArgs(args)

	b.mtx.RLock()
	defer b.mtx.RUnlock()

	if hs, ok := b.handlers[topic]; ok {
		for _, h := range hs {
			h.queue <- rArgs
		}
	}
}

func (b *messageBus) Subscribe(topic string, fn any) error {
	if err := isValidHandler(fn); err != nil {
		return err
	}

	h := &handler{
		callback: reflect.ValueOf(fn),
		queue:    make(chan []reflect.Value, b.handlerQueueSize),
	}

	go func() {
		for args := range h.queue {
			h.callback.Call(args)
		}
	}()

	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.handlers[topic] = append(b.handlers[topic], h)

	return nil
}

func (b *messageBus) Unsubscribe(topic string, fn any) error {
	if err := isValidHandler(fn); err != nil {
		return err
	}

	rv := reflect.ValueOf(fn)

	b.mtx.Lock()
	defer b.mtx.Unlock()

	if _, ok := b.handlers[topic]; ok {
		for i, h := range b.handlers[topic] {
			if h.callback == rv {
				close(h.queue)

				if len(b.handlers[topic]) == 1 {
					delete(b.handlers, topic)
				} else {
					b.handlers[topic] = append(b.handlers[topic][:i], b.handlers[topic][i+1:]...)
				}
			}
		}

		return nil
	}

	return fmt.Errorf("topic %s doesn't exist", topic)
}

func (b *messageBus) Close(topic string) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if _, ok := b.handlers[topic]; ok {
		for _, h := range b.handlers[topic] {
			close(h.queue)
		}

		delete(b.handlers, topic)

		return
	}
}

func isValidHandler(fn any) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	return nil
}

func buildHandlerArgs(args []any) []reflect.Value {
	reflectedArgs := make([]reflect.Value, 0)

	for _, arg := range args {
		reflectedArgs = append(reflectedArgs, reflect.ValueOf(arg))
	}

	return reflectedArgs
}
