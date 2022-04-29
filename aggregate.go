package evol

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"
)

type AggregateType string

func (a AggregateType) String() string {
	return string(a)
}

type Aggregate interface {
	// Entity Provides AggregateIdentity for command routing when loading from repository
	Entity
	// AggregateType specific a class of Aggregate for command routing
	AggregateType() AggregateType
	// CommandHandler for Aggregate handle command
	CommandHandler

	DomainEvents() []Event
}

//AggregateStore Saves Aggregate Status
type AggregateStore interface {
	//Load Aggregate latest status from repository with AggregateType and unique identity
	Load(ctx context.Context, aggregateType AggregateType, aggregateIdentity string) (Aggregate, error)
	//Save Aggregate latest status in repository
	Save(ctx context.Context, aggregate Aggregate) error
}

var aggFactories = make(map[AggregateType]func(string) Aggregate)
var aggFactoriesMu sync.RWMutex

func RegisterAggregate(aggregateType AggregateType, factory func(string) Aggregate) {
	ftype := reflect.TypeOf(factory)
	if ftype.Kind() != reflect.Func {
		panic("[evol] RegisterAggregate parameter not a func")
	}
	if ftype.NumIn() != 1 {
		panic("[evol] RegisterAggregate factory func parameter not match")
	}

	aggFactoriesMu.Lock()
	defer aggFactoriesMu.Unlock()

	if _, ok := aggFactories[aggregateType]; ok {
		panic(fmt.Sprintf("[evol] RegisterAggregate register duplicated aggregate: %s", aggregateType))
	}

	aggFactories[aggregateType] = factory

}

func CreateAggregate(aggType AggregateType, identity string) (Aggregate, error) {
	aggFactoriesMu.RLock()
	defer aggFactoriesMu.RUnlock()

	if factory, ok := aggFactories[aggType]; ok {
		return factory(identity), nil
	}

	return nil, errors.New("[evol] CreateAggregate aggregate not registered")
}

type BaseAggregate struct {
	identity string
	aType    AggregateType
	events   []Event
}

func (b *BaseAggregate) EntityIdentity() string {
	return b.identity
}

func (b *BaseAggregate) AggregateType() AggregateType {
	return b.aType
}

// PublishEvent publish a domain codec to current aggregate immediately with the esHandler
// and scheduled for other handlers when aggregate store(such as codec bus )
func (b *BaseAggregate) PublishEvent(topic Topic, data interface{}, timestamp time.Time, esHandler EventSourcingHandler, options ...EventOption) Event {
	options = append(options, ForAggregate(
		b.AggregateType(),
		b.EntityIdentity(),
	))
	e := NewEvent(topic, data, timestamp, options...)
	b.events = append(b.events, e)

	if err := esHandler.HandleSourcingEvent(context.Background(), e); err != nil {
		log.Println(err)
	}
	return e
}

func (b *BaseAggregate) DomainEvents() []Event {
	return b.events
}

func NewBaseAggregate(aggregateType AggregateType, identity string) *BaseAggregate {
	return &BaseAggregate{
		identity: identity,
		aType:    aggregateType,
	}
}
