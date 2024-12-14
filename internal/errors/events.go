package errors

import (
	"context"
	"sync"
)

type ErrorEventHandler interface {
	HandleError(ctx context.Context, err *Error)
}

type ErrorEventBus struct {
	mu       sync.RWMutex
	handlers map[ErrorCategory][]ErrorEventHandler
}

func NewErrorEventBus() *ErrorEventBus {
	return &ErrorEventBus{
		handlers: make(map[ErrorCategory][]ErrorEventHandler),
	}
}

func (bus *ErrorEventBus) Subscribe(category ErrorCategory, handler ErrorEventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if _, exists := bus.handlers[category]; !exists {
		bus.handlers[category] = make([]ErrorEventHandler, 0)
	}
	bus.handlers[category] = append(bus.handlers[category], handler)
}

func (bus *ErrorEventBus) Publish(ctx context.Context, err *Error) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	if handlers, exists := bus.handlers[err.Category]; exists {
		for _, handler := range handlers {
			go handler.HandleError(ctx, err)
		}
	}
} 