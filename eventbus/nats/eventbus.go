package nats

import (
	"context"
	"errors"
	"evol"
	"evol/codec"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

type EventBus struct {
	appID      string
	streamName string
	conn       *nats.Conn
	js         nats.JetStreamContext
	stream     *nats.StreamInfo
	connOpts   []nats.Option

	errCh  chan error
	cctx   context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	codec  evol.EventCodec
}

// NewEventBus creates an EventBus, with optional settings.
func NewEventBus(url, appID string, options ...Option) (*EventBus, error) {
	ctx, cancel := context.WithCancel(context.Background())

	b := &EventBus{
		appID:      appID,
		streamName: appID + "_events",

		errCh:  make(chan error, 100),
		cctx:   ctx,
		cancel: cancel,
		codec:  &codec.JsonEventCodec{},
	}

	// Apply configuration options.
	for _, option := range options {
		if option == nil {
			continue
		}

		if err := option(b); err != nil {
			return nil, fmt.Errorf("error while applying option: %w", err)
		}
	}

	// Create the NATS connection.
	var err error
	if b.conn, err = nats.Connect(url, b.connOpts...); err != nil {
		return nil, fmt.Errorf("could not create NATS connection: %w", err)
	}

	// Create Jetstream context.
	if b.js, err = b.conn.JetStream(); err != nil {
		return nil, fmt.Errorf("could not create Jetstream context: %w", err)
	}

	if b.stream, err = b.js.StreamInfo(b.streamName); err == nil {
		return b, nil
	}

	// Create the stream, which stores messages received on the subject.
	subjects := b.streamName + ".*.*"
	cfg := &nats.StreamConfig{
		Name:     b.streamName,
		Subjects: []string{subjects},
		Storage:  nats.FileStorage,
	}

	if b.stream, err = b.js.AddStream(cfg); err != nil {
		return nil, fmt.Errorf("could not create NATS stream: %w", err)
	}

	return b, nil
}

// Option is an option setter used to configure creation.
type Option func(*EventBus) error

// WithCodec uses the specified codec for encoding events.
func WithCodec(codec evol.EventCodec) Option {
	return func(b *EventBus) error {
		b.codec = codec

		return nil
	}
}

// WithNATSOptions adds the NATS options to the underlying client.
func WithNATSOptions(opts ...nats.Option) Option {
	return func(b *EventBus) error {
		b.connOpts = opts

		return nil
	}
}

func (b *EventBus) HandleEvent(ctx context.Context, event evol.Event) error {
	data, err := b.codec.MarshalEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("could not marshal codec: %w", err)
	}

	subject := fmt.Sprintf("%s.%s.%s", b.streamName, event.AggregateType(), event.Topic())
	if _, err := b.js.Publish(subject, data); err != nil {
		return fmt.Errorf("could not publish codec: %w", err)
	}

	return nil
}

func (b *EventBus) RegisterHandler(ctx context.Context, topic evol.Topic, eh evol.EventHandler) error {
	if topic == "" {
		return errors.New("missing codec topic")
	}

	if eh == nil {
		return errors.New("missing codec handler")
	}

	// Create a consumer.
	subject := fmt.Sprintf("%s.%s.%s", b.streamName, "*", topic)
	consumerName := fmt.Sprintf("%s_%s", b.appID, topic)

	sub, err := b.js.QueueSubscribe(subject, consumerName, b.handler(b.cctx, eh),
		nats.Durable(consumerName),
		nats.DeliverNew(),
		nats.ManualAck(),
		nats.AckExplicit(),
		nats.AckWait(60*time.Second),
		nats.MaxDeliver(10),
	)
	if err != nil {
		return fmt.Errorf("could not subscribe to queue: %w", err)
	}

	b.wg.Add(1)

	// Handle until context is cancelled.
	go b.handle(sub)

	return nil
}

func (b *EventBus) Close() error {
	b.cancel()
	b.wg.Wait()

	b.conn.Close()

	return nil
}

func (b *EventBus) handle(sub *nats.Subscription) {
	defer b.wg.Done()

	for {
		select {
		case <-b.cctx.Done():
			if b.cctx.Err() != context.Canceled {
				log.Printf("eventhorizon: context error in NATS codec bus: %s", b.cctx.Err())
			}

			return
		}
	}
}

func (b *EventBus) handler(ctx context.Context, eh evol.EventHandler) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		event, err := b.codec.UnmarshalEvent(ctx, msg.Data)
		if err != nil {
			err = fmt.Errorf("could not unmarshal codec: %w", err)
			msg.Nak()

			return
		}

		if err := eh.HandleEvent(ctx, event); err != nil {
			err = fmt.Errorf("could not handle codec : %w", err)

			msg.Nak()

			return
		}

		msg.AckSync()
	}
}
