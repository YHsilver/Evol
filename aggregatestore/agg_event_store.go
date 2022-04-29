package aggregatestore

import (
	"context"
	"errors"
	"evol"
	"fmt"
)

//AggregateEventStore stores aggregate's all codec instead of the latest status for codec sourcing
type AggregateEventStore struct {
	repo evol.EventRepo
}

func NewAggregateEventStore(repo evol.EventRepo) *AggregateEventStore {
	return &AggregateEventStore{repo: repo}
}

func (s *AggregateEventStore) Load(ctx context.Context, aggregateType evol.AggregateType, aggregateIdentity string) (evol.Aggregate, error) {
	//TODO: set tag as a snapshot of aggregate, thus the number of events needs tobe apple will be reduced
	agg, err := evol.CreateAggregate(aggregateType, aggregateIdentity)
	if err != nil {
		return nil, err
	}
	events, err := s.repo.Load(ctx, agg.EntityIdentity())
	if err != nil && !errors.Is(err, evol.ErrAggregateNotFound) {
		return nil, err
	}

	esAgg, ok := agg.(evol.EventSourcingHandler)
	if !ok {
		return nil, fmt.Errorf("[evol] AggregateEventStore Load error: aggregate %s is not a codec sourcing aggregate", agg)
	}

	if err := s.applyEvents(ctx, esAgg, events); err != nil {
		return nil, err
	}
	return agg, nil
}

func (s *AggregateEventStore) Save(ctx context.Context, agg evol.Aggregate) error {
	esAgg, ok := agg.(evol.EventSourcingHandler)
	if !ok {
		return fmt.Errorf("[evol] AggregateEventStore Save error: aggregate %s is not a codec sourcing aggregate", agg)
	}
	events := esAgg.DomainEvents()

	if err := s.repo.Save(ctx, events); err != nil {
		return err
	}
	return nil
}

func (s *AggregateEventStore) applyEvents(ctx context.Context, agg evol.EventSourcingHandler, events []evol.Event) error {
	for _, event := range events {
		if event.AggregateType() != agg.AggregateType() {
			return fmt.Errorf("[evol] aggeventrepo applyEvents aggregateType not match%s", event)
		}

		if err := agg.HandleSourcingEvent(ctx, event); err != nil {
			return fmt.Errorf("[evol] aggeventrepo applyEvents apply codec error %s: %w", event, err)
		}

	}
	return nil
}
