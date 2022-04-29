package memory

import (
	"context"
	"errors"
	"evol"
	"sync"
)

type AggregateEventRepo struct {
	db   map[string]AggregateRecord
	dbMu sync.RWMutex
}

func NewAggregateEventRepo() *AggregateEventRepo {
	return &AggregateEventRepo{
		db: make(map[string]AggregateRecord),
	}
}

type AggregateRecord struct {
	identity string
	aType    evol.AggregateType
	events   []evol.Event
}

func (r *AggregateEventRepo) Save(ctx context.Context, events []evol.Event) error {
	r.dbMu.Lock()
	defer r.dbMu.Unlock()

	if events == nil || len(events) == 0 {
		return errors.New("save events are empty")
	}

	id := events[0].AggregateIdentity()
	at := events[0].AggregateType()

	ar, ok := r.db[id]
	if !ok {
		ar = AggregateRecord{
			identity: id,
			aType:    at,
			events:   events,
		}
		r.db[id] = ar
	} else {
		ar.events = append(ar.events, events...)
		r.db[id] = ar
	}

	return nil
}

func (r *AggregateEventRepo) Load(ctx context.Context, id string) ([]evol.Event, error) {
	ar, ok := r.db[id]
	if !ok {
		return nil, evol.ErrAggregateNotFound
	}

	events := make([]evol.Event, len(ar.events))
	for i, e := range ar.events {
		events[i] = e
	}

	return events, nil
}
