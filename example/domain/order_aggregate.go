package domain

import (
	"context"
	"evol"
	"github.com/mitchellh/mapstructure"
	"time"
)

func init() {
	// Register Aggregate
	evol.RegisterAggregate(OrderAggregateType, func(s string) evol.Aggregate {
		return &OrderAggregate{
			BaseAggregate: evol.NewBaseAggregate(OrderAggregateType, s),
			OrderId:       s, //OrderId as Aggregate Identity
		}
	})
}

var OrderAggregateType evol.AggregateType = "OrderAggregate"

//OrderAggregate Aggregate
type OrderAggregate struct {
	*evol.BaseAggregate
	OrderId    string
	BuyerId    string
	TotalPrice float32
	ProductIds []string
	Status     string
}

func (o *OrderAggregate) HandleCommand(ctx context.Context, c evol.Command) error {
	switch cmd := c.(type) {
	case *CreateOrderCmd:
		o.PublishEvent(OrderCreatedEventTopic, &OrderCreatedEvent{
			OrderId:    cmd.OrderId,
			BuyerId:    cmd.BuyerId,
			TotalPrice: cmd.Price,
			ProductIds: cmd.Goods,
		}, time.Now(), o)

	}
	return nil
}

func (o *OrderAggregate) HandleSourcingEvent(ctx context.Context, e evol.Event) error {
	switch e.Topic() {
	case OrderCreatedEventTopic:
		o.Status = "CREATED"
		ne := new(OrderCreatedEvent)
		err := mapstructure.Decode(e.Data(), ne)
		if err != nil {
			return err
		}
		o.OrderId = ne.OrderId
		o.BuyerId = ne.BuyerId
		o.ProductIds = ne.ProductIds
		o.TotalPrice = ne.TotalPrice

	}
	return nil
}
