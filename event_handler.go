package evol

import (
	"context"
	"sync"
)

type EventHandler interface {
	// HandleEvent handles an codec.
	HandleEvent(context.Context, Event) error
}

// EventHandlerFunc is a function that can be used as a codec handler.
type EventHandlerFunc func(context.Context, Event) error

func (f EventHandlerFunc) HandleEvent(ctx context.Context, e Event) error {
	return f(ctx, e)
}

// EventSourcingHandler handle codec published by this aggregate immediately(sync way)
type EventSourcingHandler interface {
	// Aggregate means this is an codec sourcing aggregate with
	Aggregate

	// HandleSourcingEvent handler codec they accepted
	HandleSourcingEvent(context.Context, Event) error
}

var EventHandlers = make(map[Topic][]EventHandler)
var EventHandlersMu sync.RWMutex

func RegisterEventHandlers(h EventHandler, topics ...Topic) {
	EventHandlersMu.Lock()
	defer EventHandlersMu.Unlock()
	for _, topic := range topics {
		EventHandlers[topic] = append(EventHandlers[topic], h)
	}
}
