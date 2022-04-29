package domain

import (
	"context"
	"evol"
	"github.com/mitchellh/mapstructure"
	"time"
)

func init() {
	// Register Aggregate
	evol.RegisterAggregate(StockAggregateType, func(s string) evol.Aggregate {
		return &StockAggregate{
			BaseAggregate: evol.NewBaseAggregate(StockAggregateType, s),
			ProductId:     s, //aggregate root identity
		}
	})
}

var StockAggregateType evol.AggregateType = "StockAggregate"

type StockAggregate struct {
	*evol.BaseAggregate
	ProductId string //aggregate root identity
	Quantity  int
}

func (s *StockAggregate) HandleCommand(ctx context.Context, c evol.Command) error {
	switch cmd := c.(type) {
	case *MakeReservationCmd:
		s.Quantity = Stock[s.ProductId] // mock, should be load from db firstly
		if s.Quantity >= cmd.Count {
			s.PublishEvent(ProductReservedEventTopic, &ProductReservedEvent{
				ProductId: cmd.ProductId,
				OrderId:   cmd.OrderId,
				Count:     cmd.Count,
			}, time.Now(), s)
		} else {
			s.PublishEvent(ProductReserveFailedEventTopic, &ProductReserveFailedEvent{
				ProductId: cmd.ProductId,
				OrderId:   cmd.OrderId,
				Count:     cmd.Count,
				Reason:    "not enough products",
			}, time.Now(), s)
		}

	}
	return nil
}

func (s *StockAggregate) HandleSourcingEvent(ctx context.Context, e evol.Event) error {
	switch e.Topic() {
	case ProductReservedEventTopic:
		evt := new(ProductReservedEvent)
		err := mapstructure.Decode(e.Data(), evt)
		if err != nil {
			return err
		}

		s.Quantity -= evt.Count
	case ProductReserveFailedEventTopic:

	}
	return nil
}
