package aggregatestore

import (
	"context"
	"errors"
	"evol"
)

//AggregateModelStore stores aggregate's latest status, usually used by the command bus(dispatcher)
type AggregateModelStore struct {
	repo evol.RWRepository
}

func (s *AggregateModelStore) Load(ctx context.Context, aggregateType evol.AggregateType, aggregateIdentity string) (evol.Aggregate, error) {
	entity, err := s.repo.Find(ctx, aggregateIdentity)
	if errors.Is(err, evol.ErrAggregateNotFound) {
		// Aggregate not found, create s new one
		aggregate, err2 := evol.CreateAggregate(aggregateType, aggregateIdentity)
		if err2 != nil {
			return nil, err2
		}
		return aggregate, nil
	} else if err != nil {
		return nil, err
	}

	//convert entity to Aggregate
	if aggregate, ok := entity.(evol.Aggregate); ok {
		return aggregate, nil
	}
	return nil, errors.New("[evol] Load Aggregate: invalid aggregate")

}

func (s *AggregateModelStore) Save(ctx context.Context, aggregate evol.Aggregate) error {
	if err := s.repo.Save(ctx, aggregate); err != nil {
		return err
	}

	return nil
}
