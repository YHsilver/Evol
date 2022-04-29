package command

import (
	"context"
	"errors"
	"evol"
)

//AggCmdHandler is command handlers for a type of aggregate
type AggCmdHandler struct {
	aggregateType evol.AggregateType
	repo          evol.AggregateStore
	eventBus      evol.EventBus
}

func (h *AggCmdHandler) HandleCommand(ctx context.Context, cmd evol.Command) error {
	a, err := h.repo.Load(ctx, h.aggregateType, cmd.TargetIdentity())
	if err != nil {
		return err
	} else if a == nil {
		return errors.New("[evol] HandleCommand: Aggregate not found")
	}

	//here must be sync handle because aggregate status is stored later
	if err = a.HandleCommand(ctx, cmd); err != nil {
		return err
	}

	events := a.DomainEvents()

	for _, e := range events {
		if err := h.eventBus.HandleEvent(ctx, e); err != nil {
			return err
		}
	}

	return h.repo.Save(ctx, a)
}

//NewAggCmdHandler Create a evol.CommandHandler for an aggregate type
func NewAggCmdHandler(aggregateType evol.AggregateType, repo evol.AggregateStore, bus evol.EventBus) (*AggCmdHandler, error) {
	if repo == nil {
		return nil, errors.New("[evol] NewAggCmdHandler: aggregatestore nil")
	}

	h := &AggCmdHandler{
		aggregateType: aggregateType,
		repo:          repo,
		eventBus:      bus,
	}

	return h, nil
}
