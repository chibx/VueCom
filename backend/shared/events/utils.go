package events

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Event struct {
	Type eventType
	Data any
}

type eventResponse struct {
	Data any
}

type Subscriber func(context.Context, *Event) any

// A very simple event emitter
// Should suffice
type eventemitter struct {
	mu   sync.Mutex
	subs map[eventType]Subscriber
	// channels map[eventType][]chan Event
	// running map[eventType]int64
	// max     int
}

func (e *eventemitter) Suscribe(etype eventType, sub Subscriber) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.subs[etype]
	if ok {
		panic(fmt.Sprintf("Only a single suscriber for event %s is needed for this struct", etype))
	}
	e.subs[etype] = sub

	return nil
}

func (e *eventemitter) Publish(ctx context.Context, event *Event) (*eventResponse, error) {
	c := make(chan any)
	var result any
	e.mu.Lock()
	sub, ok := e.subs[event.Type]
	e.mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("No Subscribers for %s event", event.Type)
	}

	go func() {
		data := sub(ctx, event)
		c <- data
	}()

	select {
	case <-ctx.Done():
		return nil, errors.New("Context Timeout occurred for event " + string(event.Type))
	case result = <-c:
		{
			if result == nil {
				return nil, nil
			}

			return &eventResponse{Data: result}, nil
		}
	}
}

var DefaultEmitter = eventemitter{
	subs: make(map[eventType]Subscriber),
}
