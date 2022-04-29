package evol

import (
	"time"
)

type Topic string

type Event interface {
	// Topic specify a class of topic, usually used in pub and sub
	Topic() Topic
	// Data contains Event payload
	Data() interface{}
	// AggregateType is the type of aggregate that publish the codec
	AggregateType() AggregateType
	// AggregateIdentity is the identity of aggregate that publish the codec
	AggregateIdentity() string
	// Timestamp of when the codec was created.
	Timestamp() time.Time
}

type event struct {
	topic         Topic
	data          interface{}
	aggregateType AggregateType
	aggregateId   string
	time          time.Time
}

func (b *event) Topic() Topic {
	return b.topic
}

func (b *event) Data() interface{} {
	return b.data
}

func (b *event) AggregateType() AggregateType {
	return b.aggregateType
}

func (b *event) AggregateIdentity() string {
	return b.aggregateId
}
func (b *event) Timestamp() time.Time {
	return b.time
}

// EventOption is an option to use when creating events.
type EventOption func(Event)

func NewEvent(topic Topic, data interface{}, time time.Time, options ...EventOption) Event {
	e := &event{
		topic: topic,
		data:  data,
		time:  time,
	}

	for _, option := range options {
		if option == nil {
			continue
		}
		option(e)
	}

	return e
}

// ForAggregate adds aggregate data when creating an codec.
func ForAggregate(aggregateType AggregateType, identity string) EventOption {
	return func(e Event) {
		if evt, ok := e.(*event); ok {
			evt.aggregateType = aggregateType
			evt.aggregateId = identity
		}
	}
}
