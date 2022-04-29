package local

import (
	"context"
	"errors"
	"evol"
	"fmt"
	"log"
	"sync"
	"time"
)

// DefaultQueueSize is the default queue size per handler for publishing events.
var DefaultQueueSize = 1000

// EventBus is a local codec bus that dispatch events to all matching handlers
type EventBus struct {
	group *Group

	ctx    context.Context
	cancel context.CancelFunc
	codec  evol.EventCodec
	wg     sync.WaitGroup
}

// NewEventBus creates a EventBus.
func NewEventBus(options ...Option) *EventBus {
	ctx, cancel := context.WithCancel(context.Background())

	b := &EventBus{
		group:  NewGroup(),
		ctx:    ctx,
		cancel: cancel,
	}

	// Apply configuration options.
	for _, option := range options {
		if option == nil {
			continue
		}

		option(b)
	}

	return b
}

type Option func(*EventBus)

func WithGroup(g *Group) Option {
	return func(b *EventBus) {
		b.group = g
	}
}

// HandleEvent implements the HandleEvent method of the codec.EventHandler interface.
func (b *EventBus) HandleEvent(ctx context.Context, event evol.Event) error {
	data, err := b.codec.MarshalEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("could not marshal codec: %w", err)
	}

	return b.group.publish(ctx, string(event.Topic()), data)
}

// RegisterHandler implements the RegisterHandler method of the codec.EventBus interface.
func (b *EventBus) RegisterHandler(ctx context.Context, topic evol.Topic, h evol.EventHandler) error {

	if h == nil {
		return errors.New("missing codec handler")
	}

	ch := b.group.channel(string(topic))

	b.wg.Add(1)

	// Handle until context is cancelled.
	go b.handle(ctx, topic, h, ch)

	return nil
}

func (b *EventBus) Close() error {
	b.cancel()
	b.wg.Wait()
	b.group.close()

	return nil
}

// Handles all events coming in on the channel.
func (b *EventBus) handle(ctx context.Context, topic evol.Topic, h evol.EventHandler, ch <-chan []byte) {
	defer b.wg.Done()

	for {
		select {
		case data := <-ch:

			time.Sleep(time.Millisecond)

			e, err := b.codec.UnmarshalEvent(ctx, data)
			if err != nil {
				err = fmt.Errorf("could not unmarshal e: %w", err)
				return
			}

			if err := h.HandleEvent(ctx, e); err != nil {
				_ = fmt.Errorf("could not handle e (%s): %s", e, err.Error())

			}
		case <-b.ctx.Done():
			return
		}
	}
}

//Group publish events separate by topic
type Group struct {
	bus   map[string]chan []byte
	busMu sync.RWMutex
}

func NewGroup() *Group {
	return &Group{
		bus: map[string]chan []byte{},
	}
}

func (g *Group) channel(id string) <-chan []byte {
	g.busMu.Lock()
	defer g.busMu.Unlock()

	if ch, ok := g.bus[id]; ok {
		return ch
	}

	ch := make(chan []byte, DefaultQueueSize)
	g.bus[id] = ch

	return ch
}

func (g *Group) publish(ctx context.Context, topic string, b []byte) error {
	g.busMu.RLock()
	defer g.busMu.RUnlock()

	ch := g.bus[topic]

	select {
	case ch <- b:
	default:
		log.Printf("[evol] publish queue full in local codec bus")
	}

	return nil
}

// Closes all the open channels after handling is done.
func (g *Group) close() {
	for _, ch := range g.bus {
		close(ch)
	}

	g.bus = nil
}
